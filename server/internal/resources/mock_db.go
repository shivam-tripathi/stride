//go:build test || integration
// +build test integration

package resources

import (
	"context"

	"quizizz.com/internal/config"
)

// MockDB is a mock implementation of DBResource for testing
type MockDB struct {
	connected bool
	config    config.DatabaseConfig
}

// NewMockDB creates a new MockDB resource
func NewMockDB(cfg *config.Config) DBResource {
	return &MockDB{
		config: cfg.Database,
	}
}

// Connect simulates establishing a connection to the database
func (d *MockDB) Connect(ctx context.Context) error {
	d.connected = true
	return nil
}

// Close simulates closing the database connection
func (d *MockDB) Close(ctx context.Context) error {
	d.connected = false
	return nil
}

// Ping simulates checking the database connection
func (d *MockDB) Ping(ctx context.Context) error {
	if !d.connected {
		return ErrResourceNotConnected
	}
	return nil
}

// Name returns the name of the resource
func (d *MockDB) Name() string {
	return "mock-database"
}

// DB returns a mock database instance (nil for now since we're using mock repositories)
func (d *MockDB) DB() interface{} {
	return nil // Mock implementation doesn't provide actual DB instance
}
