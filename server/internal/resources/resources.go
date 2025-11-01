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

// resourceInitResult holds the result of a resource initialization
type resourceInitResult struct {
	name     string
	resource Resource
	err      error
	duration time.Duration
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

// InitResources initializes all resources concurrently
func InitResources(ctx context.Context, resources *Resources) error {
	startTime := time.Now()
	logger.Info("Initializing resources concurrently")

	// Create a list of all resources to initialize
	resourcesList := []Resource{
		resources.DB,
		resources.Redis,
	}

	// Channel to collect initialization results
	resultsChan := make(chan resourceInitResult, len(resourcesList))

	// Connect each resource concurrently
	for _, res := range resourcesList {
		go func(resource Resource) {
			resStart := time.Now()
			name := resource.Name()

			logger.Info("Connecting to resource", zap.String("resource", name))

			err := resource.Connect(ctx)
			duration := time.Since(resStart)

			resultsChan <- resourceInitResult{
				name:     name,
				resource: resource,
				err:      err,
				duration: duration,
			}
		}(res)
	}

	// Collect all results
	var initErrors []error
	successCount := 0

	for i := 0; i < len(resourcesList); i++ {
		result := <-resultsChan

		if result.err != nil {
			logger.Error("Failed to connect to resource",
				zap.String("resource", result.name),
				zap.Error(result.err),
				zap.Duration("duration", result.duration),
			)
			initErrors = append(initErrors,
				errors.New(result.name+": "+result.err.Error()))
		} else {
			logger.Info("Successfully connected to resource",
				zap.String("resource", result.name),
				zap.Duration("duration", result.duration),
			)
			successCount++
		}
	}

	// If any initialization failed, return error with all failures
	if len(initErrors) > 0 {
		errorMsg := "failed to initialize resources: "
		for i, err := range initErrors {
			if i > 0 {
				errorMsg += "; "
			}
			errorMsg += err.Error()
		}
		return errors.New(errorMsg)
	}

	totalDuration := time.Since(startTime)
	logger.Info("All resources initialized successfully",
		zap.Int("count", successCount),
		zap.Duration("total_duration", totalDuration),
	)

	return nil
}

// CloseResources closes all resources concurrently
func CloseResources(ctx context.Context, resources *Resources) {
	startTime := time.Now()
	logger.Info("Closing resources")

	// Create a list of all resources to close
	resourcesList := []Resource{
		resources.DB,
		resources.Redis,
	}

	// Channel to collect close results
	resultsChan := make(chan resourceInitResult, len(resourcesList))

	// Close each resource concurrently
	for _, res := range resourcesList {
		go func(resource Resource) {
			resStart := time.Now()
			name := resource.Name()

			logger.Info("Closing resource", zap.String("resource", name))

			err := resource.Close(ctx)
			duration := time.Since(resStart)

			resultsChan <- resourceInitResult{
				name:     name,
				resource: resource,
				err:      err,
				duration: duration,
			}
		}(res)
	}

	// Collect all results
	successCount := 0
	failureCount := 0

	for range resourcesList {
		result := <-resultsChan

		if result.err != nil {
			logger.Error("Failed to close resource",
				zap.String("resource", result.name),
				zap.Error(result.err),
				zap.Duration("duration", result.duration),
			)
			failureCount++
		} else {
			logger.Info("Successfully closed resource",
				zap.String("resource", result.name),
				zap.Duration("duration", result.duration),
			)
			successCount++
		}
	}

	totalDuration := time.Since(startTime)
	logger.Info("Resource cleanup completed",
		zap.Int("success", successCount),
		zap.Int("failures", failureCount),
		zap.Duration("total_duration", totalDuration),
	)
}
