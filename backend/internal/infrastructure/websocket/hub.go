package websocket

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	redisMessaging "github.com/blytz/live/backend/internal/infrastructure/messaging/redis"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

// Hub manages WebSocket connections across instances using Redis Pub/Sub
type Hub struct {
	// Local connections
	rooms map[string]*Room // auction_id -> Room
	mu    sync.RWMutex
	
	// Redis for cross-instance communication
	redisClient *redis.Client
	eventBus    *redisMessaging.EventBus
	subscription *redisMessaging.Subscription
	
	// WebSocket upgrader
	upgrader websocket.Upgrader
}

// Room represents an auction room with connected clients
type Room struct {
	AuctionID   string
	clients     map[*Client]bool
	mu          sync.RWMutex
	viewerCount int
}

// Client represents a WebSocket client
type Client struct {
	hub       *Hub
	conn      *websocket.Conn
	send      chan []byte
	room      *Room
	userID    string
	auctionID string
}

// Message represents a WebSocket message
type Message struct {
	Type      string                 `json:"type"`
	AuctionID string                 `json:"auction_id"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}

// NewHub creates a new WebSocket hub
func NewHub(redisClient *redis.Client, eventBus *redisMessaging.EventBus) *Hub {
	return &Hub{
		rooms:       make(map[string]*Room),
		redisClient: redisClient,
		eventBus:    eventBus,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Allow all origins in development
				return true
			},
		},
	}
}

// Start starts the hub and subscribes to Redis events
func (h *Hub) Start(ctx context.Context) error {
	// Subscribe to global events
	sub, err := h.eventBus.SubscribeGlobal(ctx)
	if err != nil {
		return err
	}
	h.subscription = sub
	
	// Process incoming events from Redis
	go h.processEvents(ctx)
	
	<-ctx.Done()
	return nil
}

// Shutdown gracefully shuts down the hub
func (h *Hub) Shutdown(ctx context.Context) error {
	if h.subscription != nil {
		return h.subscription.Close()
	}
	return nil
}

// HandleConnection handles WebSocket upgrade and connection
func (h *Hub) HandleConnection(w http.ResponseWriter, r *http.Request, auctionID, userID string) {
	// Upgrade connection
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	
	// Get or create room
	room := h.getOrCreateRoom(auctionID)
	
	// Create client
	client := &Client{
		hub:       h,
		conn:      conn,
		send:      make(chan []byte, 256),
		room:      room,
		userID:    userID,
		auctionID: auctionID,
	}
	
	// Register client
	room.addClient(client)
	h.broadcastViewerCount(auctionID)
	
	// Start goroutines
	go client.writePump()
	go client.readPump()
}

// processEvents processes events from Redis Pub/Sub
func (h *Hub) processEvents(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-h.subscription.Channel():
			var event redisMessaging.Event
			if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
				log.Printf("Failed to unmarshal event: %v", err)
				continue
			}
			
			h.handleEvent(event)
		}
	}
}

// handleEvent handles a single event
func (h *Hub) handleEvent(event redisMessaging.Event) {
	switch event.Type {
	case redisMessaging.EventBidPlaced:
		h.broadcastToRoom(event.AuctionID, Message{
			Type:      "bid",
			AuctionID: event.AuctionID,
			Data:      event.Payload,
			Timestamp: event.Timestamp,
		})
		
	case redisMessaging.EventAuctionStarted:
		h.broadcastToRoom(event.AuctionID, Message{
			Type:      "auction_started",
			AuctionID: event.AuctionID,
			Data:      event.Payload,
			Timestamp: event.Timestamp,
		})
		
	case redisMessaging.EventAuctionEnded:
		h.broadcastToRoom(event.AuctionID, Message{
			Type:      "auction_ended",
			AuctionID: event.AuctionID,
			Data:      event.Payload,
			Timestamp: event.Timestamp,
		})
		
	case redisMessaging.EventAuctionExtended:
		h.broadcastToRoom(event.AuctionID, Message{
			Type:      "auction_extended",
			AuctionID: event.AuctionID,
			Data:      event.Payload,
			Timestamp: event.Timestamp,
		})
	}
}

// broadcastToRoom broadcasts a message to all clients in a room
func (h *Hub) broadcastToRoom(auctionID string, msg Message) {
	room := h.getRoom(auctionID)
	if room == nil {
		return
	}
	
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return
	}
	
	room.broadcast(data)
}

// broadcastViewerCount broadcasts viewer count update
func (h *Hub) broadcastViewerCount(auctionID string) {
	room := h.getRoom(auctionID)
	if room == nil {
		return
	}
	
	count := room.getViewerCount()
	msg := Message{
		Type:      "viewer_count",
		AuctionID: auctionID,
		Data: map[string]interface{}{
			"count": count,
		},
		Timestamp: time.Now(),
	}
	
	data, _ := json.Marshal(msg)
	room.broadcast(data)
}

// getRoom gets a room by auction ID
func (h *Hub) getRoom(auctionID string) *Room {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.rooms[auctionID]
}

// getOrCreateRoom gets or creates a room
func (h *Hub) getOrCreateRoom(auctionID string) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	if room, ok := h.rooms[auctionID]; ok {
		return room
	}
	
	room := &Room{
		AuctionID: auctionID,
		clients:   make(map[*Client]bool),
	}
	h.rooms[auctionID] = room
	return room
}

// Room methods

func (r *Room) addClient(c *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[c] = true
	r.viewerCount++
}

func (r *Room) removeClient(c *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.clients[c]; ok {
		delete(r.clients, c)
		r.viewerCount--
		close(c.send)
	}
}

func (r *Room) getViewerCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.viewerCount
}

func (r *Room) broadcast(message []byte) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	for client := range r.clients {
		select {
		case client.send <- message:
		default:
			// Client send buffer full, skip
		}
	}
}

// Client methods

func (c *Client) readPump() {
	defer func() {
		c.room.removeClient(c)
		c.conn.Close()
		c.hub.broadcastViewerCount(c.auctionID)
	}()
	
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})
	
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			
			c.conn.WriteMessage(websocket.TextMessage, message)
			
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}