package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/blytz/live/backend/internal/domain/auction"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// AuctionCache implements auction.Cache interface
type AuctionCache struct {
	client *Client
	prefix string
}

// NewAuctionCache creates a new auction cache
func NewAuctionCache(client *Client) *AuctionCache {
	return &AuctionCache{
		client: client,
		prefix: "auction:",
	}
}

// GetAuctionState retrieves auction state from cache
func (c *AuctionCache) GetAuctionState(ctx context.Context, auctionID uuid.UUID) (*auction.AuctionState, error) {
	key := c.key(auctionID)
	var state auction.AuctionState
	if err := c.client.Get(ctx, key, &state); err != nil {
		return nil, err
	}
	return &state, nil
}

// SetAuctionState stores auction state in cache
func (c *AuctionCache) SetAuctionState(ctx context.Context, auctionID uuid.UUID, state *auction.AuctionState, ttl time.Duration) error {
	key := c.key(auctionID)
	return c.client.Set(ctx, key, state, ttl)
}

// DeleteAuctionState removes auction state from cache
func (c *AuctionCache) DeleteAuctionState(ctx context.Context, auctionID uuid.UUID) error {
	key := c.key(auctionID)
	return c.client.Delete(ctx, key)
}

// IncrementViewerCount increments viewer count for an auction
func (c *AuctionCache) IncrementViewerCount(ctx context.Context, auctionID uuid.UUID) (int, error) {
	key := c.key(auctionID) + ":viewers"
	count, err := c.client.Increment(ctx, key)
	if err != nil {
		return 0, err
	}
	// Set expiration
	c.client.Expire(ctx, key, time.Hour)
	return int(count), nil
}

// DecrementViewerCount decrements viewer count for an auction
func (c *AuctionCache) DecrementViewerCount(ctx context.Context, auctionID uuid.UUID) (int, error) {
	key := c.key(auctionID) + ":viewers"
	count, err := c.client.Decrement(ctx, key)
	if err != nil {
		return 0, err
	}
	if count < 0 {
		count = 0
		c.client.Set(ctx, key, 0, time.Hour)
	}
	return int(count), nil
}

// GetViewerCount gets current viewer count
func (c *AuctionCache) GetViewerCount(ctx context.Context, auctionID uuid.UUID) (int, error) {
	key := c.key(auctionID) + ":viewers"
	
	data, err := c.client.GetClient().Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	
	var count int
	if err := json.Unmarshal([]byte(data), &count); err != nil {
		return 0, err
	}
	return count, nil
}

// key generates cache key
func (c *AuctionCache) key(auctionID uuid.UUID) string {
	return fmt.Sprintf("%s%s:state", c.prefix, auctionID.String())
}