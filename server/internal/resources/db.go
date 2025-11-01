package resources

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"quizizz.com/internal/config"
	"quizizz.com/internal/logger"
)

// DB implements the DBResource interface using MongoDB
type DB struct {
	client   *mongo.Client
	database *mongo.Database
	config   config.MongoDBConfig
	tracer   trace.Tracer
}

// NewDB creates a new DB resource
func NewDB(cfg *config.Config) DBResource {
	return &DB{
		config: cfg.MongoDB,
		tracer: otel.Tracer("mongodb"),
	}
}

// Connect establishes a connection to the database
func (d *DB) Connect(ctx context.Context) error {
	ctx, span := d.tracer.Start(ctx, "MongoDB.Connect",
		trace.WithAttributes(
			attribute.String("db.system", "mongodb"),
			attribute.String("db.name", d.config.Database),
			attribute.String("db.connection_string", d.config.URI),
		),
	)
	defer span.End()

	logger.InfoCtx(ctx, "Connecting to MongoDB",
		zap.String("uri", d.config.URI),
		zap.String("database", d.config.Database),
	)

	// Create a connection timeout context
	connectCtx, cancel := context.WithTimeout(ctx, d.config.ConnectTimeout)
	defer cancel()

	// Configure MongoDB client options
	clientOptions := options.Client().
		ApplyURI(d.config.URI).
		SetMaxPoolSize(d.config.MaxPoolSize).
		SetMinPoolSize(d.config.MinPoolSize).
		SetServerSelectionTimeout(d.config.ConnectTimeout).
		SetMonitor(otelmongo.NewMonitor())

	// Connect to MongoDB
	client, err := mongo.Connect(connectCtx, clientOptions)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to connect to MongoDB", zap.Error(err))
		span.RecordError(err)
		return fmt.Errorf("failed to connect to mongodb: %w", err)
	}

	d.client = client
	d.database = client.Database(d.config.Database)

	// Verify the connection
	if err := d.Ping(ctx); err != nil {
		span.RecordError(err)
		return err
	}

	logger.InfoCtx(ctx, "Successfully connected to MongoDB")
	return nil
}

// Close closes the database connection
func (d *DB) Close(ctx context.Context) error {
	ctx, span := d.tracer.Start(ctx, "MongoDB.Close")
	defer span.End()

	if d.client != nil {
		logger.InfoCtx(ctx, "Closing MongoDB connection")
		return d.client.Disconnect(ctx)
	}
	return nil
}

// Ping checks the database connection
func (d *DB) Ping(ctx context.Context) error {
	ctx, span := d.tracer.Start(ctx, "MongoDB.Ping")
	defer span.End()

	if d.client == nil {
		err := fmt.Errorf("mongodb connection not established")
		span.RecordError(err)
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, d.config.Timeout)
	defer cancel()

	err := d.client.Ping(ctx, nil)
	if err != nil {
		span.RecordError(err)
	}
	return err
}

// Name returns the name of the resource
func (d *DB) Name() string {
	return "mongodb"
}

// DB returns the database instance
func (d *DB) DB() any {
	return d.database
}

// GetDatabase returns the underlying mongo.Database instance
func (d *DB) GetDatabase() *mongo.Database {
	return d.database
}

// GetClient returns the underlying mongo.Client instance
func (d *DB) GetClient() *mongo.Client {
	return d.client
}

// Collection returns a handle to a MongoDB collection
func (d *DB) Collection(name string) *mongo.Collection {
	return d.database.Collection(name)
}

// WithContext creates a new traced context for database operations
func (d *DB) WithContext(ctx context.Context, operation string) (context.Context, trace.Span) {
	return d.tracer.Start(ctx, operation,
		trace.WithAttributes(
			attribute.String("db.system", "mongodb"),
			attribute.String("db.name", d.config.Database),
			attribute.String("db.operation", operation),
		),
	)
}

// WithTransaction executes a function within a MongoDB transaction
func (d *DB) WithTransaction(ctx context.Context, fn func(sessCtx mongo.SessionContext) error) error {
	ctx, span := d.tracer.Start(ctx, "MongoDB.Transaction")
	defer span.End()

	session, err := d.client.StartSession()
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	// Execute the transaction
	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		return nil, fn(sessCtx)
	})

	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

// EnsureIndexes creates indexes for a collection
func (d *DB) EnsureIndexes(ctx context.Context, collectionName string, indexes []mongo.IndexModel) error {
	ctx, span := d.tracer.Start(ctx, "MongoDB.EnsureIndexes",
		trace.WithAttributes(
			attribute.String("collection", collectionName),
		),
	)
	defer span.End()

	collection := d.Collection(collectionName)
	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		span.RecordError(err)
		logger.ErrorCtx(ctx, "Failed to create indexes",
			zap.String("collection", collectionName),
			zap.Error(err),
		)
		return fmt.Errorf("failed to create indexes for %s: %w", collectionName, err)
	}

	logger.InfoCtx(ctx, "Successfully created indexes",
		zap.String("collection", collectionName),
		zap.Int("count", len(indexes)),
	)

	return nil
}

// HealthCheck performs a comprehensive health check
func (d *DB) HealthCheck(ctx context.Context) error {
	ctx, span := d.tracer.Start(ctx, "MongoDB.HealthCheck")
	defer span.End()

	// Check connection
	if err := d.Ping(ctx); err != nil {
		return err
	}

	// Run a simple command to ensure database is accessible
	result := d.database.RunCommand(ctx, bson.D{{Key: "ping", Value: 1}})
	if result.Err() != nil {
		span.RecordError(result.Err())
		return result.Err()
	}

	return nil
}
