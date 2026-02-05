package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/blytz/live/backend/internal/infrastructure/cache/redis"
	"github.com/gin-gonic/gin"
)

// RateLimiter implements Redis-backed rate limiting
type RateLimiter struct {
	client *redis.Client
	window time.Duration
	max    int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(client *redis.Client, max int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		client: client,
		window: window,
		max:    max,
	}
}

// Allow checks if request should be allowed
func (rl *RateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	pipe := rl.client.GetClient().Pipeline()
	
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, rl.window)
	
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}
	
	count := incr.Val()
	return count <= int64(rl.max), nil
}

// Middleware creates Gin middleware
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := fmt.Sprintf("ratelimit:%s:%s", c.ClientIP(), c.FullPath())
		
		allowed, err := rl.Allow(c.Request.Context(), key)
		if err != nil {
			// Fail open on error
			c.Next()
			return
		}
		
		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error":   "RATE_LIMIT_EXCEEDED",
				"message": "too many requests",
				"retry_after": int(rl.window.Seconds()),
			})
			return
		}
		
		c.Next()
	}
}

// AuthRateLimit creates stricter rate limiting for auth endpoints
func AuthRateLimit(client *redis.Client) gin.HandlerFunc {
	// 5 requests per minute for auth
	rl := NewRateLimiter(client, 5, time.Minute)
	return rl.Middleware()
}

// GeneralRateLimit creates general rate limiting
func GeneralRateLimit(client *redis.Client) gin.HandlerFunc {
	// 100 requests per minute general
	rl := NewRateLimiter(client, 100, time.Minute)
	return rl.Middleware()
}

// AuctionBidRateLimit creates rate limiting for bid placement
func AuctionBidRateLimit(client *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		auctionID := c.Param("id")
		
		key := fmt.Sprintf("ratelimit:bid:%s:%s", userID.(string), auctionID)
		
		// 10 bids per minute per user per auction
		rl := NewRateLimiter(client, 10, time.Minute)
		
		allowed, err := rl.Allow(c.Request.Context(), key)
		if err != nil {
			c.Next()
			return
		}
		
		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error":   "RATE_LIMIT_EXCEEDED",
				"message": "too many bids, please slow down",
			})
			return
		}
		
		c.Next()
	}
}