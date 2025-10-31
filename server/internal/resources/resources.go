// Package resources provides interfaces and implementations for external resources
package resources

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"quizizz.com/internal/logger"
)

// Common errors
var (
	ErrResourceNotConnected = errors.New("resource not connected")
)

// Resources holds all the application resources
type Resources struct {
	DB    DBResource
	Redis RedisResource
}

// Resource defines the interface that all resources must implement
type Resource interface {
	// Connect establishes a connection to the resource
	Connect(ctx context.Context) error

	// Close closes the connection to the resource
	Close(ctx context.Context) error

	// Ping checks if the resource is available
	Ping(ctx context.Context) error

	// Name returns the name of the resource
	Name() string
}

// HealthCheck performs a health check on a resource
type HealthCheck struct {
	Name    string    `json:"name"`
	Status  string    `json:"status"`
	Message string    `json:"message,omitempty"`
	Time    time.Time `json:"time"`
}

// CheckHealth checks the health of a resource
func CheckHealth(ctx context.Context, res Resource) HealthCheck {
	start := time.Now()
	err := res.Ping(ctx)

	health := HealthCheck{
		Name: res.Name(),
		Time: time.Now(),
	}

	if err != nil {
		health.Status = "error"
		health.Message = err.Error()
		logger.Error("Resource health check failed",
			zap.String("resource", res.Name()),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
	} else {
		health.Status = "ok"
		logger.Debug("Resource health check passed",
			zap.String("resource", res.Name()),
			zap.Duration("duration", time.Since(start)),
		)
	}

	return health
}

// DBResource defines the interface for database resources
type DBResource interface {
	Resource

	// DB returns the database instance
	DB() interface{}
}

// RedisResource defines the interface for Redis resources
type RedisResource interface {
	Resource

	// Client returns the Redis client
	Client() interface{}
}

// InitResources initializes all resources
func InitResources(ctx context.Context, resources *Resources) error {
	logger.Info("Initializing resources")

	// Initialize DB
	if err := resources.DB.Connect(ctx); err != nil {
		return err
	}

	// Initialize Redis
	if err := resources.Redis.Connect(ctx); err != nil {
		return err
	}

	logger.Info("All resources initialized successfully")
	return nil
}

// CloseResources closes all resources
func CloseResources(ctx context.Context, resources *Resources) {
	logger.Info("Closing resources")

	// Close DB
	if err := resources.DB.Close(ctx); err != nil {
		logger.Error("Failed to close DB", zap.Error(err))
	}

	// Close Redis
	if err := resources.Redis.Close(ctx); err != nil {
		logger.Error("Failed to close Redis", zap.Error(err))
	}

	logger.Info("All resources closed")
}
