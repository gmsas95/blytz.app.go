package models

import (
	"time"

	"github.com/blytz.live.remake/backend/internal/common"
	"github.com/google/uuid"
)

// Auction represents a live auction session
type Auction struct {
	common.BaseModel
	ProductID    uuid.UUID  `gorm:"not null;references:ID" json:"product_id"`
	Product      Product    `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	SellerID     uuid.UUID  `gorm:"not null;references:ID" json:"seller_id"`
	Seller       User       `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
	Title        string     `gorm:"not null" json:"title"`
	Description  *string    `json:"description"`
	StartTime    time.Time  `gorm:"not null" json:"start_time"`
	EndTime      time.Time  `gorm:"not null" json:"end_time"`
	Status       string     `gorm:"default:'scheduled'" json:"status"` // scheduled, live, ended, cancelled
	StartPrice   float64    `gorm:"not null" json:"start_price"`
	ReservePrice *float64   `json:"reserve_price"`
	BuyNowPrice  *float64   `json:"buy_now_price"`
	CurrentBid   *float64   `json:"current_bid"`
	BidCount     int        `gorm:"default:0" json:"bid_count"`
	WinnerID     *uuid.UUID `gorm:"references:ID" json:"winner_id"`
	Winner       *User      `gorm:"foreignKey:WinnerID" json:"winner,omitempty"`
	LiveKitRoom  string     `gorm:"uniqueIndex;not null" json:"livekit_room"`
	StreamKey    *string    `json:"stream_key"`
	AutoExtend   bool       `gorm:"default:true" json:"auto_extend"` // Auto-extend if bid in last 5 minutes
	ExtendTime   int        `gorm:"default:300" json:"extend_time"`  // Extend time in seconds
	IsFeatured   bool       `gorm:"default:false" json:"is_featured"`
	Bids         []Bid      `gorm:"foreignKey:AuctionID" json:"bids,omitempty"`
}

// Bid represents a bid in an auction
type Bid struct {
	common.BaseModel
	AuctionID uuid.UUID `gorm:"not null;references:ID;index" json:"auction_id"`
	Auction   Auction   `gorm:"foreignKey:AuctionID" json:"auction,omitempty"`
	UserID    uuid.UUID `gorm:"not null;references:ID;index" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Amount    float64   `gorm:"not null" json:"amount"`
	IsAutoBid bool      `gorm:"default:false" json:"is_auto_bid"`
	IsWinning bool      `gorm:"default:false" json:"is_winning"`
	BidTime   time.Time `gorm:"autoCreateTime" json:"bid_time"`
}

// AutoBid represents an automatic bidding configuration
type AutoBid struct {
	common.BaseModel
	AuctionID    uuid.UUID `gorm:"not null;references:ID;index" json:"auction_id"`
	Auction      Auction   `gorm:"foreignKey:AuctionID" json:"auction,omitempty"`
	UserID       uuid.UUID `gorm:"not null;references:ID;index" json:"user_id"`
	User         User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	MaxAmount    float64   `gorm:"not null" json:"max_amount"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CurrentBid   *float64  `json:"current_bid"`
	BidIncrement float64   `gorm:"not null;default:5.0" json:"bid_increment"`
	LastBidTime  *time.Time `json:"last_bid_time"`
}

// AuctionWatch represents users watching an auction
type AuctionWatch struct {
	common.BaseModel
	AuctionID   uuid.UUID `gorm:"not null;references:ID;index" json:"auction_id"`
	Auction     Auction   `gorm:"foreignKey:AuctionID" json:"auction,omitempty"`
	UserID      uuid.UUID `gorm:"not null;references:ID;index" json:"user_id"`
	User        User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	JoinedAt    time.Time `gorm:"autoCreateTime" json:"joined_at"`
	LastActive  time.Time `gorm:"autoUpdateTime" json:"last_active"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	NotificationSettings string `gorm:"type:jsonb" json:"notification_settings"` // JSON: {"outbid": true, "ending_soon": true}
}

// AuctionStats represents auction statistics
type AuctionStats struct {
	common.BaseModel
	AuctionID      uuid.UUID `gorm:"not null;uniqueIndex;references:ID" json:"auction_id"`
	Auction        Auction   `gorm:"foreignKey:AuctionID" json:"auction,omitempty"`
	TotalBids      int       `gorm:"default:0" json:"total_bids"`
	TotalBidders   int       `gorm:"default:0" json:"total_bidders"`
	UniqueViewers  int       `gorm:"default:0" json:"unique_viewers"`
	PeakViewers    int       `gorm:"default:0" json:"peak_viewers"`
	AvgWatchTime   int       `gorm:"default:0" json:"avg_watch_time"` // in minutes
	EngagementRate float64   `gorm:"default:0" json:"engagement_rate"`
	Revenue        float64   `gorm:"default:0" json:"revenue"`
}

// LiveStream represents a live streaming session
type LiveStream struct {
	common.BaseModel
	AuctionID     uuid.UUID  `gorm:"not null;uniqueIndex;references:ID" json:"auction_id"`
	Auction       Auction    `gorm:"foreignKey:AuctionID" json:"auction,omitempty"`
	StreamURL     string     `gorm:"not null" json:"stream_url"`
	StreamKey     string     `gorm:"not null" json:"stream_key"`
	PlaybackURL   *string    `json:"playback_url"`
	Status        string     `gorm:"default:'waiting'" json:"status"` // waiting, live, ended, error
	StartedAt     *time.Time `json:"started_at"`
	EndedAt       *time.Time `json:"ended_at"`
	Duration      int        `gorm:"default:0" json:"duration"` // in seconds
	ViewerCount   int        `gorm:"default:0" json:"viewer_count"`
	Latency       int        `gorm:"default:0" json:"latency"`   // in milliseconds
	Bandwidth     int        `gorm:"default:0" json:"bandwidth"` // in kbps
	RecordingURL  *string    `json:"recording_url"`
	ThumbnailURL  *string    `json:"thumbnail_url"`
	IsRecording   bool       `gorm:"default:false" json:"is_recording"`
}

// ChatMessage represents a chat message during a live auction
type ChatMessage struct {
	common.BaseModel
	AuctionID   uuid.UUID `gorm:"not null;references:ID;index" json:"auction_id"`
	Auction     Auction   `gorm:"foreignKey:AuctionID" json:"auction,omitempty"`
	UserID      uuid.UUID `gorm:"not null;references:ID;index" json:"user_id"`
	User        User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Message     string    `gorm:"type:text;not null" json:"message"`
	MessageType string    `gorm:"default:'user'" json:"message_type"` // user, system, moderator
	IsModerated bool      `gorm:"default:false" json:"is_moderated"`
	Timestamp   time.Time `gorm:"autoCreateTime" json:"timestamp"`
	ReplyTo     *uuid.UUID `gorm:"references:ID" json:"reply_to,omitempty"` // Reply to another message
}