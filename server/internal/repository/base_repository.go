package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"quizizz.com/internal/logger"
)

// Common repository errors
var (
	ErrNotFound      = errors.New("document not found")
	ErrAlreadyExists = errors.New("document already exists")
	ErrInvalidID     = errors.New("invalid document ID")
	ErrInvalidInput  = errors.New("invalid input")
)

// BaseRepository provides common MongoDB operations using generics for type safety
// T is the document type (e.g., userDocument, productDocument)
type BaseRepository[T any] struct {
	collection *mongo.Collection
	tracer     trace.Tracer
	entityName string // For better error messages
}

// BaseRepositoryConfig configures a BaseRepository
type BaseRepositoryConfig struct {
	Collection *mongo.Collection
	EntityName string // e.g., "user", "product" - used in error messages
}

// NewBaseRepository creates a new BaseRepository with generic type
func NewBaseRepository[T any](collection *mongo.Collection) *BaseRepository[T] {
	return &BaseRepository[T]{
		collection: collection,
		tracer:     otel.Tracer("repository"),
		entityName: collection.Name(),
	}
}

// NewBaseRepositoryWithConfig creates a new BaseRepository with configuration
func NewBaseRepositoryWithConfig[T any](cfg BaseRepositoryConfig) *BaseRepository[T] {
	entityName := cfg.EntityName
	if entityName == "" {
		entityName = cfg.Collection.Name()
	}

	return &BaseRepository[T]{
		collection: cfg.Collection,
		tracer:     otel.Tracer("repository"),
		entityName: entityName,
	}
}

// EntityName returns the entity name for this repository
func (r *BaseRepository[T]) EntityName() string {
	return r.entityName
}

// FindByID finds a document by its ID and returns it
func (r *BaseRepository[T]) FindByID(ctx context.Context, id string) (*T, error) {
	ctx, span := r.tracer.Start(ctx, "BaseRepository.FindByID",
		trace.WithAttributes(
			attribute.String("collection", r.collection.Name()),
			attribute.String("id", id),
		),
	)
	defer span.End()

	// Convert string ID to ObjectID if needed
	var filter bson.M
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		// If not a valid ObjectID, search by string ID
		filter = bson.M{"_id": id}
	} else {
		filter = bson.M{"_id": objectID}
	}

	var result T
	err = r.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			span.RecordError(ErrNotFound)
			return nil, ErrNotFound
		}
		span.RecordError(err)
		logger.ErrorCtx(ctx, fmt.Sprintf("Failed to find %s by ID", r.entityName),
			zap.String("entity", r.entityName),
			zap.String("id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to find %s: %w", r.entityName, err)
	}

	return &result, nil
}

