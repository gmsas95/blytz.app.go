package auction

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusScheduled Status = "scheduled"
	StatusLive      Status = "live"
	StatusEnded     Status = "ended"
	StatusCancelled Status = "cancelled"
)

type Auction struct {
	ID           uuid.UUID
	ProductID    uuid.UUID
	SellerID     uuid.UUID
	Title        string
	Description  string
	StartTime    time.Time
	EndTime      time.Time
	Status       Status
	StartPrice   float64
	ReservePrice *float64
	BuyNowPrice  *float64
	CurrentBid   *Bid
	BidCount     int
	WinnerID     *uuid.UUID
	LiveKitRoom  string
	StreamKey    string
	AutoExtend   bool
	ExtendTime   time.Duration
	IsFeatured   bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (a *Auction) CanPlaceBid(amount float64, now time.Time) error {
	if a.Status != StatusLive {
		return errors.New("auction is not live")
	}
	if now.After(a.EndTime) {
		return errors.New("auction has ended")
	}
	minBid := a.StartPrice
	if a.CurrentBid != nil {
		minBid = a.CurrentBid.Amount + 1.0
	}
	if amount < minBid {
		return errors.New("bid amount too low")
	}
	return nil
}

func (a *Auction) PlaceBid(bidderID uuid.UUID, amount float64, now time.Time) (*Bid, error) {
	if err := a.CanPlaceBid(amount, now); err != nil {
		return nil, err
	}
	bid := &Bid{
		ID:        uuid.New(),
		AuctionID: a.ID,
		UserID:    bidderID,
		Amount:    amount,
		IsWinning: true,
		BidTime:   now,
	}
	a.CurrentBid = bid
	a.BidCount++
	if a.AutoExtend && time.Until(a.EndTime) <= 5*time.Minute {
		a.EndTime = a.EndTime.Add(a.ExtendTime)
	}
	a.UpdatedAt = now
	return bid, nil
}

func (a *Auction) End(now time.Time) error {
	if a.Status != StatusLive {
		return errors.New("auction is not live")
	}
	a.Status = StatusEnded
	a.EndTime = now
	a.UpdatedAt = now
	if a.CurrentBid != nil {
		if a.ReservePrice == nil || a.CurrentBid.Amount >= *a.ReservePrice {
			a.WinnerID = &a.CurrentBid.UserID
		}
	}
	return nil
}

type Bid struct {
	ID        uuid.UUID
	AuctionID uuid.UUID
	UserID    uuid.UUID
	User      *Bidder
	Amount    float64
	IsAutoBid bool
	IsWinning bool
	BidTime   time.Time
}

type Bidder struct {
	ID        uuid.UUID
	Email     string
	FirstName string
	LastName  string
	AvatarURL string
}

type AutoBid struct {
	ID           uuid.UUID
	AuctionID    uuid.UUID
	UserID       uuid.UUID
	MaxAmount    float64
	IsActive     bool
	CurrentBid   *float64
	BidIncrement float64
	LastBidTime  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (ab *AutoBid) ShouldBid(currentBid float64) (bool, float64) {
	if !ab.IsActive {
		return false, 0
	}
	if currentBid >= ab.MaxAmount {
		return false, 0
	}
	bidAmount := currentBid + ab.BidIncrement
	if bidAmount > ab.MaxAmount {
		bidAmount = ab.MaxAmount
	}
	return true, bidAmount
}

type Repository interface {
	Create(ctx context.Context, auction *Auction) error
	Update(ctx context.Context, auction *Auction) error
	Delete(ctx context.Context, id uuid.UUID) error
	AddBid(ctx context.Context, bid *Bid) error
	UpdateBidWinningStatus(ctx context.Context, auctionID, userID uuid.UUID, isWinning bool) error
	CreateAutoBid(ctx context.Context, autoBid *AutoBid) error
	UpdateAutoBid(ctx context.Context, autoBid *AutoBid) error
	GetByID(ctx context.Context, id uuid.UUID) (*Auction, error)
	GetWithBids(ctx context.Context, id uuid.UUID) (*Auction, error)
	GetLiveAuctions(ctx context.Context, limit, offset int) ([]*Auction, error)
	GetScheduledAuctions(ctx context.Context, limit, offset int) ([]*Auction, error)
	GetBidsByAuction(ctx context.Context, auctionID uuid.UUID, limit, offset int) ([]*Bid, error)
	GetActiveAutoBids(ctx context.Context, auctionID uuid.UUID) ([]*AutoBid, error)
	GetBidCount(ctx context.Context, auctionID uuid.UUID) (int, error)
}

type Cache interface {
	GetAuctionState(ctx context.Context, auctionID uuid.UUID) (*AuctionState, error)
	SetAuctionState(ctx context.Context, auctionID uuid.UUID, state *AuctionState, ttl time.Duration) error
	DeleteAuctionState(ctx context.Context, auctionID uuid.UUID) error
	IncrementViewerCount(ctx context.Context, auctionID uuid.UUID) (int, error)
	DecrementViewerCount(ctx context.Context, auctionID uuid.UUID) (int, error)
	GetViewerCount(ctx context.Context, auctionID uuid.UUID) (int, error)
}

type AuctionState struct {
	AuctionID   uuid.UUID `json:"auction_id"`
	CurrentBid  *Bid      `json:"current_bid,omitempty"`
	BidCount    int       `json:"bid_count"`
	Status      Status    `json:"status"`
	EndTime     time.Time `json:"end_time"`
	ViewerCount int       `json:"viewer_count"`
	LastUpdated time.Time `json:"last_updated"`
}

type EventBus interface {
	PublishBidPlaced(ctx context.Context, auctionID uuid.UUID, bid *Bid) error
	PublishAuctionStarted(ctx context.Context, auctionID uuid.UUID) error
	PublishAuctionEnded(ctx context.Context, auctionID uuid.UUID, winnerID *uuid.UUID) error
	PublishAuctionExtended(ctx context.Context, auctionID uuid.UUID, newEndTime time.Time) error
}