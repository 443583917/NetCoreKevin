package lock

import (
	"context"
	"time"
)

// Lock is the interface for distributed lock
type Lock interface {
	// Acquire acquires the lock with timeout
	Acquire(ctx context.Context, key string, ttl time.Duration) (bool, error)

	// Release releases the lock
	Release(ctx context.Context, key string) error

	// IsLocked checks if the key is locked
	IsLocked(ctx context.Context, key string) (bool, error)
}
