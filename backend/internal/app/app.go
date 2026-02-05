package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/blytz/live/backend/internal/application/auction"
	"github.com/blytz/live/backend/internal/application/auth"
	"github.com/blytz/live/backend/internal/application/category"
	"github.com/blytz/live/backend/internal/application/product"
	userDomain "github.com/blytz/live/backend/internal/domain/user"
	"github.com/blytz/live/backend/internal/infrastructure/cache/redis"
	httpInfra "github.com/blytz/live/backend/internal/infrastructure/http"
	redisMessaging "github.com/blytz/live/backend/internal/infrastructure/messaging/redis"
	"github.com/blytz/live/backend/internal/infrastructure/persistence/postgres"
	"github.com/blytz/live/backend/internal/infrastructure/websocket"
	"github.com/blytz/live/backend/internal/interfaces/http/handlers"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

// Application is the main application container
type Application struct {
	config *Config
	db     *gorm.DB
	redis  *redis.Client
	
	// Services
	authService     *auth.Service
	auctionService  *auction.Service
	productService  *product.Service
	categoryService *category.Service
	
	// Infrastructure
	httpServer  *httpInfra.Server
	wsHub       *websocket.Hub
	eventBus    *redisMessaging.EventBus
	tokenManager userDomain.TokenManager
}

// Config holds application configuration
type Config struct {
	Environment string
	Port        string
	Database    postgres.Config
	Redis       redis.Config
	JWTSecret   string
}

// New creates a new Application instance
func New(cfg *Config) (*Application, error) {
	app := &Application{
		config: cfg,
	}

	// Initialize dependencies
	if err := app.initDatabase(); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	if err := app.initRedis(); err != nil {
		return nil, fmt.Errorf("failed to initialize Redis: %w", err)
	}

	if err := app.initTokenManager(); err != nil {
		return nil, fmt.Errorf("failed to initialize token manager: %w", err)
	}

	if err := app.initEventBus(); err != nil {
		return nil, fmt.Errorf("failed to initialize event bus: %w", err)
	}

	if err := app.initServices(); err != nil {
		return nil, fmt.Errorf("failed to initialize services: %w", err)
	}

	if err := app.initHTTPServer(); err != nil {
		return nil, fmt.Errorf("failed to initialize HTTP server: %w", err)
	}

	if err := app.initWebSocketHub(); err != nil {
		return nil, fmt.Errorf("failed to initialize WebSocket hub: %w", err)
	}

	return app, nil
}

// Run starts the application and blocks until shutdown
func (a *Application) Run(ctx context.Context) error {
	log.Println("Starting application...")

	g, ctx := errgroup.WithContext(ctx)

	// Start HTTP server
	g.Go(func() error {
		log.Printf("HTTP server starting on port %s", a.config.Port)
		return a.httpServer.Start(ctx)
	})

	// Start WebSocket hub
	g.Go(func() error {
		log.Println("WebSocket hub starting...")
		return a.wsHub.Start(ctx)
	})

	// Wait for shutdown signal
	<-ctx.Done()
	log.Println("Shutdown signal received, gracefully stopping...")

	// Give services time to finish
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}

	return g.Wait()
}

// Shutdown gracefully shuts down the application
func (a *Application) Shutdown(ctx context.Context) error {
	var errs []error

	if a.httpServer != nil {
		if err := a.httpServer.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("http server shutdown: %w", err))
		}
	}

	if a.wsHub != nil {
		if err := a.wsHub.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("websocket hub shutdown: %w", err))
		}
	}

	if a.redis != nil {
		if err := a.redis.Close(); err != nil {
			errs = append(errs, fmt.Errorf("redis close: %w", err))
		}
	}

	if a.db != nil {
		sqlDB, err := a.db.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				errs = append(errs, fmt.Errorf("database close: %w", err))
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("shutdown errors: %v", errs)
	}

	return nil
}

// initDatabase initializes the database connection
func (a *Application) initDatabase() error {
	db, err := postgres.NewConnection(a.config.Database)
	if err != nil {
		return err
	}

	// Run migrations
	if err := postgres.AutoMigrate(db); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	a.db = db
	log.Println("Database connected and migrated")
	return nil
}

// initRedis initializes Redis connection
func (a *Application) initRedis() error {
	client, err := redis.NewClient(a.config.Redis)
	if err != nil {
		return err
	}

	a.redis = client
	log.Println("Redis connected")
	return nil
}

// initTokenManager initializes JWT token manager
func (a *Application) initTokenManager() error {
	a.tokenManager = httpInfra.NewJWTTokenManager(
		a.config.JWTSecret,
		time.Hour,          // Access token expiry
		7*24*time.Hour,     // Refresh token expiry
	)
	return nil
}

// initEventBus initializes the event bus
func (a *Application) initEventBus() error {
	a.eventBus = redisMessaging.NewEventBus(a.redis.GetClient())
	log.Println("Event bus initialized")
	return nil
}

// initServices initializes domain services
func (a *Application) initServices() error {
	// Initialize repositories
	userRepo := postgres.NewUserRepository(a.db)
	auctionRepo := postgres.NewAuctionRepository(a.db)
	
	// Initialize repositories
	productRepo := postgres.NewProductRepository(a.db)
	categoryRepo := postgres.NewCategoryRepository(a.db)
	
	// Initialize auth service
	a.authService = auth.NewService(
		userRepo,
		nil, // Session repo - TODO: implement Redis session store
		a.tokenManager,
	)
	
	// Initialize auction service
	a.auctionService = auction.NewService(
		auctionRepo,
		nil, // Cache - TODO: implement Redis cache
		a.eventBus,
	)
	
	// Initialize product service
	a.productService = product.NewService(productRepo, categoryRepo)
	
	// Initialize category service
	a.categoryService = category.NewService(categoryRepo)

	log.Println("Services initialized")
	return nil
}

// initHTTPServer initializes the HTTP server
func (a *Application) initHTTPServer() error {
	handlers := &httpInfra.Handlers{
		Auth:     handlers.NewAuthHandler(a.authService),
		Auction:  handlers.NewAuctionHandler(a.auctionService),
		Product:  handlers.NewProductHandler(a.productService),
		Category: handlers.NewCategoryHandler(a.categoryService),
	}

	a.httpServer = httpInfra.NewServer(
		a.config.Port,
		handlers,
		a.tokenManager,
		a.redis,
	)
	return nil
}

// initWebSocketHub initializes the WebSocket hub
func (a *Application) initWebSocketHub() error {
	a.wsHub = websocket.NewHub(a.redis.GetClient(), a.eventBus)
	return nil
}