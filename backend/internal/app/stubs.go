package app

import (
	"context"
	"net/http"

	"github.com/blytz/live/backend/internal/infrastructure/messaging/redis"
	redisclient "github.com/redis/go-redis/v9"
)

// HTTPServer stub - will be implemented in infrastructure/http
type HTTPServer struct {
	port     string
	handlers interface{} // Will be http.Handlers
}

func NewHTTPServer(port string, handlers interface{}) *HTTPServer {
	return &HTTPServer{port: port, handlers: handlers}
}

func (s *HTTPServer) Start(ctx context.Context) error {
	// TODO: Implement HTTP server
	<-ctx.Done()
	return nil
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return nil
}

// WebSocketHub stub - will be implemented in infrastructure/websocket
type WebSocketHub struct {
	redisClient *redisclient.Client
	eventBus    *redis.EventBus
}

func NewWebSocketHub(redisClient *redisclient.Client, eventBus *redis.EventBus) *WebSocketHub {
	return &WebSocketHub{
		redisClient: redisClient,
		eventBus:    eventBus,
	}
}

func (h *WebSocketHub) Start(ctx context.Context) error {
	// TODO: Implement WebSocket hub
	<-ctx.Done()
	return nil
}

func (h *WebSocketHub) Shutdown(ctx context.Context) error {
	return nil
}

// Handler interface stub
type Handler interface {
	RegisterRoutes(mux *http.ServeMux)
}