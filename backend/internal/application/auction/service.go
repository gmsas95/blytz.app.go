package auction

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/blytz/live/backend/internal/domain/auction"
	appErrors "github.com/blytz/live/backend/pkg/errors"
	"github.com/google/uuid"
)

// Service provides auction use cases
type Service struct {
	repo      auction.Repository
	cache     auction.Cache
	eventBus  auction.EventBus
}

// NewService creates a new auction service
func NewService(repo auction.Repository, cache auction.Cache, eventBus auction.EventBus) *Service {
	return &Service{
		repo:     repo,
		cache:    cache,
		eventBus: eventBus,
	}
}

// CreateAuction creates a new auction
func (s *Service) CreateAuction(ctx context.Context, req *CreateAuctionRequest) (*auction.Auction, error) {
	// Validate request
	if req.StartTime.Before(time.Now()) {
		return nil, appErrors.New(appErrors.ErrValidation, "start time cannot be in the past")
	}
	if req.EndTime.Before(req.StartTime) {
		return nil, appErrors.New(appErrors.ErrValidation, "end time must be after start time")
	}

	a := &auction.Auction{
		ID:           uuid.New(),
		ProductID:    req.ProductID,
		SellerID:     req.SellerID,
		Title:        req.Title,
		Description:  req.Description,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		Status:       auction.StatusScheduled,
		StartPrice:   req.StartPrice,
		ReservePrice: req.ReservePrice,
		BuyNowPrice:  req.BuyNowPrice,
		LiveKitRoom:  fmt.Sprintf("auction-%s", uuid.New().String()),
		AutoExtend:   true,
		ExtendTime:   300, // 5 minutes
		IsFeatured:   req.IsFeatured,
	}

	if err := s.repo.Create(ctx, a); err != nil {
		return nil, appErrors.Wrap(err, appErrors.ErrInternal, "failed to create auction")
	}

	// Cache auction state
	if s.cache != nil {
		state := &auction.AuctionState{
			AuctionID:   a.ID,
			Status:      a.Status,
			EndTime:     a.EndTime,
			BidCount:    0,
			LastUpdated: time.Now(),
		}
		s.cache.SetAuctionState(ctx, a.ID, state, time.Hour)
	}

	return a, nil
}

// GetAuction gets auction by ID
func (s *Service) GetAuction(ctx context.Context, id uuid.UUID) (*auction.Auction, error) {
	// Try cache first
	if s.cache != nil {
		state, err := s.cache.GetAuctionState(ctx, id)
		if err == nil && time.Since(state.LastUpdated) < 5*time.Second {
			// Return from cache if fresh
			a := &auction.Auction{
				ID:        id,
				Status:    state.Status,
				EndTime:   state.EndTime,
				BidCount:  state.BidCount,
				CurrentBid: state.CurrentBid,
			}
			return a, nil
		}
	}

	// Get from database
	a, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update cache
	if s.cache != nil {
		state := &auction.AuctionState{
			AuctionID:   a.ID,
			Status:      a.Status,
			EndTime:     a.EndTime,
			BidCount:    a.BidCount,
			CurrentBid:  a.CurrentBid,
			LastUpdated: time.Now(),
		}
		s.cache.SetAuctionState(ctx, id, state, time.Hour)
	}

	return a, nil
}

// PlaceBid places a bid on an auction
func (s *Service) PlaceBid(ctx context.Context, auctionID, userID uuid.UUID, amount float64, isAutoBid bool) (*auction.Bid, error) {
	// Get auction with lock
	a, err := s.repo.GetWithBids(ctx, auctionID)
	if err != nil {
		return nil, err
	}

	// Validate and place bid using domain logic
	bid, err := a.PlaceBid(userID, amount, time.Now())
	if err != nil {
		// Map domain errors to app errors
		switch err.Error() {
		case "auction is not live":
			return nil, appErrors.New(appErrors.ErrAuctionNotLive, err.Error())
		case "auction has ended":
			return nil, appErrors.New(appErrors.ErrAuctionEnded, err.Error())
		case "bid amount too low":
			return nil, appErrors.New(appErrors.ErrBidTooLow, err.Error())
		default:
			return nil, appErrors.New(appErrors.ErrInvalidBid, err.Error())
		}
	}
	bid.IsAutoBid = isAutoBid

	// Persist bid and update auction
	if err := s.repo.AddBid(ctx, bid); err != nil {
		return nil, appErrors.Wrap(err, appErrors.ErrInternal, "failed to save bid")
	}

	// Update auction state
	if err := s.repo.Update(ctx, a); err != nil {
		log.Printf("Failed to update auction state: %v", err)
	}

	// Publish event
	if s.eventBus != nil {
		if err := s.eventBus.PublishBidPlaced(ctx, auctionID, bid); err != nil {
			log.Printf("Failed to publish bid event: %v", err)
		}
	}

	// Update cache
	if s.cache != nil {
		state := &auction.AuctionState{
			AuctionID:   a.ID,
			CurrentBid:  bid,
			BidCount:    a.BidCount,
			Status:      a.Status,
			EndTime:     a.EndTime,
			LastUpdated: time.Now(),
		}
		s.cache.SetAuctionState(ctx, auctionID, state, time.Hour)
	}

	// Process auto-bids asynchronously
	go s.processAutoBids(context.Background(), auctionID, amount)

	return bid, nil
}

