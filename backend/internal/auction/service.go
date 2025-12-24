package auction

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/blytz.live.remake/backend/internal/models"
	"github.com/blytz.live.remake/backend/pkg/logging"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service provides auction business logic
type Service struct {
	db         *gorm.DB
	logger     *logging.Logger
	wsManager  *WebSocketManager
}

// NewService creates a new auction service
func NewService(db *gorm.DB) *Service {
	logger := logging.NewLogger()
	return &Service{
		db:     db,
		logger: logger,
	}
}

// SetWebSocketManager sets the WebSocket manager for real-time notifications
func (s *Service) SetWebSocketManager(wsManager *WebSocketManager) {
	s.wsManager = wsManager
}

// CreateAuction creates a new auction
func (s *Service) CreateAuction(ctx context.Context, auction *models.Auction) error {
	if auction.StartTime.Before(time.Now()) {
		return errors.New("start time cannot be in the past")
	}
	
	if auction.EndTime.Before(auction.StartTime) {
		return errors.New("end time must be after start time")
	}
	
	// Generate LiveKit room name
	if auction.LiveKitRoom == "" {
		auction.LiveKitRoom = fmt.Sprintf("auction-%s", auction.ID.String())
	}
	
	return s.db.WithContext(ctx).Create(auction).Error
}

// GetAuction retrieves an auction by ID
func (s *Service) GetAuction(ctx context.Context, id uuid.UUID) (*models.Auction, error) {
	var auction models.Auction
	err := s.db.WithContext(ctx).
		Preload("Product").
		Preload("Seller").
		Preload("Bids", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		Preload("Bids.User").
		First(&auction, "id = ?", id).Error
	
	if err != nil {
		return nil, err
	}
	
	// Update view count for active auctions
	if auction.Status == "live" {
		s.db.WithContext(ctx).Model(&auction).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))
	}
	
	return &auction, nil
}

// ListAuctions retrieves paginated auctions
func (s *Service) ListAuctions(ctx context.Context, status string, page, limit int) ([]models.Auction, int64, error) {
	var auctions []models.Auction
	var total int64
	
	query := s.db.WithContext(ctx).Model(&models.Auction{}).Preload("Product").Preload("Seller")
	
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Get paginated results
	offset := (page - 1) * limit
	err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&auctions).Error
	
	return auctions, total, err
}

// PlaceBid places a bid on an auction
func (s *Service) PlaceBid(ctx context.Context, auctionID, userID uuid.UUID, amount float64, isAutoBid bool) (*models.Bid, error) {
	// Start transaction
	tx := s.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	
	// Get auction and lock it
	var auction models.Auction
	if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&auction, "id = ?", auctionID).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	
	// Validate auction status
	if auction.Status != "live" {
		tx.Rollback()
		return nil, errors.New("auction is not live")
	}
	
	// Check if auction has ended
	if time.Now().After(auction.EndTime) {
		// Update auction status
		tx.Model(&auction).Update("status", "ended")
		tx.Rollback()
		return nil, errors.New("auction has ended")
	}
	
	// Validate bid amount
	minimumBid := auction.StartPrice
	if auction.CurrentBid != nil {
		minimumBid = *auction.CurrentBid + 1.0 // Minimum increment of $1
	}
	
	if amount < minimumBid {
		tx.Rollback()
		return nil, fmt.Errorf("bid amount must be at least $%.2f", minimumBid)
	}
	
	// Check if user is trying to outbid themselves
	if auction.CurrentBid != nil {
		var lastBid models.Bid
		if err := tx.Where("auction_id = ? AND user_id = ?", auctionID, userID).
			Order("created_at DESC").First(&lastBid).Error; err == nil {
			if lastBid.Amount >= amount {
				tx.Rollback()
				return nil, errors.New("cannot outbid yourself")
			}
		}
	}
	
	// Create bid
	bid := &models.Bid{
		AuctionID: auctionID,
		UserID:    userID,
		Amount:    amount,
		IsAutoBid: isAutoBid,
		BidTime:   time.Now(),
		IsWinning: true,
	}
	
	if err := tx.Create(bid).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	
	// Update previous winning bids
	if auction.CurrentBid != nil {
		tx.Model(&models.Bid{}).
			Where("auction_id = ? AND user_id != ? AND is_winning = ?", auctionID, userID, true).
			Update("is_winning", false)
	}
	
	// Update auction
	updates := map[string]interface{}{
		"current_bid": amount,
		"bid_count":   auction.BidCount + 1,
	}
	
	// Auto-extend auction if bid is in last 5 minutes
	if auction.AutoExtend && time.Until(auction.EndTime) <= 5*time.Minute {
		newEndTime := auction.EndTime.Add(time.Duration(auction.ExtendTime) * time.Second)
		updates["end_time"] = newEndTime
	}
	
	if err := tx.Model(&auction).Updates(updates).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	
	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	
	s.logger.Info("Bid placed successfully", map[string]interface{}{
		"auction_id": auctionID,
		"user_id":    userID,
		"amount":     amount,
		"is_auto_bid": isAutoBid,
	})
	
	// Notify WebSocket clients about the new bid
	if s.wsManager != nil {
		s.wsManager.NotifyBidPlaced(auctionID, bid)
	}
	
	// Process auto-bids
	go s.processAutoBids(context.Background(), auctionID, amount)
	
	return bid, nil
}

