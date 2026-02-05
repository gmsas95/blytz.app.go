package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	userDomain "github.com/blytz/live/backend/internal/domain/user"
	"github.com/blytz/live/backend/internal/infrastructure/cache/redis"
	"github.com/blytz/live/backend/internal/interfaces/http/handlers"
	"github.com/blytz/live/backend/internal/interfaces/middleware"
	"github.com/gin-gonic/gin"
)

// Server represents the HTTP server
type Server struct {
	router  *gin.Engine
	server  *http.Server
	handlers *Handlers
}

// Handlers holds all HTTP handlers
type Handlers struct {
	Auth      *handlers.AuthHandler
	Auction   *handlers.AuctionHandler
	Product   *handlers.ProductHandler
	Category  *handlers.CategoryHandler
	AuctionWS *handlers.AuctionWSHandler
	Upload    *handlers.UploadHandler
}

// NewServer creates a new HTTP server
func NewServer(port string, h *Handlers, tokenManager user.TokenManager, redisClient *redis.Client) *Server {
	gin.SetMode(gin.ReleaseMode)
	
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())
	router.Use(loggerMiddleware())

	s := &Server{
		router:   router,
		handlers: h,
	}

	s.setupRoutes(tokenManager, redisClient)

	s.server = &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	return s
}

// Start starts the HTTP server
func (s *Server) Start(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.server.Shutdown(shutdownCtx)
	}()

	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// setupRoutes configures all routes
func (s *Server) setupRoutes(tokenManager user.TokenManager, redisClient *redis.Client) {
	// Health check
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API v1
	v1 := s.router.Group("/api/v1")
	
	// Public routes
	auth := v1.Group("/auth")
	auth.Use(middleware.AuthRateLimit(redisClient))
	{
		auth.POST("/register", s.handlers.Auth.Register)
		auth.POST("/login", s.handlers.Auth.Login)
		auth.POST("/refresh", s.handlers.Auth.Refresh)
	}

	// Public auction routes
	auctions := v1.Group("/auctions")
	auctions.Use(middleware.GeneralRateLimit(redisClient))
	{
		auctions.GET("", s.handlers.Auction.ListLiveAuctions)
		auctions.GET("/live", s.handlers.Auction.ListLiveAuctions)
		auctions.GET("/:id", s.handlers.Auction.GetAuction)
	}

	// WebSocket endpoint for auctions (public, but auth recommended)
	s.router.GET("/ws/auctions/:id", func(c *gin.Context) {
		s.handlers.AuctionWS.HandleWebSocket(c)
	})

	// Upload endpoints (protected)
	uploads := v1.Group("/uploads")
	uploads.Use(middleware.AuthMiddleware(tokenManager))
	uploads.Use(middleware.GeneralRateLimit(redisClient))
	{
		uploads.POST("/product-image", s.handlers.Upload.UploadProductImage)
		uploads.POST("/avatar", s.handlers.Upload.UploadAvatar)
		uploads.POST("/stream-thumbnail", s.handlers.Upload.UploadStreamThumbnail)
		uploads.POST("/:folder", s.handlers.Upload.UploadGeneric)
		uploads.DELETE("", s.handlers.Upload.DeleteFile)
	}

	// Public product routes
	products := v1.Group("/products")
	products.Use(middleware.GeneralRateLimit(redisClient))
	{
		products.GET("", s.handlers.Product.List)
		products.GET("/slug/:slug", s.handlers.Product.GetBySlug)
		products.GET("/:id", s.handlers.Product.Get)
	}

	// Public category routes
	categories := v1.Group("/categories")
	categories.Use(middleware.GeneralRateLimit(redisClient))
	{
		categories.GET("", s.handlers.Category.List)
		categories.GET("/tree", s.handlers.Category.GetTree)
		categories.GET("/slug/:slug", s.handlers.Category.GetBySlug)
		categories.GET("/:id", s.handlers.Category.Get)
	}

	// Protected routes
	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware(tokenManager))
	protected.Use(middleware.GeneralRateLimit(redisClient))
	{
		// Auth
		protected.GET("/auth/profile", s.handlers.Auth.GetProfile)
		protected.POST("/auth/change-password", s.handlers.Auth.ChangePassword)
		protected.POST("/auth/logout", s.handlers.Auth.Logout)

		// Auctions (protected)
		protected.POST("/auctions", s.handlers.Auction.CreateAuction)
		protected.POST("/auctions/:id/bid", middleware.AuctionBidRateLimit(redisClient), s.handlers.Auction.PlaceBid)
		protected.POST("/auctions/:id/start", s.handlers.Auction.StartAuction)
		protected.POST("/auctions/:id/end", s.handlers.Auction.EndAuction)

		// Products (protected - seller only)
		protected.POST("/products", middleware.RequireRole(userDomain.RoleSeller), s.handlers.Product.Create)
		protected.PUT("/products/:id", middleware.RequireRole(userDomain.RoleSeller), s.handlers.Product.Update)
		protected.DELETE("/products/:id", middleware.RequireRole(userDomain.RoleSeller), s.handlers.Product.Delete)
		protected.POST("/products/:id/publish", middleware.RequireRole(userDomain.RoleSeller), s.handlers.Product.Publish)
		protected.POST("/products/:id/archive", middleware.RequireRole(userDomain.RoleSeller), s.handlers.Product.Archive)
		protected.GET("/my-products", middleware.RequireRole(userDomain.RoleSeller), s.handlers.Product.GetMyProducts)
	}

	// Admin routes
	admin := v1.Group("/admin")
	admin.Use(middleware.AuthMiddleware(tokenManager))
	admin.Use(middleware.RequireRole(userDomain.RoleAdmin))
	{
		// Categories (admin only)
		admin.POST("/categories", s.handlers.Category.Create)
		admin.PUT("/categories/:id", s.handlers.Category.Update)
		admin.DELETE("/categories/:id", s.handlers.Category.Delete)
	}

	// Admin routes
	admin := v1.Group("/admin")
	admin.Use(middleware.AuthMiddleware(tokenManager))
	admin.Use(middleware.RequireRole(userDomain.RoleAdmin))
	{
		// Admin endpoints here
	}
}

// corsMiddleware handles CORS
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	}
}

// loggerMiddleware logs requests
func loggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}