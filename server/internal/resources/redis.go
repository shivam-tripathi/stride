package resources

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"quizizz.com/internal/config"
	"quizizz.com/internal/logger"
)

// Redis implements the RedisResource interface using go-redis
type Redis struct {
	client *redis.Client
	config config.RedisConfig
	tracer trace.Tracer
}

// NewRedis creates a new Redis resource
func NewRedis(cfg *config.Config) RedisResource {
	return &Redis{
		config: cfg.Redis,
		tracer: otel.Tracer("redis"),
	}
}

// Connect establishes a connection to Redis
func (r *Redis) Connect(ctx context.Context) error {
	ctx, span := r.tracer.Start(ctx, "Redis.Connect",
		trace.WithAttributes(
			semconv.DBSystemRedis,
			attribute.String("redis.host", r.config.Host),
			attribute.String("redis.port", r.config.Port),
			attribute.Int("redis.db", r.config.DB),
		),
	)
	defer span.End()

	logger.InfoCtx(ctx, "Connecting to Redis",
		zap.String("host", r.config.Host),
		zap.String("port", r.config.Port),
		zap.Int("db", r.config.DB),
	)

	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", r.config.Host, r.config.Port),
		Password:     r.config.Password,
		DB:           r.config.DB,
		DialTimeout:  r.config.Timeout,
		ReadTimeout:  r.config.Timeout,
		WriteTimeout: r.config.Timeout,
	})

	r.client = client

	// Verify the connection
	if err := r.Ping(ctx); err != nil {
		span.RecordError(err)
		return err
	}

	logger.InfoCtx(ctx, "Successfully connected to Redis")
	return nil
}

// Close closes the Redis connection
func (r *Redis) Close(ctx context.Context) error {
	ctx, span := r.tracer.Start(ctx, "Redis.Close")
	defer span.End()

	if r.client != nil {
		logger.InfoCtx(ctx, "Closing Redis connection")
		return r.client.Close()
	}
	return nil
}

// Ping checks the Redis connection
func (r *Redis) Ping(ctx context.Context) error {
	ctx, span := r.tracer.Start(ctx, "Redis.Ping")
	defer span.End()

	if r.client == nil {
		err := fmt.Errorf("Redis connection not established")
		span.RecordError(err)
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, r.config.Timeout)
	defer cancel()

	_, err := r.client.Ping(ctx).Result()
	if err != nil {
		span.RecordError(err)
	}
	return err
}

// Name returns the name of the resource
func (r *Redis) Name() string {
	return "redis"
}

// Client returns the Redis client
func (r *Redis) Client() interface{} {
	return r.client
}

// GetClient returns the underlying redis.Client instance
func (r *Redis) GetClient() *redis.Client {
	return r.client
}

// WithContext creates a new traced context for Redis operations
func (r *Redis) WithContext(ctx context.Context, operation string) (context.Context, trace.Span) {
	return r.tracer.Start(ctx, operation,
		trace.WithAttributes(
			semconv.DBSystemRedis,
			attribute.String("redis.operation", operation),
			attribute.Int("redis.db", r.config.DB),
		),
	)
}
