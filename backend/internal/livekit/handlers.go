package livekit

import (
	"net/http"
	"strconv"

	"github.com/blytz.live.remake/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler provides LiveKit HTTP handlers
type Handler struct {
	service *Service
}

// NewHandler creates a new LiveKit handler
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// StartStreamRequest represents request body for starting a stream
type StartStreamRequest struct {
	AuctionID   uuid.UUID `json:"auction_id" binding:"required"`
	RecordStream bool      `json:"record_stream"`
}

// UpdateMetricsRequest represents request body for updating stream metrics
type UpdateMetricsRequest struct {
	ViewerCount int  `json:"viewer_count"`
	Latency     int  `json:"latency"`     // in milliseconds
	Bandwidth   int  `json:"bandwidth"`   // in kbps
	CPUUsage    int  `json:"cpu_usage"`    // percentage
	MemoryUsage int  `json:"memory_usage"` // percentage
}

// CreateAuctionRoom creates a LiveKit room for an auction
func (h *Handler) CreateAuctionRoom(c *gin.Context) {
	auctionID, err := uuid.Parse(c.Param("auction_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auction ID"})
		return
	}

	room, err := h.service.CreateAuctionRoom(c.Request.Context(), auctionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, room)
}

// GetViewerToken generates a token for viewers to join auction stream
func (h *Handler) GetViewerToken(c *gin.Context) {
	auctionID, err := uuid.Parse(c.Param("auction_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auction ID"})
		return
	}

	// Get user ID from context (optional for viewers)
	var userID *uuid.UUID
	if uid, exists := c.Get("user_id"); exists {
		if id, ok := uid.(uuid.UUID); ok {
			userID = &id
		}
	}

	token, err := h.service.GenerateViewerToken(c.Request.Context(), auctionID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      token,
		"auction_id": auctionID,
		"user_id":    userID,
	})
}

// GetHostToken generates a token for auction host/seller
func (h *Handler) GetHostToken(c *gin.Context) {
	auctionID, err := uuid.Parse(c.Param("auction_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auction ID"})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	token, err := h.service.GenerateHostToken(c.Request.Context(), auctionID, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      token,
		"auction_id": auctionID,
		"user_id":    userID,
	})
}

// StartAuctionStream starts an auction stream
func (h *Handler) StartAuctionStream(c *gin.Context) {
	auctionID, err := uuid.Parse(c.Param("auction_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auction ID"})
		return
	}

	// Check user permissions (must be seller or admin)
	role, _ := c.Get("role")
	if role != "seller" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	err = h.service.StartAuctionStream(c.Request.Context(), auctionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Auction stream started successfully"})
}

// EndAuctionStream ends an auction stream
func (h *Handler) EndAuctionStream(c *gin.Context) {
	auctionID, err := uuid.Parse(c.Param("auction_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auction ID"})
		return
	}

	// Check user permissions (must be seller or admin)
	role, _ := c.Get("role")
	if role != "seller" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	err = h.service.EndAuctionStream(c.Request.Context(), auctionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Auction stream ended successfully"})
}

// GetStreamInfo gets current stream information for an auction
func (h *Handler) GetStreamInfo(c *gin.Context) {
	auctionID, err := uuid.Parse(c.Param("auction_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auction ID"})
		return
	}

	stream, err := h.service.GetStreamInfo(c.Request.Context(), auctionID)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Stream not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stream)
}

