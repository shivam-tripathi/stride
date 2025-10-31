package resources

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"quizizz.com/internal/config"
	"quizizz.com/internal/logger"
)

// DB implements the DBResource interface using sqlx
type DB struct {
	db     *sqlx.DB
	config config.DatabaseConfig
	tracer trace.Tracer
}

// NewDB creates a new DB resource
func NewDB(cfg *config.Config) DBResource {
	return &DB{
		config: cfg.Database,
		tracer: otel.Tracer("db"),
	}
}

// Connect establishes a connection to the database
func (d *DB) Connect(ctx context.Context) error {
	ctx, span := d.tracer.Start(ctx, "DB.Connect",
		trace.WithAttributes(
			semconv.DBSystemPostgreSQL,
			attribute.String("db.name", d.config.Name),
			attribute.String("db.user", d.config.User),
			attribute.String("db.host", d.config.Host),
			attribute.String("db.port", d.config.Port),
		),
	)
	defer span.End()

	logger.InfoCtx(ctx, "Connecting to database",
		zap.String("host", d.config.Host),
		zap.String("port", d.config.Port),
		zap.String("user", d.config.User),
		zap.String("database", d.config.Name),
		zap.String("sslmode", d.config.SSLMode),
	)

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.config.Host,
		d.config.Port,
		d.config.User,
		d.config.Password,
		d.config.Name,
		d.config.SSLMode,
	)

	db, err := sqlx.ConnectContext(ctx, "postgres", dsn)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to connect to database", zap.Error(err))
		span.RecordError(err)
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(d.config.MaxOpen)
	db.SetMaxIdleConns(d.config.MaxIdle)
	db.SetConnMaxLifetime(time.Minute)

	d.db = db

	// Verify the connection
	if err := d.Ping(ctx); err != nil {
		span.RecordError(err)
		return err
	}

	logger.InfoCtx(ctx, "Successfully connected to database")
	return nil
}

// Close closes the database connection
func (d *DB) Close(ctx context.Context) error {
	ctx, span := d.tracer.Start(ctx, "DB.Close")
	defer span.End()

	if d.db != nil {
		logger.InfoCtx(ctx, "Closing database connection")
		return d.db.Close()
	}
	return nil
}

// Ping checks the database connection
func (d *DB) Ping(ctx context.Context) error {
	ctx, span := d.tracer.Start(ctx, "DB.Ping")
	defer span.End()

	if d.db == nil {
		err := fmt.Errorf("database connection not established")
		span.RecordError(err)
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, d.config.Timeout)
	defer cancel()

	err := d.db.PingContext(ctx)
	if err != nil {
		span.RecordError(err)
	}
	return err
}

// Name returns the name of the resource
func (d *DB) Name() string {
	return "postgres"
}

// DB returns the database instance
func (d *DB) DB() any {
	return d.db
}

// GetDB returns the underlying sqlx.DB instance
func (d *DB) GetDB() *sqlx.DB {
	return d.db
}

// WithContext creates a new traced context for database operations
func (d *DB) WithContext(ctx context.Context, operation string) (context.Context, trace.Span) {
	return d.tracer.Start(ctx, operation,
		trace.WithAttributes(
			semconv.DBSystemPostgreSQL,
			attribute.String("db.name", d.config.Name),
			attribute.String("db.operation", operation),
		),
	)
}
