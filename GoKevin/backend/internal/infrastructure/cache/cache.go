package cache

import (
	"context"
	"time"
)

// Cache is the interface for cache operations
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, keys ...string) (bool, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	GetJSON(ctx context.Context, key string, dest interface{}) error
	SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error
}
