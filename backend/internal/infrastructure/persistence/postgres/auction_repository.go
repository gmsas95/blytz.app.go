package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/blytz/live/backend/internal/domain/auction"
	appErrors "github.com/blytz/live/backend/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuctionRepository implements auction.Repository
type AuctionRepository struct {
	db *gorm.DB
}

// NewAuctionRepository creates a new auction repository
func NewAuctionRepository(db *gorm.DB) *AuctionRepository {
	return &AuctionRepository{db: db}
}

// Create creates a new auction
func (r *AuctionRepository) Create(ctx context.Context, a *auction.Auction) error {
	model := toAuctionModel(a)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	a.ID = model.ID
	a.CreatedAt = model.CreatedAt
	a.UpdatedAt = model.UpdatedAt
	return nil
}

// Update updates an auction
func (r *AuctionRepository) Update(ctx context.Context, a *auction.Auction) error {
	model := toAuctionModel(a)
	return r.db.WithContext(ctx).Save(model).Error
}

// Delete soft-deletes an auction
func (r *AuctionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&Auction{}, "id = ?", id).Error
}

// GetByID gets auction by ID
func (r *AuctionRepository) GetByID(ctx context.Context, id uuid.UUID) (*auction.Auction, error) {
	var model Auction
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.New(appErrors.ErrNotFound, "auction not found")
		}
		return nil, err
	}
	return toAuctionDomain(&model), nil
}

// GetWithBids gets auction with current bid
func (r *AuctionRepository) GetWithBids(ctx context.Context, id uuid.UUID) (*auction.Auction, error) {
	var model Auction
	if err := r.db.WithContext(ctx).Preload("CurrentBid").First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.New(appErrors.ErrNotFound, "auction not found")
		}
		return nil, err
	}
	return toAuctionDomain(&model), nil
}

// GetLiveAuctions gets currently live auctions
func (r *AuctionRepository) GetLiveAuctions(ctx context.Context, limit, offset int) ([]*auction.Auction, error) {
	var models []Auction
	err := r.db.WithContext(ctx).
		Where("status = ? AND end_time > ?", string(auction.StatusLive), time.Now()).
		Order("end_time ASC").
		Limit(limit).
		Offset(offset).
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	auctions := make([]*auction.Auction, len(models))
	for i, m := range models {
		auctions[i] = toAuctionDomain(&m)
	}
	return auctions, nil
}

// GetScheduledAuctions gets upcoming auctions
func (r *AuctionRepository) GetScheduledAuctions(ctx context.Context, limit, offset int) ([]*auction.Auction, error) {
	var models []Auction
	err := r.db.WithContext(ctx).
		Where("status = ? AND start_time > ?", string(auction.StatusScheduled), time.Now()).
		Order("start_time ASC").
		Limit(limit).
		Offset(offset).
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	auctions := make([]*auction.Auction, len(models))
	for i, m := range models {
		auctions[i] = toAuctionDomain(&m)
	}
	return auctions, nil
}

// AddBid creates a new bid and updates auction state
func (r *AuctionRepository) AddBid(ctx context.Context, bid *auction.Bid) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create bid
		bidModel := toBidModel(bid)
		if err := tx.Create(bidModel).Error; err != nil {
			return fmt.Errorf("failed to create bid: %w", err)
		}
		bid.ID = bidModel.ID

		// Update previous winning bids
		if err := tx.Model(&Bid{}).
			Where("auction_id = ? AND user_id != ? AND is_winning = ?", 
				bid.AuctionID, bid.UserID, true).
			Update("is_winning", false).Error; err != nil {
			return fmt.Errorf("failed to update previous bids: %w", err)
		}

		return nil
	})
}

// UpdateBidWinningStatus updates winning status
func (r *AuctionRepository) UpdateBidWinningStatus(ctx context.Context, auctionID, userID uuid.UUID, isWinning bool) error {
	return r.db.WithContext(ctx).Model(&Bid{}).
		Where("auction_id = ? AND user_id = ?", auctionID, userID).
		Update("is_winning", isWinning).Error
}

// GetBidsByAuction gets bids for an auction
func (r *AuctionRepository) GetBidsByAuction(ctx context.Context, auctionID uuid.UUID, limit, offset int) ([]*auction.Bid, error) {
	var models []Bid
	err := r.db.WithContext(ctx).
		Where("auction_id = ?", auctionID).
		Order("bid_time DESC").
		Limit(limit).
		Offset(offset).
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	bids := make([]*auction.Bid, len(models))
	for i, m := range models {
		bids[i] = toBidDomain(&m)
	}
	return bids, nil
}

