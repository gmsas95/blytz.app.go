package handlers

import (
	"net/http"

	"github.com/blytz/live/backend/internal/infrastructure/websocket"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuctionWSHandler handles WebSocket connections for auctions
type AuctionWSHandler struct {
	hub *websocket.Hub
}

// NewAuctionWSHandler creates a new auction WebSocket handler
func NewAuctionWSHandler(hub *websocket.Hub) *AuctionWSHandler {
	return &AuctionWSHandler{hub: hub}
}

// HandleWebSocket handles WebSocket upgrade for auction room
func (h *AuctionWSHandler) HandleWebSocket(c *gin.Context) {
	auctionID := c.Param("id")
	
	// Validate auction ID
	if _, err := uuid.Parse(auctionID); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid auction id"})
		return
	}
	
	// Get user ID from context (may be empty for anonymous viewers)
	userID, _ := c.Get("user_id")
	userIDStr := ""
	if userID != nil {
		userIDStr = userID.(uuid.UUID).String()
	}
	
	// Upgrade to WebSocket
	h.hub.HandleConnection(c.Writer, c.Request, auctionID, userIDStr)
}