// processAutoBids processes automatic bids for an auction
func (s *Service) processAutoBids(ctx context.Context, auctionID uuid.UUID, currentBid float64) {
	// Get all active auto-bids for this auction
	var autoBids []models.AutoBid
	err := s.db.WithContext(ctx).
		Preload("User").
		Where("auction_id = ? AND is_active = ? AND max_amount > ?", auctionID, true, currentBid).
		Order("created_at ASC").
		Find(&autoBids).Error
	
	if err != nil || len(autoBids) == 0 {
		return
	}
	
	// Place bids for auto-bidders
	for _, autoBid := range autoBids {
		if autoBid.UserID == uuid.Nil { // Skip if user ID is invalid
			continue
		}
		
		// Calculate bid amount
		bidAmount := currentBid + autoBid.BidIncrement
		
		// Don't exceed max amount
		if bidAmount > autoBid.MaxAmount {
			bidAmount = autoBid.MaxAmount
		}
		
		// Place the bid
		_, err := s.PlaceBid(ctx, auctionID, autoBid.UserID, bidAmount, true)
		if err != nil {
			s.logger.Error("Failed to place auto-bid", map[string]interface{}{
				"auto_bid_id": autoBid.ID,
				"error":       err.Error(),
			})
			continue
		}
		
		// Update auto-bid record
		s.db.WithContext(ctx).Model(&autoBid).Updates(map[string]interface{}{
			"current_bid":  bidAmount,
			"last_bid_time": time.Now(),
		})
		
		currentBid = bidAmount
	}
}

// StartAuction starts an auction
func (s *Service) StartAuction(ctx context.Context, auctionID uuid.UUID) error {
	return s.db.WithContext(ctx).Model(&models.Auction{}).
		Where("id = ? AND status = ?", auctionID, "scheduled").
		Updates(map[string]interface{}{
			"status":     "live",
			"start_time": time.Now(),
		}).Error
}

// EndAuction ends an auction and determines the winner
func (s *Service) EndAuction(ctx context.Context, auctionID uuid.UUID) error {
	// Get auction with bids
	var auction models.Auction
	err := s.db.WithContext(ctx).
		Preload("Bids", func(db *gorm.DB) *gorm.DB {
			return db.Order("amount DESC, created_at ASC")
		}).
		First(&auction, "id = ?", auctionID).Error
	
	if err != nil {
		return err
	}
	
	// Determine winner
	var winnerID *uuid.UUID
	if len(auction.Bids) > 0 {
		highestBid := auction.Bids[0]
		
		// Check if reserve price is met
		if auction.ReservePrice != nil && highestBid.Amount < *auction.ReservePrice {
			// Reserve not met, no winner
		} else {
			winnerID = &highestBid.UserID
		}
	}
	
	// Update auction
	updates := map[string]interface{}{
		"status":   "ended",
		"end_time": time.Now(),
	}
	
	if winnerID != nil {
		updates["winner_id"] = winnerID
	}
	
	return s.db.WithContext(ctx).Model(&auction).Updates(updates).Error
}

