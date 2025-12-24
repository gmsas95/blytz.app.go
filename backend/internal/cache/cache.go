package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/blytz.live.remake/backend/pkg/logging"
	"github.com/redis/go-redis/v9"
)

// Cache provides Redis caching functionality
type Cache struct {
	client *redis.Client
	logger *logging.Logger
}

// NewCache creates a new Redis cache client
func NewCache(redisURL string) (*Cache, error) {
	logger := logging.NewLogger()
	
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opt)
	
	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	
	logger.Info("Redis cache connected successfully")
	
	return &Cache{
		client: client,
		logger: logger,
	}, nil
}

// Set stores a value in cache with expiration
func (c *Cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	// Serialize value to JSON
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}
	
	// Store in Redis
	err = c.client.Set(ctx, key, jsonValue, expiration).Err()
	if err != nil {
		c.logger.Error("Failed to set cache value", map[string]interface{}{
			"key":   key,
			"error": err.Error(),
		})
		return fmt.Errorf("failed to set cache value: %w", err)
	}
	
	c.logger.Debug("Cache value set", map[string]interface{}{
		"key":        key,
		"expiration": expiration,
	})
	
	return nil
}

// Get retrieves a value from cache
func (c *Cache) Get(ctx context.Context, key string, dest interface{}) error {
	// Get from Redis
	result, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			// Key not found is not an error
			return fmt.Errorf("cache miss")
		}
		c.logger.Error("Failed to get cache value", map[string]interface{}{
			"key":   key,
			"error": err.Error(),
		})
		return fmt.Errorf("failed to get cache value: %w", err)
	}
	
	// Deserialize JSON value
	err = json.Unmarshal([]byte(result), dest)
	if err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}
	
	c.logger.Debug("Cache value retrieved", map[string]interface{}{
		"key": key,
	})
	
	return nil
}

// Delete removes a value from cache
func (c *Cache) Delete(ctx context.Context, key string) error {
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		c.logger.Error("Failed to delete cache value", map[string]interface{}{
			"key":   key,
			"error": err.Error(),
		})
		return fmt.Errorf("failed to delete cache value: %w", err)
	}
	
	c.logger.Debug("Cache value deleted", map[string]interface{}{
		"key": key,
	})
	
	return nil
}

// SetMultiple stores multiple key-value pairs
func (c *Cache) SetMultiple(ctx context.Context, items map[string]interface{}, expiration time.Duration) error {
	pipe := c.client.Pipeline()
	
	for key, value := range items {
		jsonValue, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value for key %s: %w", key, err)
		}
		pipe.Set(ctx, key, jsonValue, expiration)
	}
	
	_, err := pipe.Exec(ctx)
	if err != nil {
		c.logger.Error("Failed to set multiple cache values", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("failed to set multiple cache values: %w", err)
	}
	
	c.logger.Debug("Multiple cache values set", map[string]interface{}{
		"count": len(items),
	})
	
	return nil
}

// GetMultiple retrieves multiple values from cache
func (c *Cache) GetMultiple(ctx context.Context, keys []string) (map[string]interface{}, error) {
	pipe := c.client.Pipeline()
	
	for _, key := range keys {
		pipe.Get(ctx, key)
	}
	
	results, err := pipe.Exec(ctx)
	if err != nil {
		c.logger.Error("Failed to get multiple cache values", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to get multiple cache values: %w", err)
	}
	
	values := make(map[string]interface{})
	
	for i, key := range keys {
		cmd := results[i]
		result, err := cmd.(*redis.StringCmd).Result()
		
		if err == redis.Nil {
			// Key not found, skip
			continue
		} else if err != nil {
			c.logger.Error("Failed to get cache result", map[string]interface{}{
				"key":   key,
				"error": err.Error(),
			})
			continue
		}
		
		var value interface{}
		err = json.Unmarshal([]byte(result), &value)
		if err != nil {
			c.logger.Error("Failed to unmarshal cache result", map[string]interface{}{
				"key":   key,
				"error": err.Error(),
			})
			continue
		}
		
		values[key] = value
	}
	
	c.logger.Debug("Multiple cache values retrieved", map[string]interface{}{
		"requested": len(keys),
		"found":     len(values),
	})
	
	return values, nil
}

// Increment increments a numeric value
func (c *Cache) Increment(ctx context.Context, key string) (int64, error) {
	result, err := c.client.Incr(ctx, key).Result()
	if err != nil {
		c.logger.Error("Failed to increment cache value", map[string]interface{}{
			"key":   key,
			"error": err.Error(),
		})
		return 0, fmt.Errorf("failed to increment cache value: %w", err)
	}
	
	c.logger.Debug("Cache value incremented", map[string]interface{}{
		"key":     key,
		"new_val": result,
	})
	
	return result, nil
}

// IncrementBy increments a numeric value by a specific amount
func (c *Cache) IncrementBy(ctx context.Context, key string, value int64) (int64, error) {
	result, err := c.client.IncrBy(ctx, key, value).Result()
	if err != nil {
		c.logger.Error("Failed to increment cache value by amount", map[string]interface{}{
			"key":   key,
			"value": value,
			"error": err.Error(),
		})
		return 0, fmt.Errorf("failed to increment cache value by amount: %w", err)
	}
	
	c.logger.Debug("Cache value incremented by amount", map[string]interface{}{
		"key":     key,
		"amount":  value,
		"new_val": result,
	})
	
	return result, nil
}

// Exists checks if a key exists in cache
func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	result, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		c.logger.Error("Failed to check cache existence", map[string]interface{}{
			"key":   key,
			"error": err.Error(),
		})
		return false, fmt.Errorf("failed to check cache existence: %w", err)
	}
	
	return result > 0, nil
}

// Expire sets expiration for a key
func (c *Cache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	err := c.client.Expire(ctx, key, expiration).Err()
	if err != nil {
		c.logger.Error("Failed to set cache expiration", map[string]interface{}{
			"key":        key,
			"expiration": expiration,
			"error":      err.Error(),
		})
		return fmt.Errorf("failed to set cache expiration: %w", err)
	}
	
	return nil
}

// GetTTL returns the remaining time to live for a key
func (c *Cache) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	result, err := c.client.TTL(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get cache TTL: %w", err)
	}
	
	return result, nil
}

// Clear removes all keys matching a pattern
func (c *Cache) Clear(ctx context.Context, pattern string) error {
	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get keys for pattern: %w", err)
	}
	
	if len(keys) == 0 {
		return nil
	}
	
	err = c.client.Del(ctx, keys...).Err()
	if err != nil {
		c.logger.Error("Failed to clear cache keys", map[string]interface{}{
			"pattern": pattern,
			"count":   len(keys),
			"error":   err.Error(),
		})
		return fmt.Errorf("failed to clear cache keys: %w", err)
	}
	
	c.logger.Info("Cache keys cleared", map[string]interface{}{
		"pattern": pattern,
		"count":   len(keys),
	})
	
	return nil
}

// Close closes the Redis connection
func (c *Cache) Close() error {
	err := c.client.Close()
	if err != nil {
		c.logger.Error("Failed to close Redis connection", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("failed to close Redis connection: %w", err)
	}
	
	c.logger.Info("Redis connection closed")
	return nil
}

// GetClient returns the underlying Redis client for advanced operations
func (c *Cache) GetClient() *redis.Client {
	return c.client
}