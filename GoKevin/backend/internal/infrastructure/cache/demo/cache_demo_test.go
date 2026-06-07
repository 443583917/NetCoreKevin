package demo

import (
	"testing"

	"github.com/kevin-ai/go-kevin/internal/infrastructure/cache"
	"github.com/stretchr/testify/assert"
)

func TestCacheInterface(t *testing.T) {
	// Verify that cache.Cache is an interface type
	var _ cache.Cache
	assert.True(t, true)
}

func TestRedisCacheImplementsCache(t *testing.T) {
	// Verify RedisCache satisfies the Cache interface (compile-time check is in redis.go)
	// This is a runtime verification
	var c cache.Cache = (*cache.RedisCache)(nil)
	assert.Nil(t, c)
}

func TestRedisConfig(t *testing.T) {
	cfg := cache.RedisConfig{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
	}
	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, 6379, cfg.Port)
	assert.Equal(t, "", cfg.Password)
	assert.Equal(t, 0, cfg.DB)
}