// FindOne finds a single document matching the filter
func (r *BaseRepository[T]) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (*T, error) {
	ctx, span := r.tracer.Start(ctx, "BaseRepository.FindOne",
		trace.WithAttributes(
			attribute.String("collection", r.collection.Name()),
		),
	)
	defer span.End()

	var result T
	err := r.collection.FindOne(ctx, filter, opts...).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		span.RecordError(err)
		logger.ErrorCtx(ctx, "Failed to find document",
			zap.String("collection", r.collection.Name()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to find document: %w", err)
	}

	return &result, nil
}

// Find finds multiple documents matching the filter
func (r *BaseRepository[T]) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) ([]T, error) {
	ctx, span := r.tracer.Start(ctx, "BaseRepository.Find",
		trace.WithAttributes(
			attribute.String("collection", r.collection.Name()),
		),
	)
	defer span.End()

	cursor, err := r.collection.Find(ctx, filter, opts...)
	if err != nil {
		span.RecordError(err)
		logger.ErrorCtx(ctx, "Failed to find documents",
			zap.String("collection", r.collection.Name()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(ctx)

	var results []T
	err = cursor.All(ctx, &results)
	if err != nil {
		span.RecordError(err)
		logger.ErrorCtx(ctx, "Failed to decode documents",
			zap.String("collection", r.collection.Name()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to decode documents: %w", err)
	}

	return results, nil
}

// FindAll finds all documents in the collection
func (r *BaseRepository[T]) FindAll(ctx context.Context, opts ...*options.FindOptions) ([]T, error) {
	return r.Find(ctx, bson.M{}, opts...)
}

// InsertOne inserts a single document
func (r *BaseRepository[T]) InsertOne(ctx context.Context, document *T) (string, error) {
	ctx, span := r.tracer.Start(ctx, "BaseRepository.InsertOne",
		trace.WithAttributes(
			attribute.String("collection", r.collection.Name()),
		),
	)
	defer span.End()

	result, err := r.collection.InsertOne(ctx, document)
	if err != nil {
		span.RecordError(err)
		logger.ErrorCtx(ctx, "Failed to insert document",
			zap.String("collection", r.collection.Name()),
			zap.Error(err),
		)
		// Check if it's a duplicate key error
		if mongo.IsDuplicateKeyError(err) {
			return "", ErrAlreadyExists
		}
		return "", fmt.Errorf("failed to insert document: %w", err)
	}

	// Extract the inserted ID
	var id string
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		id = oid.Hex()
	} else {
		id = fmt.Sprintf("%v", result.InsertedID)
	}

	return id, nil
}

// InsertMany inserts multiple documents
func (r *BaseRepository[T]) InsertMany(ctx context.Context, documents []*T) ([]string, error) {
	ctx, span := r.tracer.Start(ctx, "BaseRepository.InsertMany",
		trace.WithAttributes(
			attribute.String("collection", r.collection.Name()),
			attribute.Int("count", len(documents)),
		),
	)
	defer span.End()

	// Convert []*T to []interface{} for MongoDB driver
	docs := make([]interface{}, len(documents))
	for i, doc := range documents {
		docs[i] = doc
	}

	result, err := r.collection.InsertMany(ctx, docs)
	if err != nil {
		span.RecordError(err)
		logger.ErrorCtx(ctx, "Failed to insert documents",
			zap.String("collection", r.collection.Name()),
			zap.Error(err),
		)
		if mongo.IsDuplicateKeyError(err) {
			return nil, ErrAlreadyExists
		}
		return nil, fmt.Errorf("failed to insert documents: %w", err)
	}

	// Extract the inserted IDs
	ids := make([]string, len(result.InsertedIDs))
	for i, insertedID := range result.InsertedIDs {
		if oid, ok := insertedID.(primitive.ObjectID); ok {
			ids[i] = oid.Hex()
		} else {
			ids[i] = fmt.Sprintf("%v", insertedID)
		}
	}

	return ids, nil
}

// UpdateByID updates a document by its ID
func (r *BaseRepository[T]) UpdateByID(ctx context.Context, id string, update interface{}) error {
	ctx, span := r.tracer.Start(ctx, "BaseRepository.UpdateByID",
		trace.WithAttributes(
			attribute.String("collection", r.collection.Name()),
			attribute.String("id", id),
		),
	)
	defer span.End()

	// Convert string ID to ObjectID if needed
	var filter bson.M
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		filter = bson.M{"_id": id}
	} else {
		filter = bson.M{"_id": objectID}
	}

	// Ensure update has the correct format
	var updateDoc bson.M
	switch v := update.(type) {
	case bson.M:
		// If the update already has operators like $set, use it as is
		if hasOperators(v) {
			updateDoc = v
		} else {
			// Wrap it in $set
			updateDoc = bson.M{"$set": v}
		}
	default:
		updateDoc = bson.M{"$set": update}
	}

	// Always update the updatedAt field
	if setDoc, ok := updateDoc["$set"].(bson.M); ok {
		setDoc["updatedAt"] = time.Now()
	}

	result, err := r.collection.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		span.RecordError(err)
		logger.ErrorCtx(ctx, "Failed to update document",
			zap.String("collection", r.collection.Name()),
			zap.String("id", id),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update document: %w", err)
	}

	if result.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}

// UpdateOne updates a single document matching the filter
func (r *BaseRepository[T]) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) error {
	ctx, span := r.tracer.Start(ctx, "BaseRepository.UpdateOne",
		trace.WithAttributes(
			attribute.String("collection", r.collection.Name()),
		),
	)
	defer span.End()

	result, err := r.collection.UpdateOne(ctx, filter, update, opts...)
	if err != nil {
		span.RecordError(err)
		logger.ErrorCtx(ctx, "Failed to update document",
			zap.String("collection", r.collection.Name()),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update document: %w", err)
	}

	if result.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}

// UpdateMany updates multiple documents matching the filter
func (r *BaseRepository[T]) UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (int64, error) {
	ctx, span := r.tracer.Start(ctx, "BaseRepository.UpdateMany",
		trace.WithAttributes(
			attribute.String("collection", r.collection.Name()),
		),
	)
	defer span.End()

	result, err := r.collection.UpdateMany(ctx, filter, update, opts...)
	if err != nil {
		span.RecordError(err)
		logger.ErrorCtx(ctx, "Failed to update documents",
			zap.String("collection", r.collection.Name()),
			zap.Error(err),
		)
		return 0, fmt.Errorf("failed to update documents: %w", err)
	}

	return result.ModifiedCount, nil
}

// DeleteByID deletes a document by its ID
func (r *BaseRepository[T]) DeleteByID(ctx context.Context, id string) error {
	ctx, span := r.tracer.Start(ctx, "BaseRepository.DeleteByID",
		trace.WithAttributes(
			attribute.String("collection", r.collection.Name()),
			attribute.String("id", id),
		),
	)
	defer span.End()

	// Convert string ID to ObjectID if needed
	var filter bson.M
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		filter = bson.M{"_id": id}
	} else {
		filter = bson.M{"_id": objectID}
	}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		span.RecordError(err)
		logger.ErrorCtx(ctx, "Failed to delete document",
			zap.String("collection", r.collection.Name()),
			zap.String("id", id),
			zap.Error(err),
		)
		return fmt.Errorf("failed to delete document: %w", err)
	}

	if result.DeletedCount == 0 {
		return ErrNotFound
	}

	return nil
}

// DeleteOne deletes a single document matching the filter
func (r *BaseRepository[T]) DeleteOne(ctx context.Context, filter interface{}) error {
	ctx, span := r.tracer.Start(ctx, "BaseRepository.DeleteOne",
		trace.WithAttributes(
			attribute.String("collection", r.collection.Name()),
		),
	)
	defer span.End()

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		span.RecordError(err)
		logger.ErrorCtx(ctx, "Failed to delete document",
			zap.String("collection", r.collection.Name()),
			zap.Error(err),
		)
		return fmt.Errorf("failed to delete document: %w", err)
	}

	if result.DeletedCount == 0 {
		return ErrNotFound
	}

	return nil
}

