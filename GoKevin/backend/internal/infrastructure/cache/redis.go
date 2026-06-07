package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache implements the Cache interface using Redis
type RedisCache struct {
	client *redis.Client
}

// RedisConfig represents Redis configuration
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// NewRedisCache creates a new Redis cache
func NewRedisCache(cfg RedisConfig) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return &RedisCache{client: client}
}

// NewRedisCacheWithClient creates a Redis cache with existing client
func NewRedisCacheWithClient(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

// Get gets a value from cache
func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

// Set sets a value in cache
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

// Delete deletes keys from cache
func (c *RedisCache) Delete(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

// Exists checks if keys exist
func (c *RedisCache) Exists(ctx context.Context, keys ...string) (bool, error) {
	result, err := c.client.Exists(ctx, keys...).Result()
	return result > 0, err
}

// Expire sets expiration for a key
func (c *RedisCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.client.Expire(ctx, key, expiration).Err()
}

// GetJSON gets a JSON value from cache
func (c *RedisCache) GetJSON(ctx context.Context, key string, dest interface{}) error {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

// SetJSON sets a JSON value in cache
func (c *RedisCache) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, expiration).Err()
}

// Ping checks Redis connection
func (c *RedisCache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

// Close closes the Redis connection
func (c *RedisCache) Close() error {
	return c.client.Close()
}

// Ensure RedisCache implements Cache at compile time
var _ Cache = (*RedisCache)(nil)
