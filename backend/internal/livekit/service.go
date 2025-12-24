package livekit

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/blytz.live.remake/backend/internal/models"
	"github.com/blytz.live.remake/backend/pkg/logging"
	"github.com/google/uuid"
	"github.com/livekit/protocol/auth"
	"github.com/livekit/protocol/livekit"
	livekitSdk "github.com/livekit/server-sdk-go"
	"gorm.io/gorm"
)

// Service provides LiveKit live streaming services
type Service struct {
	db        *gorm.DB
	logger    logging.Logger
	host      string
	apiKey    string
	apiSecret string
	roomConf  *livekit.RoomConfiguration
}

// NewService creates a new LiveKit service
func NewService(db *gorm.DB, host, apiKey, apiSecret string) *Service {
	logger := logging.NewLogger("livekit-service")
	
	return &Service{
		db:        db,
		logger:    logger,
		host:      host,
		apiKey:    apiKey,
		apiSecret: apiSecret,
		roomConf: &livekit.RoomConfiguration{
			MaxParticipants: 1000, // Max viewers per auction
		},
	}
}

// CreateAuctionRoom creates a LiveKit room for an auction
func (s *Service) CreateAuctionRoom(ctx context.Context, auctionID uuid.UUID) (*LiveKitRoom, error) {
	roomName := fmt.Sprintf("auction-%s", auctionID.String())
	
	// Connect to LiveKit server
	roomClient := livekitSdk.NewRoomServiceClient(s.host, s.apiKey, s.apiSecret)
	
	// Create room
	room, err := roomClient.CreateRoom(ctx, &livekit.CreateRoomRequest{
		Name:             roomName,
		EmptyTimeout:      300, // 5 minutes
		MaxParticipants:   1000,
		EnabledCodecs: []*livekit.Codec{
			{
				Mime:     "video/vp8",
				Disabled: false,
			},
			{
				Mime:     "video/h264",
				Disabled: false,
			},
			{
				Mime:     "audio/opus",
				Disabled: false,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create LiveKit room: %w", err)
	}
	
	// Create room record in database
	liveKitRoom := &models.LiveStream{
		AuctionID:   auctionID,
		StreamURL:    fmt.Sprintf("%s/room/%s", s.host, roomName),
		StreamKey:    uuid.New().String(),
		Status:       "waiting",
		ViewerCount:  0,
		IsRecording:  true,
		ThumbnailURL: nil,
	}
	
	if err := s.db.WithContext(ctx).Create(liveKitRoom).Error; err != nil {
		return nil, fmt.Errorf("failed to save live stream record: %w", err)
	}
	
	// Generate room data
	roomData := &LiveKitRoom{
		RoomName:    room.Name,
		RoomID:      room.Sid,
		StreamURL:    liveKitRoom.StreamURL,
		StreamKey:    liveKitRoom.StreamKey,
		PlaybackURL:  fmt.Sprintf("%s/playback?room=%s", s.host, roomName),
		Status:       liveKitRoom.Status,
		CreatedAt:    room.CreationTime,
		DatabaseID:   liveKitRoom.ID,
	}
	
	s.logger.Info("LiveKit room created", map[string]interface{}{
		"auction_id": auctionID,
		"room_name":  roomName,
		"room_id":   room.Sid,
	})
	
	return roomData, nil
}

// GenerateViewerToken generates a token for viewers to join auction stream
func (s *Service) GenerateViewerToken(ctx context.Context, auctionID uuid.UUID, userID *uuid.UUID) (string, error) {
	roomName := fmt.Sprintf("auction-%s", auctionID.String())
	
	// Set token permissions
	grant := &auth.VideoGrant{
		RoomJoin: true,
		Room:      roomName,
		// Viewers can only subscribe, not publish
		CanPublish:   false,
		CanSubscribe: true,
	}
	
	// Add user identity if provided
	identity := ""
	if userID != nil {
		identity = userID.String()
		grant.CanPublish = true // Allow publishing chat messages
	}
	
	// Generate token
	at := auth.NewAccessToken(s.apiKey, s.apiSecret)
	at.SetIdentity(identity)
	at.SetVideoGrant(grant)
	at.SetValidFor(2 * time.Hour) // 2 hour validity
	
	token, err := at.ToJWT()
	if err != nil {
		return "", fmt.Errorf("failed to generate viewer token: %w", err)
	}
	
	return token, nil
}

// GenerateHostToken generates a token for auction host/seller
func (s *Service) GenerateHostToken(ctx context.Context, auctionID, userID uuid.UUID) (string, error) {
	roomName := fmt.Sprintf("auction-%s", auctionID.String())
	
	// Set token permissions for host
	grant := &auth.VideoGrant{
		RoomJoin:     true,
		Room:         roomName,
		CanPublish:   true,   // Can stream video/audio
		CanSubscribe: true,   // Can see viewers
		RoomAdmin:    true,    // Can manage room
	}
	
	// Generate token
	at := auth.NewAccessToken(s.apiKey, s.apiSecret)
	at.SetIdentity(userID.String())
	at.SetName("auction-host")
	at.SetVideoGrant(grant)
	at.SetValidFor(4 * time.Hour) // 4 hour validity for auction
	
	token, err := at.ToJWT()
	if err != nil {
		return "", fmt.Errorf("failed to generate host token: %w", err)
	}
	
	return token, nil
}

// StartAuctionStream starts the auction stream
func (s *Service) StartAuctionStream(ctx context.Context, auctionID uuid.UUID) error {
	roomName := fmt.Sprintf("auction-%s", auctionID.String())
	
	// Update database record
	updates := map[string]interface{}{
		"status":    "live",
		"started_at": time.Now(),
	}
	
	err := s.db.WithContext(ctx).
		Model(&models.LiveStream{}).
		Where("auction_id = ?", auctionID).
		Updates(updates).Error
	
	if err != nil {
		return fmt.Errorf("failed to update live stream status: %w", err)
	}
	
	s.logger.Info("Auction stream started", map[string]interface{}{
		"auction_id": auctionID,
		"room_name":  roomName,
	})
	
	return nil
}

// EndAuctionStream ends the auction stream
func (s *Service) EndAuctionStream(ctx context.Context, auctionID uuid.UUID) error {
	roomName := fmt.Sprintf("auction-%s", auctionID.String())
	
	// Get stream record
	var stream models.LiveStream
	err := s.db.WithContext(ctx).First(&stream, "auction_id = ?", auctionID).Error
	if err != nil {
		return fmt.Errorf("stream not found: %w", err)
	}
	
	// Calculate duration
	duration := 0
	if stream.StartedAt != nil {
		duration = int(time.Since(*stream.StartedAt).Seconds())
	}
	
	// Update database record
	updates := map[string]interface{}{
		"status":    "ended",
		"ended_at":  time.Now(),
		"duration":  duration,
	}
	
	err = s.db.WithContext(ctx).Model(&stream).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("failed to update live stream status: %w", err)
	}
	
	// Optional: Delete LiveKit room to free resources
	// roomClient := livekitSdk.NewRoomServiceClient(s.host, s.apiKey, s.apiSecret)
	// roomClient.DeleteRoom(ctx, &livekit.DeleteRoomRequest{Room: roomName})
	
	s.logger.Info("Auction stream ended", map[string]interface{}{
		"auction_id": auctionID,
		"room_name":  roomName,
		"duration":   duration,
	})
	
	return nil
}

// GetStreamInfo gets current stream information for an auction
func (s *Service) GetStreamInfo(ctx context.Context, auctionID uuid.UUID) (*models.LiveStream, error) {
	var stream models.LiveStream
	err := s.db.WithContext(ctx).
		Preload("Auction").
		Preload("Auction.Product").
		First(&stream, "auction_id = ?", auctionID).Error
	
	return &stream, err
}

// UpdateViewerCount updates the viewer count for a stream
func (s *Service) UpdateViewerCount(ctx context.Context, auctionID uuid.UUID, viewerCount int) error {
	err := s.db.WithContext(ctx).
		Model(&models.LiveStream{}).
		Where("auction_id = ?", auctionID).
		Update("viewer_count", viewerCount).Error
	
	return err
}

// UpdateStreamMetrics updates stream performance metrics
func (s *Service) UpdateStreamMetrics(ctx context.Context, auctionID uuid.UUID, metrics *StreamMetrics) error {
	updates := map[string]interface{}{
		"viewer_count": metrics.ViewerCount,
		"latency":     metrics.Latency,
		"bandwidth":   metrics.Bandwidth,
	}
	
	err := s.db.WithContext(ctx).
		Model(&models.LiveStream{}).
		Where("auction_id = ?", auctionID).
		Updates(updates).Error
	
	return err
}

// ListActiveStreams lists all currently active streams
func (s *Service) ListActiveStreams(ctx context.Context) ([]models.LiveStream, error) {
	var streams []models.LiveStream
	err := s.db.WithContext(ctx).
		Preload("Auction").
		Preload("Auction.Product").
		Preload("Auction.Seller").
		Where("status = ?", "live").
		Find(&streams).Error
	
	return streams, err
}

// RecordStream enables recording for a stream
func (s *Service) RecordStream(ctx context.Context, auctionID uuid.UUID, enabled bool) error {
	updates := map[string]interface{}{
		"is_recording": enabled,
	}
	
	err := s.db.WithContext(ctx).
		Model(&models.LiveStream{}).
		Where("auction_id = ?", auctionID).
		Updates(updates).Error
	
	return err
}

// GenerateStreamRecordingURL generates a URL for recorded stream
func (s *Service) GenerateStreamRecordingURL(ctx context.Context, auctionID uuid.UUID) (string, error) {
	var stream models.LiveStream
	err := s.db.WithContext(ctx).First(&stream, "auction_id = ?", auctionID).Error
	if err != nil {
		return "", fmt.Errorf("stream not found: %w", err)
	}
	
	if !stream.IsRecording {
		return "", fmt.Errorf("recording not enabled for this stream")
	}
	
	// Generate recording URL (this would integrate with your recording storage)
	recordingURL := fmt.Sprintf("%s/recordings/%s.mp4", s.host, auctionID.String())
	
	// Update stream with recording URL
	s.db.WithContext(ctx).
		Model(&stream).
		Update("recording_url", recordingURL)
	
	return recordingURL, nil
}

// LiveKitRoom represents LiveKit room information
type LiveKitRoom struct {
	RoomName    string    `json:"room_name"`
	RoomID      string    `json:"room_id"`
	StreamURL    string    `json:"stream_url"`
	StreamKey    string    `json:"stream_key"`
	PlaybackURL  string    `json:"playback_url"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	DatabaseID   uuid.UUID `json:"database_id"`
}

// StreamMetrics represents stream performance metrics
type StreamMetrics struct {
	ViewerCount int  `json:"viewer_count"`
	Latency     int  `json:"latency"`     // in milliseconds
	Bandwidth   int  `json:"bandwidth"`   // in kbps
	CPUUsage    int  `json:"cpu_usage"`    // percentage
	MemoryUsage int  `json:"memory_usage"` // percentage
}

// GetRoomParticipants gets current participants in a room
func (s *Service) GetRoomParticipants(ctx context.Context, auctionID uuid.UUID) ([]*ParticipantInfo, error) {
	roomName := fmt.Sprintf("auction-%s", auctionID.String())
	
	roomClient := livekitSdk.NewRoomServiceClient(s.host, s.apiKey, s.apiSecret)
	
	// Get room participants
	participants, err := roomClient.ListParticipants(ctx, &livekit.ListParticipantsRequest{
		Room: roomName,
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to get room participants: %w", err)
	}
	
	// Convert to our format
	var participantInfos []*ParticipantInfo
	for _, p := range participants.Participants {
		info := &ParticipantInfo{
			SID:       p.Sid,
			Identity:  p.Identity,
			Name:      p.Name,
			IsPublisher: p.Permission != nil && p.Permission.CanPublish,
			JoinedAt:   time.Unix(p.JoinedAt/1000, 0),
		}
		participantInfos = append(participantInfos, info)
	}
	
	return participantInfos, nil
}

// ParticipantInfo represents a room participant
type ParticipantInfo struct {
	SID         string    `json:"sid"`
	Identity    string    `json:"identity"`
	Name        string    `json:"name"`
	IsPublisher bool      `json:"is_publisher"`
	JoinedAt    time.Time `json:"joined_at"`
}

// RemoveParticipant removes a participant from a room
func (s *Service) RemoveParticipant(ctx context.Context, auctionID uuid.UUID, participantSID string) error {
	roomName := fmt.Sprintf("auction-%s", auctionID.String())
	
	roomClient := livekitSdk.NewRoomServiceClient(s.host, s.apiKey, s.apiSecret)
	
	// Remove participant
	err := roomClient.RemoveParticipant(ctx, &livekit.RoomParticipantIdentity{
		Room:        roomName,
		Identity:     participantSID,
	})
	
	if err != nil {
		return fmt.Errorf("failed to remove participant: %w", err)
	}
	
	s.logger.Info("Participant removed from auction room", map[string]interface{}{
		"auction_id":      auctionID,
		"participant_sid": participantSID,
	})
	
	return nil
}

// MuteParticipant mutes/unmutes a participant
func (s *Service) MuteParticipant(ctx context.Context, auctionID uuid.UUID, participantSID string, muted bool) error {
	roomName := fmt.Sprintf("auction-%s", auctionID.String())
	
	roomClient := livekitSdk.NewRoomServiceClient(s.host, s.apiKey, s.apiSecret)
	
	// Mute participant
	_, err := roomClient.MutePublishedTrack(ctx, &livekit.MuteRoomTrackRequest{
		Room:     roomName,
		Identity:  participantSID,
		Muted:    muted,
	})
	
	if err != nil {
		return fmt.Errorf("failed to mute participant: %w", err)
	}
	
	return nil
}