// UpdateViewerCount updates viewer count for a stream
func (h *Handler) UpdateViewerCount(c *gin.Context) {
	auctionID, err := uuid.Parse(c.Param("auction_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auction ID"})
		return
	}

	var req struct {
		ViewerCount int `json:"viewer_count" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.UpdateViewerCount(c.Request.Context(), auctionID, req.ViewerCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Viewer count updated successfully"})
}

// UpdateStreamMetrics updates stream performance metrics
func (h *Handler) UpdateStreamMetrics(c *gin.Context) {
	auctionID, err := uuid.Parse(c.Param("auction_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auction ID"})
		return
	}

	var req UpdateMetricsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	metrics := &StreamMetrics{
		ViewerCount: req.ViewerCount,
		Latency:     req.Latency,
		Bandwidth:   req.Bandwidth,
		CPUUsage:    req.CPUUsage,
		MemoryUsage: req.MemoryUsage,
	}

	err = h.service.UpdateStreamMetrics(c.Request.Context(), auctionID, metrics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stream metrics updated successfully"})
}

// ListActiveStreams lists all currently active streams
func (h *Handler) ListActiveStreams(c *gin.Context) {
	streams, err := h.service.ListActiveStreams(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"streams": streams})
}

// RecordStream enables/disables recording for a stream
func (h *Handler) RecordStream(c *gin.Context) {
	auctionID, err := uuid.Parse(c.Param("auction_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auction ID"})
		return
	}

	var req struct {
		Enabled bool `json:"enabled" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check user permissions (must be seller or admin)
	role, _ := c.Get("role")
	if role != "seller" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	err = h.service.RecordStream(c.Request.Context(), auctionID, req.Enabled)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stream recording updated successfully"})
}

// GenerateStreamRecordingURL generates a URL for recorded stream
func (h *Handler) GenerateStreamRecordingURL(c *gin.Context) {
	auctionID, err := uuid.Parse(c.Param("auction_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auction ID"})
		return
	}

	// Check user permissions (must be seller or admin)
	role, _ := c.Get("role")
	if role != "seller" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	recordingURL, err := h.service.GenerateStreamRecordingURL(c.Request.Context(), auctionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"recording_url": recordingURL,
		"auction_id":   auctionID,
	})
}

// GetRoomParticipants gets current participants in a room
func (h *Handler) GetRoomParticipants(c *gin.Context) {
	auctionID, err := uuid.Parse(c.Param("auction_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auction ID"})
		return
	}

	// Check user permissions (must be seller or admin)
	role, _ := c.Get("role")
	if role != "seller" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	participants, err := h.service.GetRoomParticipants(c.Request.Context(), auctionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"participants": participants})
}

// RemoveParticipant removes a participant from a room
func (h *Handler) RemoveParticipant(c *gin.Context) {
	auctionID, err := uuid.Parse(c.Param("auction_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auction ID"})
		return
	}

	participantSID := c.Param("participant_sid")
	if participantSID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Participant SID is required"})
		return
	}

	// Check user permissions (must be seller or admin)
	role, _ := c.Get("role")
	if role != "seller" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	err = h.service.RemoveParticipant(c.Request.Context(), auctionID, participantSID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Participant removed successfully"})
}

// MuteParticipant mutes/unmutes a participant
func (h *Handler) MuteParticipant(c *gin.Context) {
	auctionID, err := uuid.Parse(c.Param("auction_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auction ID"})
		return
	}

	participantSID := c.Param("participant_sid")
	if participantSID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Participant SID is required"})
		return
	}

	muted := c.DefaultQuery("muted", "true")
	isMuted := muted == "true"

	// Check user permissions (must be seller or admin)
	role, _ := c.Get("role")
	if role != "seller" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	err = h.service.MuteParticipant(c.Request.Context(), auctionID, participantSID, isMuted)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Participant muted/unmuted successfully",
		"participant_sid": participantSID,
		"muted":        isMuted,
	})
}

// GetStreamRecording gets a recorded stream
func (h *Handler) GetStreamRecording(c *gin.Context) {
	auctionID, err := uuid.Parse(c.Param("auction_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auction ID"})
		return
	}

	stream, err := h.service.GetStreamInfo(c.Request.Context(), auctionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stream not found"})
		return
	}

	if stream.RecordingURL == nil || *stream.RecordingURL == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recording not available"})
		return
	}

	// Redirect to recording URL
	c.Redirect(http.StatusFound, *stream.RecordingURL)
}