package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/blytz/live/backend/internal/app"
	"github.com/blytz/live/backend/internal/infrastructure/cache/redis"
	"github.com/blytz/live/backend/internal/infrastructure/persistence/postgres"
)

func main() {
	cfg := &app.Config{
		Environment: getEnv("ENV", "development"),
		Port:        getEnv("PORT", "8080"),
		JWTSecret:   getEnv("JWT_SECRET", "dev-secret-change-in-production"),
		Database: postgres.Config{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", ""),
			Database:        getEnv("DB_NAME", "blytz"),
			SSLMode:         getEnv("DB_SSL_MODE", "disable"),
			MaxOpenConns:    100,
			MaxIdleConns:    50,
		},
		Redis: redis.Config{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       0,
		},
	}

	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal")
		cancel()
	}()

	log.Printf("Starting server on port %s", cfg.Port)
	if err := application.Run(ctx); err != nil {
		log.Fatalf("Application error: %v", err)
	}

	log.Println("Server stopped gracefully")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