// DeleteMany deletes multiple documents matching the filter
func (r *BaseRepository[T]) DeleteMany(ctx context.Context, filter interface{}) (int64, error) {
	ctx, span := r.tracer.Start(ctx, "BaseRepository.DeleteMany",
		trace.WithAttributes(
			attribute.String("collection", r.collection.Name()),
		),
	)
	defer span.End()

	result, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		span.RecordError(err)
		logger.ErrorCtx(ctx, "Failed to delete documents",
			zap.String("collection", r.collection.Name()),
			zap.Error(err),
		)
		return 0, fmt.Errorf("failed to delete documents: %w", err)
	}

	return result.DeletedCount, nil
}

// Count counts documents matching the filter
func (r *BaseRepository[T]) Count(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	ctx, span := r.tracer.Start(ctx, "BaseRepository.Count",
		trace.WithAttributes(
			attribute.String("collection", r.collection.Name()),
		),
	)
	defer span.End()

	count, err := r.collection.CountDocuments(ctx, filter, opts...)
	if err != nil {
		span.RecordError(err)
		logger.ErrorCtx(ctx, "Failed to count documents",
			zap.String("collection", r.collection.Name()),
			zap.Error(err),
		)
		return 0, fmt.Errorf("failed to count documents: %w", err)
	}

	return count, nil
}

// Exists checks if a document matching the filter exists
func (r *BaseRepository[T]) Exists(ctx context.Context, filter interface{}) (bool, error) {
	count, err := r.Count(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Aggregate performs an aggregation pipeline
func (r *BaseRepository[T]) Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) ([]T, error) {
	ctx, span := r.tracer.Start(ctx, "BaseRepository.Aggregate",
		trace.WithAttributes(
			attribute.String("collection", r.collection.Name()),
		),
	)
	defer span.End()

	cursor, err := r.collection.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		span.RecordError(err)
		logger.ErrorCtx(ctx, "Failed to aggregate documents",
			zap.String("collection", r.collection.Name()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to aggregate documents: %w", err)
	}
	defer cursor.Close(ctx)

	var results []T
	err = cursor.All(ctx, &results)
	if err != nil {
		span.RecordError(err)
		logger.ErrorCtx(ctx, "Failed to decode aggregation results",
			zap.String("collection", r.collection.Name()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to decode aggregation results: %w", err)
	}

	return results, nil
}

// Collection returns the underlying MongoDB collection
func (r *BaseRepository[T]) Collection() *mongo.Collection {
	return r.collection
}

// hasOperators checks if the update document has MongoDB update operators
func hasOperators(update bson.M) bool {
	for key := range update {
		if len(key) > 0 && key[0] == '$' {
			return true
		}
	}
	return false
}