// GetActiveAutoBids gets active auto-bids for auction
func (r *AuctionRepository) GetActiveAutoBids(ctx context.Context, auctionID uuid.UUID) ([]*auction.AutoBid, error) {
	var models []AutoBid
	err := r.db.WithContext(ctx).
		Where("auction_id = ? AND is_active = ?", auctionID, true).
		Order("created_at ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	autoBids := make([]*auction.AutoBid, len(models))
	for i, m := range models {
		autoBids[i] = toAutoBidDomain(&m)
	}
	return autoBids, nil
}

// CreateAutoBid creates an auto-bid
func (r *AuctionRepository) CreateAutoBid(ctx context.Context, autoBid *auction.AutoBid) error {
	model := toAutoBidModel(autoBid)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	autoBid.ID = model.ID
	autoBid.CreatedAt = model.CreatedAt
	return nil
}

// UpdateAutoBid updates an auto-bid
func (r *AuctionRepository) UpdateAutoBid(ctx context.Context, autoBid *auction.AutoBid) error {
	model := toAutoBidModel(autoBid)
	return r.db.WithContext(ctx).Save(model).Error
}

// GetBidCount gets total bid count for auction
func (r *AuctionRepository) GetBidCount(ctx context.Context, auctionID uuid.UUID) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&Bid{}).Where("auction_id = ?", auctionID).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

// Helper functions
func toAuctionModel(a *auction.Auction) *Auction {
	var currentBidID *uuid.UUID
	if a.CurrentBid != nil {
		currentBidID = &a.CurrentBid.ID
	}

	return &Auction{
		BaseModel: BaseModel{
			ID:        a.ID,
			CreatedAt: a.CreatedAt,
			UpdatedAt: a.UpdatedAt,
		},
		ProductID:    a.ProductID,
		SellerID:     a.SellerID,
		Title:        a.Title,
		Description:  a.Description,
		StartTime:    a.StartTime,
		EndTime:      a.EndTime,
		Status:       string(a.Status),
		StartPrice:   a.StartPrice,
		ReservePrice: a.ReservePrice,
		BuyNowPrice:  a.BuyNowPrice,
		CurrentBidID: currentBidID,
		BidCount:     a.BidCount,
		WinnerID:     a.WinnerID,
		LiveKitRoom:  a.LiveKitRoom,
		StreamKey:    a.StreamKey,
		AutoExtend:   a.AutoExtend,
		ExtendTime:   int(a.ExtendTime.Seconds()),
		IsFeatured:   a.IsFeatured,
	}
}

func toAuctionDomain(m *Auction) *auction.Auction {
	a := &auction.Auction{
		ID:           m.ID,
		ProductID:    m.ProductID,
		SellerID:     m.SellerID,
		Title:        m.Title,
		Description:  m.Description,
		StartTime:    m.StartTime,
		EndTime:      m.EndTime,
		Status:       auction.Status(m.Status),
		StartPrice:   m.StartPrice,
		ReservePrice: m.ReservePrice,
		BuyNowPrice:  m.BuyNowPrice,
		BidCount:     m.BidCount,
		WinnerID:     m.WinnerID,
		LiveKitRoom:  m.LiveKitRoom,
		StreamKey:    m.StreamKey,
		AutoExtend:   m.AutoExtend,
		ExtendTime:   time.Duration(m.ExtendTime) * time.Second,
		IsFeatured:   m.IsFeatured,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}

	// Current bid would be populated via preload or separate query
	return a
}

func toBidModel(b *auction.Bid) *Bid {
	return &Bid{
		BaseModel: BaseModel{
			ID: b.ID,
		},
		AuctionID: b.AuctionID,
		UserID:    b.UserID,
		Amount:    b.Amount,
		IsAutoBid: b.IsAutoBid,
		IsWinning: b.IsWinning,
		BidTime:   b.BidTime,
	}
}

func toBidDomain(m *Bid) *auction.Bid {
	return &auction.Bid{
		ID:        m.ID,
		AuctionID: m.AuctionID,
		UserID:    m.UserID,
		Amount:    m.Amount,
		IsAutoBid: m.IsAutoBid,
		IsWinning: m.IsWinning,
		BidTime:   m.BidTime,
	}
}

func toAutoBidModel(a *auction.AutoBid) *AutoBid {
	return &AutoBid{
		BaseModel: BaseModel{
			ID:        a.ID,
			CreatedAt: a.CreatedAt,
			UpdatedAt: a.UpdatedAt,
		},
		AuctionID:    a.AuctionID,
		UserID:       a.UserID,
		MaxAmount:    a.MaxAmount,
		IsActive:     a.IsActive,
		CurrentBid:   a.CurrentBid,
		BidIncrement: a.BidIncrement,
		LastBidTime:  a.LastBidTime,
	}
}

func toAutoBidDomain(m *AutoBid) *auction.AutoBid {
	return &auction.AutoBid{
		ID:           m.ID,
		AuctionID:    m.AuctionID,
		UserID:       m.UserID,
		MaxAmount:    m.MaxAmount,
		IsActive:     m.IsActive,
		CurrentBid:   m.CurrentBid,
		BidIncrement: m.BidIncrement,
		LastBidTime:  m.LastBidTime,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}