// StartAuction starts an auction
func (s *Service) StartAuction(ctx context.Context, auctionID uuid.UUID) error {
	a, err := s.repo.GetByID(ctx, auctionID)
	if err != nil {
		return err
	}

	if a.Status != auction.StatusScheduled {
		return appErrors.New(appErrors.ErrValidation, "auction is not scheduled")
	}

	a.Status = auction.StatusLive
	a.StartTime = time.Now()

	if err := s.repo.Update(ctx, a); err != nil {
		return appErrors.Wrap(err, appErrors.ErrInternal, "failed to start auction")
	}

	// Publish event
	if s.eventBus != nil {
		s.eventBus.PublishAuctionStarted(ctx, auctionID)
	}

	// Update cache
	if s.cache != nil {
		state := &auction.AuctionState{
			AuctionID:   a.ID,
			Status:      a.Status,
			EndTime:     a.EndTime,
			LastUpdated: time.Now(),
		}
		s.cache.SetAuctionState(ctx, auctionID, state, time.Hour)
	}

	return nil
}

// EndAuction ends an auction
func (s *Service) EndAuction(ctx context.Context, auctionID uuid.UUID) error {
	a, err := s.repo.GetWithBids(ctx, auctionID)
	if err != nil {
		return err
	}

	if err := a.End(time.Now()); err != nil {
		return err
	}

	if err := s.repo.Update(ctx, a); err != nil {
		return appErrors.Wrap(err, appErrors.ErrInternal, "failed to end auction")
	}

	// Publish event
	if s.eventBus != nil {
		s.eventBus.PublishAuctionEnded(ctx, auctionID, a.WinnerID)
	}

	// Update cache
	if s.cache != nil {
		s.cache.DeleteAuctionState(ctx, auctionID)
	}

	return nil
}

// ListLiveAuctions lists currently live auctions
func (s *Service) ListLiveAuctions(ctx context.Context, page, pageSize int) ([]*auction.Auction, error) {
	return s.repo.GetLiveAuctions(ctx, pageSize, (page-1)*pageSize)
}

// SetAutoBid sets up automatic bidding
func (s *Service) SetAutoBid(ctx context.Context, auctionID, userID uuid.UUID, maxAmount, increment float64) (*auction.AutoBid, error) {
	// Check if auto-bid exists
	autoBids, err := s.repo.GetActiveAutoBids(ctx, auctionID)
	if err != nil {
		return nil, err
	}

	for _, ab := range autoBids {
		if ab.UserID == userID {
			// Update existing
			ab.MaxAmount = maxAmount
			ab.BidIncrement = increment
			if err := s.repo.UpdateAutoBid(ctx, ab); err != nil {
				return nil, err
			}
			return ab, nil
		}
	}

	// Create new
	autoBid := &auction.AutoBid{
		ID:           uuid.New(),
		AuctionID:    auctionID,
		UserID:       userID,
		MaxAmount:    maxAmount,
		IsActive:     true,
		BidIncrement: increment,
	}

	if err := s.repo.CreateAutoBid(ctx, autoBid); err != nil {
		return nil, err
	}

	return autoBid, nil
}

// processAutoBids processes automatic bids
func (s *Service) processAutoBids(ctx context.Context, auctionID uuid.UUID, currentBid float64) {
	autoBids, err := s.repo.GetActiveAutoBids(ctx, auctionID)
	if err != nil || len(autoBids) == 0 {
		return
	}

	for _, autoBid := range autoBids {
		shouldBid, bidAmount := autoBid.ShouldBid(currentBid)
		if !shouldBid {
			continue
		}

		_, err := s.PlaceBid(ctx, auctionID, autoBid.UserID, bidAmount, true)
		if err != nil {
			log.Printf("Failed to place auto-bid: %v", err)
			continue
		}

		// Update auto-bid record
		autoBid.CurrentBid = &bidAmount
		now := time.Now()
		autoBid.LastBidTime = &now
		s.repo.UpdateAutoBid(ctx, autoBid)

		currentBid = bidAmount
	}
}

// Request/Response types

type CreateAuctionRequest struct {
	ProductID    uuid.UUID
	SellerID     uuid.UUID
	Title        string
	Description  string
	StartTime    time.Time
	EndTime      time.Time
	StartPrice   float64
	ReservePrice *float64
	BuyNowPrice  *float64
	IsFeatured   bool
}