package handlers

import (
	"net/http"
	"strconv"
	"time"

	auctionApp "github.com/blytz/live/backend/internal/application/auction"
	auctionDomain "github.com/blytz/live/backend/internal/domain/auction"
	appErrors "github.com/blytz/live/backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuctionHandler handles auction HTTP requests
type AuctionHandler struct {
	service *auctionApp.Service
}

// NewAuctionHandler creates a new auction handler
func NewAuctionHandler(service *auctionApp.Service) *AuctionHandler {
	return &AuctionHandler{service: service}
}

// CreateAuctionRequest represents auction creation request
type CreateAuctionRequest struct {
	ProductID    string  `json:"product_id" binding:"required"`
	Title        string  `json:"title" binding:"required"`
	Description  string  `json:"description"`
	StartTime    string  `json:"start_time" binding:"required"` // RFC3339
	EndTime      string  `json:"end_time" binding:"required"`
	StartPrice   float64 `json:"start_price" binding:"required,gt=0"`
	ReservePrice *float64 `json:"reserve_price"`
	BuyNowPrice  *float64 `json:"buy_now_price"`
	IsFeatured   bool    `json:"is_featured"`
}

// PlaceBidRequest represents bid placement request
type PlaceBidRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

// SetAutoBidRequest represents auto-bid setup request
type SetAutoBidRequest struct {
	MaxAmount    float64 `json:"max_amount" binding:"required,gt=0"`
	BidIncrement float64 `json:"bid_increment" binding:"required,gt=0"`
}

// AuctionResponse represents auction response
type AuctionResponse struct {
	ID           string     `json:"id"`
	ProductID    string     `json:"product_id"`
	SellerID     string     `json:"seller_id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	StartTime    time.Time  `json:"start_time"`
	EndTime      time.Time  `json:"end_time"`
	Status       string     `json:"status"`
	StartPrice   float64    `json:"start_price"`
	CurrentBid   *BidResponse `json:"current_bid,omitempty"`
	BidCount     int        `json:"bid_count"`
	LiveKitRoom  string     `json:"livekit_room"`
	IsFeatured   bool       `json:"is_featured"`
	CreatedAt    time.Time  `json:"created_at"`
}

// BidResponse represents bid response
type BidResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Amount    float64   `json:"amount"`
	IsAutoBid bool      `json:"is_auto_bid"`
	BidTime   time.Time `json:"bid_time"`
}

// CreateAuction creates a new auction
func (h *AuctionHandler) CreateAuction(c *gin.Context) {
	sellerID, _ := c.Get("user_id")

	var req CreateAuctionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, appErrors.New(appErrors.ErrValidation, err.Error()))
		return
	}

	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		respondError(c, appErrors.New(appErrors.ErrValidation, "invalid product_id"))
		return
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		respondError(c, appErrors.New(appErrors.ErrValidation, "invalid start_time format"))
		return
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		respondError(c, appErrors.New(appErrors.ErrValidation, "invalid end_time format"))
		return
	}

	sellerUUID, _ := uuid.Parse(sellerID.(string))

	a, err := h.service.CreateAuction(c.Request.Context(), &auctionApp.CreateAuctionRequest{
		ProductID:    productID,
		SellerID:     sellerUUID,
		Title:        req.Title,
		Description:  req.Description,
		StartTime:    startTime,
		EndTime:      endTime,
		StartPrice:   req.StartPrice,
		ReservePrice: req.ReservePrice,
		BuyNowPrice:  req.BuyNowPrice,
		IsFeatured:   req.IsFeatured,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	respondJSON(c, http.StatusCreated, toAuctionResponse(a))
}

// GetAuction gets auction by ID
func (h *AuctionHandler) GetAuction(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondError(c, appErrors.New(appErrors.ErrValidation, "invalid auction id"))
		return
	}

	a, err := h.service.GetAuction(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}

	respondJSON(c, http.StatusOK, toAuctionResponse(a))
}

// ListLiveAuctions lists live auctions
func (h *AuctionHandler) ListLiveAuctions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	auctions, err := h.service.ListLiveAuctions(c.Request.Context(), page, pageSize)
	if err != nil {
		respondError(c, err)
		return
	}

	responses := make([]*AuctionResponse, len(auctions))
	for i, a := range auctions {
		responses[i] = toAuctionResponse(a)
	}

	respondJSON(c, http.StatusOK, responses)
}

// PlaceBid places a bid on an auction
func (h *AuctionHandler) PlaceBid(c *gin.Context) {
	auctionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondError(c, appErrors.New(appErrors.ErrValidation, "invalid auction id"))
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	var req PlaceBidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, appErrors.New(appErrors.ErrValidation, err.Error()))
		return
	}

	bid, err := h.service.PlaceBid(c.Request.Context(), auctionID, userID, req.Amount, false)
	if err != nil {
		respondError(c, err)
		return
	}

	respondJSON(c, http.StatusCreated, toBidResponse(bid))
}

// StartAuction starts an auction
func (h *AuctionHandler) StartAuction(c *gin.Context) {
	auctionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondError(c, appErrors.New(appErrors.ErrValidation, "invalid auction id"))
		return
	}

	if err := h.service.StartAuction(c.Request.Context(), auctionID); err != nil {
		respondError(c, err)
		return
	}

	respondJSON(c, http.StatusOK, gin.H{"message": "auction started"})
}

// EndAuction ends an auction
func (h *AuctionHandler) EndAuction(c *gin.Context) {
	auctionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondError(c, appErrors.New(appErrors.ErrValidation, "invalid auction id"))
		return
	}

	if err := h.service.EndAuction(c.Request.Context(), auctionID); err != nil {
		respondError(c, err)
		return
	}

	respondJSON(c, http.StatusOK, gin.H{"message": "auction ended"})
}

// Helper functions
func toAuctionResponse(a *auctionDomain.Auction) *AuctionResponse {
	resp := &AuctionResponse{
		ID:          a.ID.String(),
		ProductID:   a.ProductID.String(),
		SellerID:    a.SellerID.String(),
		Title:       a.Title,
		Description: a.Description,
		StartTime:   a.StartTime,
		EndTime:     a.EndTime,
		Status:      string(a.Status),
		StartPrice:  a.StartPrice,
		BidCount:    a.BidCount,
		LiveKitRoom: a.LiveKitRoom,
		IsFeatured:  a.IsFeatured,
		CreatedAt:   a.CreatedAt,
	}

	if a.CurrentBid != nil {
		resp.CurrentBid = toBidResponse(a.CurrentBid)
	}

	return resp
}

func toBidResponse(b *auctionDomain.Bid) *BidResponse {
	return &BidResponse{
		ID:        b.ID.String(),
		UserID:    b.UserID.String(),
		Amount:    b.Amount,
		IsAutoBid: b.IsAutoBid,
		BidTime:   b.BidTime,
	}
}