// SetAutoBid sets up automatic bidding for a user
func (s *Service) SetAutoBid(ctx context.Context, autoBid *models.AutoBid) error {
	// Check if auto-bid already exists
	var existing models.AutoBid
	err := s.db.WithContext(ctx).
		Where("auction_id = ? AND user_id = ?", autoBid.AuctionID, autoBid.UserID).
		First(&existing).Error
	
	if err == nil {
		// Update existing
		return s.db.WithContext(ctx).Model(&existing).Updates(autoBid).Error
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new
		return s.db.WithContext(ctx).Create(autoBid).Error
	} else {
		return err
	}
}

// GetAuctionStats gets statistics for an auction
func (s *Service) GetAuctionStats(ctx context.Context, auctionID uuid.UUID) (*models.AuctionStats, error) {
	var stats models.AuctionStats
	err := s.db.WithContext(ctx).
		Preload("Auction").
		First(&stats, "auction_id = ?", auctionID).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create stats if not exists
			var auction models.Auction
			if err := s.db.WithContext(ctx).First(&auction, "id = ?", auctionID).Error; err != nil {
				return nil, err
			}
			
			stats = models.AuctionStats{
				AuctionID: auctionID,
				TotalBids: auction.BidCount,
			}
			
			// Count unique bidders
			var uniqueBidders int64
			s.db.WithContext(ctx).Model(&models.Bid{}).
				Where("auction_id = ?", auctionID).
				Distinct("user_id").
				Count(&uniqueBidders)
			stats.TotalBidders = int(uniqueBidders)
			
			err = s.db.WithContext(ctx).Create(&stats).Error
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	
	return &stats, nil
}

// GetLiveAuctions gets currently live auctions
func (s *Service) GetLiveAuctions(ctx context.Context) ([]models.Auction, error) {
	var auctions []models.Auction
	err := s.db.WithContext(ctx).
		Preload("Product").
		Preload("Seller").
		Where("status = ? AND end_time > ?", "live", time.Now()).
		Order("end_time ASC").
		Find(&auctions).Error
	
	return auctions, err
}

// JoinAuction allows a user to join an auction audience
func (s *Service) JoinAuction(ctx context.Context, auctionID, userID uuid.UUID) error {
	// Check if already watching
	var watch models.AuctionWatch
	err := s.db.WithContext(ctx).
		Where("auction_id = ? AND user_id = ?", auctionID, userID).
		First(&watch).Error
	
	if err == nil {
		// Update existing
		watch.IsActive = true
		watch.LastActive = time.Now()
		return s.db.WithContext(ctx).Save(&watch).Error
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new
		watch = models.AuctionWatch{
			AuctionID: auctionID,
			UserID:    userID,
			IsActive:  true,
			JoinedAt:  time.Now(),
			LastActive: time.Now(),
			NotificationSettings: `{"outbid": true, "ending_soon": true}`,
		}
		return s.db.WithContext(ctx).Create(&watch).Error
	} else {
		return err
	}
}

// LeaveAuction removes a user from auction audience
func (s *Service) LeaveAuction(ctx context.Context, auctionID, userID uuid.UUID) error {
	return s.db.WithContext(ctx).Model(&models.AuctionWatch{}).
		Where("auction_id = ? AND user_id = ?", auctionID, userID).
		Updates(map[string]interface{}{
			"is_active":  false,
			"last_active": time.Now(),
		}).Error
}

// GetAuctionBids gets all bids for an auction
func (s *Service) GetAuctionBids(ctx context.Context, auctionID uuid.UUID, page, limit int) ([]models.Bid, int64, error) {
	var bids []models.Bid
	var total int64
	
	query := s.db.WithContext(ctx).Model(&models.Bid{}).Preload("User").Where("auction_id = ?", auctionID)
	
	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Get paginated results
	offset := (page - 1) * limit
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&bids).Error
	
	return bids, total, err
}