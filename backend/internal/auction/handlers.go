package auction

import (
	"net/http"
	"strconv"
	"time"

	"github.com/blytz.live.remake/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler provides auction HTTP handlers
type Handler struct {
	service *Service
}

// NewHandler creates a new auction handler
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// CreateAuctionRequest represents request body for creating an auction
type CreateAuctionRequest struct {
	ProductID    uuid.UUID  `json:"product_id" binding:"required"`
	Title        string     `json:"title" binding:"required,min=1,max=255"`
	Description  *string    `json:"description"`
	StartTime    time.Time  `json:"start_time" binding:"required"`
	EndTime      time.Time  `json:"end_time" binding:"required"`
	StartPrice   float64    `json:"start_price" binding:"required,gt=0"`
	ReservePrice *float64   `json:"reserve_price"`
	BuyNowPrice  *float64   `json:"buy_now_price"`
	AutoExtend   bool       `json:"auto_extend"`
	ExtendTime   int        `json:"extend_time"`
	IsFeatured   bool       `json:"is_featured"`
}

// PlaceBidRequest represents request body for placing a bid
type PlaceBidRequest struct {
	Amount    float64 `json:"amount" binding:"required,gt=0"`
	IsAutoBid bool    `json:"is_auto_bid"`
}

// AutoBidRequest represents request body for setting auto-bid
type AutoBidRequest struct {
	MaxAmount    float64 `json:"max_amount" binding:"required,gt=0"`
	BidIncrement float64 `json:"bid_increment"`
	IsActive     bool    `json:"is_active"`
}

// CreateAuction creates a new auction
func (h *Handler) CreateAuction(c *gin.Context) {
	var req CreateAuctionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	auction := &models.Auction{
		ProductID:    req.ProductID,
		SellerID:     userID.(uuid.UUID),
		Title:        req.Title,
		Description:  req.Description,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		StartPrice:   req.StartPrice,
		ReservePrice: req.ReservePrice,
		BuyNowPrice:  req.BuyNowPrice,
		Status:       "scheduled",
		AutoExtend:   req.AutoExtend,
		ExtendTime:   req.ExtendTime,
		IsFeatured:   req.IsFeatured,
		LiveKitRoom:  "auction-" + uuid.New().String(),
	}

	if err := h.service.CreateAuction(c.Request.Context(), auction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, auction)
}

// GetAuction gets an auction by ID
func (h *Handler) GetAuction(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auction ID"})
		return
	}

	auction, err := h.service.GetAuction(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Auction not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, auction)
}

// ListAuctions lists auctions with pagination
func (h *Handler) ListAuctions(c *gin.Context) {
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	auctions, total, err := h.service.ListAuctions(c.Request.Context(), status, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"auctions": auctions,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

// GetLiveAuctions gets currently live auctions
func (h *Handler) GetLiveAuctions(c *gin.Context) {
	auctions, err := h.service.GetLiveAuctions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"auctions": auctions})
}

// PlaceBid places a bid on an auction
func (h *Handler) PlaceBid(c *gin.Context) {
	auctionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auction ID"})
		return
	}

	var req PlaceBidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	bid, err := h.service.PlaceBid(c.Request.Context(), auctionID, userID.(uuid.UUID), req.Amount, req.IsAutoBid)
	if err != nil {
		if err.Error() == "auction is not live" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "auction has ended" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, bid)
}

// StartAuction starts an auction (seller/admin only)
func (h *Handler) StartAuction(c *gin.Context) {
	auctionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auction ID"})
		return
	}

	// Get user info from context
	role, _ := c.Get("role")

	// TODO: Verify user owns the auction or is admin
	// For now, just check if user is seller or admin
	if role != "seller" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	if err := h.service.StartAuction(c.Request.Context(), auctionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Auction started successfully"})
}

// EndAuction ends an auction (seller/admin only)
func (h *Handler) EndAuction(c *gin.Context) {
	auctionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auction ID"})
		return
	}

	// Get user info from context
	role, _ := c.Get("role")

	// TODO: Verify user owns the auction or is admin
	if role != "seller" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	if err := h.service.EndAuction(c.Request.Context(), auctionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Auction ended successfully"})
}

// SetAutoBid sets up automatic bidding
func (h *Handler) SetAutoBid(c *gin.Context) {
	auctionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auction ID"})
		return
	}

	var req AutoBidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	autoBid := &models.AutoBid{
		AuctionID:    auctionID,
		UserID:       userID.(uuid.UUID),
		MaxAmount:    req.MaxAmount,
		BidIncrement: req.BidIncrement,
		IsActive:     req.IsActive,
	}

	// Set default bid increment if not provided
	if req.BidIncrement == 0 {
		autoBid.BidIncrement = 5.0
	}

	if err := h.service.SetAutoBid(c.Request.Context(), autoBid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, autoBid)
}

// GetAuctionStats gets auction statistics
func (h *Handler) GetAuctionStats(c *gin.Context) {
	auctionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auction ID"})
		return
	}

	stats, err := h.service.GetAuctionStats(c.Request.Context(), auctionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// JoinAuction allows a user to join an auction audience
func (h *Handler) JoinAuction(c *gin.Context) {
	auctionID, err := uuid.Parse(c.Param("id"))
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

	if err := h.service.JoinAuction(c.Request.Context(), auctionID, userID.(uuid.UUID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Joined auction successfully"})
}

// LeaveAuction allows a user to leave an auction audience
func (h *Handler) LeaveAuction(c *gin.Context) {
	auctionID, err := uuid.Parse(c.Param("id"))
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

	if err := h.service.LeaveAuction(c.Request.Context(), auctionID, userID.(uuid.UUID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Left auction successfully"})
}

// GetAuctionBids gets all bids for an auction
func (h *Handler) GetAuctionBids(c *gin.Context) {
	auctionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auction ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	bids, total, err := h.service.GetAuctionBids(c.Request.Context(), auctionID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"bids":  bids,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}