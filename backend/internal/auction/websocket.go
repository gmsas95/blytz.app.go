package auction

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/blytz.live.remake/backend/internal/models"
	"github.com/blytz.live.remake/backend/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

// WebSocketManager manages WebSocket connections for auctions
type WebSocketManager struct {
	connections map[string]map[*websocket.Conn]bool // auction_id -> connections
	mutex       sync.RWMutex
	logger      *logging.Logger
	db          *gorm.DB
	service     *Service
}

// WebSocketMessage represents a WebSocket message
type WebSocketMessage struct {
	Type      string      `json:"type"` // bid, auction_update, chat, user_count
	AuctionID string      `json:"auction_id"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// NewWebSocketManager creates a new WebSocket manager
func NewWebSocketManager(db *gorm.DB, service *Service) *WebSocketManager {
	logger := logging.NewLogger()
	return &WebSocketManager{
		connections: make(map[string]map[*websocket.Conn]bool),
		logger:      logger,
		db:          db,
		service:     service,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// In production, implement proper origin checking
		return true
	},
}

// HandleWebSocket handles WebSocket connections for auctions
func (wsm *WebSocketManager) HandleWebSocket(c *gin.Context) {
	// Get auction ID from query params
	auctionIDStr := c.Query("auction_id")
	if auctionIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "auction_id is required"})
		return
	}

	// Validate auction ID
	auctionID, err := uuid.Parse(auctionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auction ID"})
		return
	}

	// Get user ID from context (if authenticated)
	var userID *uuid.UUID
	if uid, exists := c.Get("user_id"); exists {
		if id, ok := uid.(uuid.UUID); ok {
			userID = &id
		}
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade WebSocket connection: %v", err)
		return
	}
	defer conn.Close()

	// Add connection to the auction room
	wsm.addConnection(auctionID.String(), conn)
	defer wsm.removeConnection(auctionID.String(), conn)

	// Send initial auction state
	wsm.sendInitialData(conn, auctionID, userID)

	// Handle incoming messages
	go wsm.handleMessages(conn, auctionID, userID)

	// Send periodic updates
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Send user count and auction status
			wsm.sendStatusUpdate(auctionID.String())
		default:
			// Handle WebSocket ping/pong
			if _, _, err := conn.ReadMessage(); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket error: %v", err)
				}
				return
			}
		}
	}
}

// addConnection adds a WebSocket connection to an auction room
func (wsm *WebSocketManager) addConnection(auctionID string, conn *websocket.Conn) {
	wsm.mutex.Lock()
	defer wsm.mutex.Unlock()

	if wsm.connections[auctionID] == nil {
		wsm.connections[auctionID] = make(map[*websocket.Conn]bool)
	}
	wsm.connections[auctionID][conn] = true

	wsm.logger.Info("User joined auction room", map[string]interface{}{
		"auction_id": auctionID,
		"connections": len(wsm.connections[auctionID]),
	})
}

// removeConnection removes a WebSocket connection from an auction room
func (wsm *WebSocketManager) removeConnection(auctionID string, conn *websocket.Conn) {
	wsm.mutex.Lock()
	defer wsm.mutex.Unlock()

	if wsm.connections[auctionID] != nil {
		delete(wsm.connections[auctionID], conn)
		if len(wsm.connections[auctionID]) == 0 {
			delete(wsm.connections, auctionID)
		}
	}

	wsm.logger.Info("User left auction room", map[string]interface{}{
		"auction_id": auctionID,
		"remaining":  len(wsm.connections[auctionID]),
	})
}

// broadcastToAuction broadcasts a message to all connections in an auction room
func (wsm *WebSocketManager) broadcastToAuction(auctionID string, message WebSocketMessage) {
	wsm.mutex.RLock()
	defer wsm.mutex.RUnlock()

	if wsm.connections[auctionID] == nil {
		return
	}

	for conn := range wsm.connections[auctionID] {
		if err := conn.WriteJSON(message); err != nil {
			log.Printf("Error sending WebSocket message: %v", err)
			// Connection will be cleaned up on next error
		}
	}
}

// NotifyBidPlaced notifies all users in an auction about a new bid
func (wsm *WebSocketManager) NotifyBidPlaced(auctionID uuid.UUID, bid *models.Bid) {
	// Get bid with user info
	bidWithUser := *bid
	if err := wsm.db.Preload("User").First(&bidWithUser, "id = ?", bid.ID).Error; err != nil {
		wsm.logger.Error("Failed to load bid user info", map[string]interface{}{
			"bid_id": bid.ID,
			"error":  err.Error(),
		})
		return
	}

	message := WebSocketMessage{
		Type:      "bid",
		AuctionID: auctionID.String(),
		Data:      bidWithUser,
		Timestamp: time.Now(),
	}

	wsm.broadcastToAuction(auctionID.String(), message)
}

// NotifyAuctionUpdate notifies users about auction status changes
func (wsm *WebSocketManager) NotifyAuctionUpdate(auctionID uuid.UUID, auction *models.Auction) {
	message := WebSocketMessage{
		Type:      "auction_update",
		AuctionID: auctionID.String(),
		Data:      auction,
		Timestamp: time.Now(),
	}

	wsm.broadcastToAuction(auctionID.String(), message)
}

// NotifyChatMessage broadcasts a chat message to all users in an auction
func (wsm *WebSocketManager) NotifyChatMessage(auctionID uuid.UUID, message *models.ChatMessage) {
	wsMessage := WebSocketMessage{
		Type:      "chat",
		AuctionID: auctionID.String(),
		Data:      message,
		Timestamp: time.Now(),
	}

	wsm.broadcastToAuction(auctionID.String(), wsMessage)
}

// sendInitialData sends initial auction data to a newly connected user
func (wsm *WebSocketManager) sendInitialData(conn *websocket.Conn, auctionID uuid.UUID, userID *uuid.UUID) {
	// Get auction details
	auction, err := wsm.service.GetAuction(context.Background(), auctionID)
	if err != nil {
		log.Printf("Failed to get auction for initial data: %v", err)
		return
	}

	// Send auction data
	initialMessage := WebSocketMessage{
		Type:      "initial_data",
		AuctionID: auctionID.String(),
		Data: gin.H{
			"auction":     auction,
			"user_count":  wsm.getConnectionCount(auctionID.String()),
			"current_bid": auction.CurrentBid,
		},
		Timestamp: time.Now(),
	}

	if err := conn.WriteJSON(initialMessage); err != nil {
		log.Printf("Failed to send initial data: %v", err)
	}
}

// handleMessages handles incoming WebSocket messages
func (wsm *WebSocketManager) handleMessages(conn *websocket.Conn, auctionID uuid.UUID, userID *uuid.UUID) {
	for {
		var message WebSocketMessage
		if err := conn.ReadJSON(&message); err != nil {
			if !websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}

		// Handle different message types
		switch message.Type {
		case "ping":
			// Respond with pong
			pongMessage := WebSocketMessage{
				Type:      "pong",
				AuctionID: auctionID.String(),
				Timestamp: time.Now(),
			}
			conn.WriteJSON(pongMessage)
			
		case "chat":
			// Handle chat message
			if userID != nil {
				wsm.handleChatMessage(auctionID, *userID, message.Data)
			}
		}
	}
}

// handleChatMessage processes incoming chat messages
func (wsm *WebSocketManager) handleChatMessage(auctionID uuid.UUID, userID uuid.UUID, data interface{}) {
	// Parse chat message data
	chatData, ok := data.(map[string]interface{})
	if !ok {
		return
	}

	messageText, ok := chatData["message"].(string)
	if !ok || messageText == "" {
		return
	}

	// Create chat message
	chatMessage := &models.ChatMessage{
		AuctionID:   auctionID,
		UserID:      userID,
		Message:     messageText,
		MessageType: "user",
		Timestamp:   time.Now(),
	}

	// Save to database
	if err := wsm.db.Create(chatMessage).Error; err != nil {
		wsm.logger.Error("Failed to save chat message", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Load user info
	if err := wsm.db.Preload("User").First(chatMessage, "id = ?", chatMessage.ID).Error; err != nil {
		return
	}

	// Broadcast to all users
	wsm.NotifyChatMessage(auctionID, chatMessage)
}

// sendStatusUpdate sends periodic status updates to auction room
func (wsm *WebSocketManager) sendStatusUpdate(auctionID string) {
	userCount := wsm.getConnectionCount(auctionID)
	
	// Get current auction info
	auctionUUID, err := uuid.Parse(auctionID)
	if err != nil {
		return
	}

	auction, err := wsm.service.GetAuction(context.Background(), auctionUUID)
	if err != nil {
		return
	}

	statusMessage := WebSocketMessage{
		Type:      "status_update",
		AuctionID: auctionID,
		Data: gin.H{
			"user_count":   userCount,
			"current_bid":  auction.CurrentBid,
			"bid_count":    auction.BidCount,
			"status":       auction.Status,
			"end_time":     auction.EndTime,
			"time_left":    time.Until(auction.EndTime),
		},
		Timestamp: time.Now(),
	}

	wsm.broadcastToAuction(auctionID, statusMessage)
}

// getConnectionCount returns the number of active connections for an auction
func (wsm *WebSocketManager) getConnectionCount(auctionID string) int {
	wsm.mutex.RLock()
	defer wsm.mutex.RUnlock()

	if wsm.connections[auctionID] == nil {
		return 0
	}
	return len(wsm.connections[auctionID])
}

// GetActiveConnections returns the number of active connections for all auctions
func (wsm *WebSocketManager) GetActiveConnections() map[string]int {
	wsm.mutex.RLock()
	defer wsm.mutex.RUnlock()

	result := make(map[string]int)
	for auctionID, connections := range wsm.connections {
		result[auctionID] = len(connections)
	}
	return result
}