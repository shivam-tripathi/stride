//go:build test || integration
// +build test integration

package resources

import (
	"context"

	"quizizz.com/internal/config"
)

// MockRedis is a mock implementation of RedisResource for testing
type MockRedis struct {
	connected bool
	config    config.RedisConfig
}

// NewMockRedis creates a new MockRedis resource
func NewMockRedis(cfg *config.Config) RedisResource {
	return &MockRedis{
		config: cfg.Redis,
	}
}

// Connect simulates establishing a connection to Redis
func (r *MockRedis) Connect(ctx context.Context) error {
	r.connected = true
	return nil
}

// Close simulates closing the Redis connection
func (r *MockRedis) Close(ctx context.Context) error {
	r.connected = false
	return nil
}

// Ping simulates checking the Redis connection
func (r *MockRedis) Ping(ctx context.Context) error {
	if !r.connected {
		return ErrResourceNotConnected
	}
	return nil
}

// Name returns the name of the resource
func (r *MockRedis) Name() string {
	return "mock-redis"
}

// Client returns a mock Redis client (nil for now since we're not using Redis in current tests)
func (r *MockRedis) Client() interface{} {
	return nil // Mock implementation doesn't provide actual Redis client
}
