package lock

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrLockNotAcquired = errors.New("lock not acquired")
	ErrLockNotHeld     = errors.New("lock not held")
)

// RedisLock implements distributed lock using Redis
type RedisLock struct {
	client *redis.Client
	prefix string
}

// NewRedisLock creates a new Redis lock
func NewRedisLock(client *redis.Client) *RedisLock {
	return &RedisLock{
		client: client,
		prefix: "lock:",
	}
}

// Acquire acquires the lock
func (l *RedisLock) Acquire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	fullKey := l.prefix + key
	result, err := l.client.SetNX(ctx, fullKey, "1", ttl).Result()
	if err != nil {
		return false, err
	}
	return result, nil
}

// Release releases the lock
func (l *RedisLock) Release(ctx context.Context, key string) error {
	fullKey := l.prefix + key
	deleted, err := l.client.Del(ctx, fullKey).Result()
	if err != nil {
		return err
	}
	if deleted == 0 {
		return ErrLockNotHeld
	}
	return nil
}

// IsLocked checks if the key is locked
func (l *RedisLock) IsLocked(ctx context.Context, key string) (bool, error) {
	fullKey := l.prefix + key
	result, err := l.client.Exists(ctx, fullKey).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}
