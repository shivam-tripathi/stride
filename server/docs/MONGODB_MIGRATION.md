# MongoDB Migration Guide

This document outlines the migration from PostgreSQL to MongoDB and how to use the new database system.

## Overview

The application has been migrated from PostgreSQL to MongoDB to leverage:
- Flexible schema design
- Horizontal scalability
- Better performance for document-based data
- Native support for nested documents and arrays
- Simplified data modeling

## Changes Made

### 1. Configuration (`internal/config/config.go`)

**Before (PostgreSQL):**
```go
type DatabaseConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    Name     string
    SSLMode  string
    MaxOpen  int
    MaxIdle  int
    Timeout  time.Duration
}
```

**After (MongoDB):**
```go
type MongoDBConfig struct {
    URI            string
    Database       string
    MaxPoolSize    uint64
    MinPoolSize    uint64
    ConnectTimeout time.Duration
    Timeout        time.Duration
}
```

### 2. Database Resource (`internal/resources/db.go`)

**Before:**
- Used `sqlx` for PostgreSQL connections
- SQL-based operations
- Connection pooling with `SetMaxOpenConns` and `SetMaxIdleConns`

**After:**
- Uses MongoDB Go driver (`mongo-driver`)
- BSON-based operations
- Connection pooling with `SetMaxPoolSize` and `SetMinPoolSize`
- Integrated OpenTelemetry tracing with `otelmongo`
- Additional helper methods:
  - `Collection(name)`: Get a collection handle
  - `WithTransaction()`: Execute operations within a transaction
  - `EnsureIndexes()`: Create indexes for collections
  - `HealthCheck()`: Comprehensive health check

### 3. Repository Layer

**New Base Repository** (`internal/repository/base_repository.go`):
- Provides common CRUD operations
- Can be embedded in domain-specific repositories
- Includes tracing, logging, and error handling
- Methods:
  - Find operations: `FindByID`, `FindOne`, `Find`, `FindAll`
  - Insert operations: `InsertOne`, `InsertMany`
  - Update operations: `UpdateByID`, `UpdateOne`, `UpdateMany`
  - Delete operations: `DeleteByID`, `DeleteOne`, `DeleteMany`
  - Utility operations: `Count`, `Exists`, `Aggregate`

**User Repository** (`internal/repository/mongo_user_repository.go`):
- Implements `UserRepository` interface
- Embeds `BaseRepository` for common operations
- Adds domain-specific methods
- Manages MongoDB document structure
- Creates necessary indexes

### 4. Docker Configuration

**docker-compose.yml** now includes:
- MongoDB service (port 27017)
- Redis service (port 6379)
- API service with proper dependencies
- Volume mounts for data persistence
- Health checks for all services

## Environment Variables

### Old PostgreSQL Variables (Removed)
```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=app
DB_SSL_MODE=disable
DB_MAX_OPEN=10
DB_MAX_IDLE=5
DB_TIMEOUT=5s
```

### New MongoDB Variables
```bash
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=app
MONGODB_MAX_POOL_SIZE=100
MONGODB_MIN_POOL_SIZE=10
MONGODB_CONNECT_TIMEOUT=10s
MONGODB_TIMEOUT=5s
```

## Running the Application

### Local Development

1. **Start MongoDB and Redis:**
```bash
docker-compose up -d mongodb redis
```

2. **Set environment variables:**
```bash
export MONGODB_URI=mongodb://admin:password@localhost:27017
export MONGODB_DATABASE=app
```

3. **Run the application:**
```bash
make run
```

### Docker Compose

Start all services (MongoDB, Redis, API):
```bash
docker-compose up
```

### Connect to MongoDB

Using `mongosh`:
```bash
mongosh mongodb://admin:password@localhost:27017/app
```

Using MongoDB Compass:
```
mongodb://admin:password@localhost:27017/app
```

## Data Model Changes

### User Collection

**Document Structure:**
```javascript
{
  "_id": ObjectId("..."),
  "name": "John Doe",
  "email": "john@example.com",
  "createdAt": ISODate("2024-01-01T00:00:00Z"),
  "updatedAt": ISODate("2024-01-01T00:00:00Z")
}
```

**Indexes:**
- `email` (unique)
- `createdAt` (descending, for efficient listing)

## Query Examples

### Before (SQL/PostgreSQL)

```sql
-- Get user by ID
SELECT * FROM users WHERE id = $1;

-- List all users
SELECT * FROM users ORDER BY created_at DESC;

-- Create user
INSERT INTO users (name, email, created_at, updated_at)
VALUES ($1, $2, $3, $4) RETURNING id;

-- Update user
UPDATE users SET name = $1, email = $2, updated_at = $3
WHERE id = $4;

-- Delete user
DELETE FROM users WHERE id = $1;
```

### After (MongoDB/BSON)

```go
// Get user by ID
filter := bson.M{"_id": objectID}
collection.FindOne(ctx, filter).Decode(&user)

// List all users
opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
collection.Find(ctx, bson.M{}, opts)

// Create user
collection.InsertOne(ctx, userDoc)

// Update user
update := bson.M{"$set": bson.M{"name": name, "email": email}}
collection.UpdateOne(ctx, filter, update)

// Delete user
collection.DeleteOne(ctx, filter)
```

## Creating New Repositories

When creating a new repository for a domain entity:

1. **Define the repository interface** in `internal/repository/`:
```go
type ProductRepository interface {
    GetByID(id string) (*domain.Product, error)
    List() ([]*domain.Product, error)
    Create(product *domain.Product) error
    Update(product *domain.Product) error
    Delete(id string) error
}
```

2. **Create MongoDB implementation** that embeds `BaseRepository`:
```go
type MongoProductRepository struct {
    *BaseRepository
    db *resources.DB
}

func NewMongoProductRepository(db resources.DBResource) ProductRepository {
    dbInstance := db.(*resources.DB)
    collection := dbInstance.Collection("products")

    return &MongoProductRepository{
        BaseRepository: NewBaseRepository(collection),
        db:             dbInstance,
    }
}
```

3. **Implement interface methods** using base repository methods:
```go
func (r *MongoProductRepository) GetByID(id string) (*domain.Product, error) {
    ctx := context.Background()
    var doc ProductDocument

    if err := r.FindByID(ctx, id, &doc); err != nil {
        if err == ErrNotFound {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to get product: %w", err)
    }

    return r.documentToProduct(&doc), nil
}
```

4. **Create indexes** for your collection:
```go
func (r *MongoProductRepository) EnsureIndexes() error {
    ctx := context.Background()

    indexes := []mongo.IndexModel{
        {
            Keys:    bson.D{{Key: "sku", Value: 1}},
            Options: options.Index().SetUnique(true),
        },
    }

    return r.db.EnsureIndexes(ctx, "products", indexes)
}
```

5. **Wire it up** in `wire/wire.go`:
```go
func provideProductRepository(db resources.DBResource) repository.ProductRepository {
    return repository.NewMongoProductRepository(db)
}
```

## Testing

### Unit Tests

Use the mock repository for unit tests:
```go
func TestService(t *testing.T) {
    repo := repository.NewMockUserRepository()
    service := service.NewUserService(repo)

    // Test service logic
}
```

### Integration Tests

For integration tests with MongoDB:

1. Use a test database
2. Set up test data
3. Run tests
4. Clean up test data

```go
func TestMongoRepository(t *testing.T) {
    cfg := &config.Config{
        MongoDB: config.MongoDBConfig{
            URI:      "mongodb://localhost:27017",
            Database: "test_db",
        },
    }

    db := resources.NewDB(cfg)
    err := db.Connect(context.Background())
    require.NoError(t, err)
    defer db.Close(context.Background())

    repo := repository.NewMongoUserRepository(db)

    // Run tests...
}
```

## Common Patterns

### Querying

```go
// Simple query
filter := bson.M{"status": "active"}
repo.Find(ctx, filter, &results)

// Complex query
filter := bson.M{
    "age": bson.M{"$gte": 18},
    "status": "active",
}

// With options
opts := options.Find().
    SetSort(bson.D{{Key: "createdAt", Value: -1}}).
    SetLimit(10).
    SetSkip(0)
repo.Find(ctx, filter, &results, opts)
```

### Updating

```go
// Simple update
update := bson.M{"status": "inactive"}
repo.UpdateByID(ctx, id, update)

// With operators
update := bson.M{
    "$set": bson.M{"status": "inactive"},
    "$inc": bson.M{"loginCount": 1},
}
repo.UpdateOne(ctx, filter, update)
```

### Aggregation

```go
pipeline := []bson.M{
    {"$match": bson.M{"status": "active"}},
    {"$group": bson.M{
        "_id": "$category",
        "count": bson.M{"$sum": 1},
    }},
    {"$sort": bson.M{"count": -1}},
}

var results []CategoryCount
repo.Aggregate(ctx, pipeline, &results)
```

### Transactions

```go
err := db.WithTransaction(ctx, func(sessCtx mongo.SessionContext) error {
    // Perform multiple operations within transaction
    if err := repo1.Create(sessCtx, doc1); err != nil {
        return err
    }
    if err := repo2.Update(sessCtx, doc2); err != nil {
        return err
    }
    return nil
})
```

## Migration Checklist

- [x] Update configuration for MongoDB
- [x] Replace PostgreSQL driver with MongoDB driver
- [x] Create base repository with common operations
- [x] Implement MongoDB user repository
- [x] Update wire providers
- [x] Update Docker Compose configuration
- [x] Remove PostgreSQL dependencies from go.mod
- [x] Add MongoDB dependencies to go.mod
- [x] Update documentation

## Performance Considerations

1. **Indexes**: Always create indexes on frequently queried fields
2. **Connection Pooling**: Configure appropriate pool sizes
3. **Projections**: Use projections to fetch only required fields
4. **Batch Operations**: Use `InsertMany`, `UpdateMany` for bulk operations
5. **Aggregation**: Use aggregation pipelines for complex queries
6. **Sharding**: Consider sharding for horizontal scaling (future consideration)

## Monitoring

MongoDB metrics to monitor:
- Connection pool usage
- Query execution time
- Index hit ratio
- Document size
- Collection size
- Replication lag (if using replica sets)

Use the built-in OpenTelemetry tracing to monitor database operations.

## Troubleshooting

### Connection Issues

```go
// Check MongoDB connection
err := db.Ping(context.Background())
if err != nil {
    log.Fatal("MongoDB connection failed:", err)
}
```

### Index Creation

```go
// Ensure indexes are created on startup
if err := userRepo.EnsureIndexes(); err != nil {
    log.Fatal("Failed to create indexes:", err)
}
```

### Query Performance

Use MongoDB's explain plan:
```javascript
db.users.find({email: "test@example.com"}).explain("executionStats")
```

## Resources

- [MongoDB Go Driver Documentation](https://www.mongodb.com/docs/drivers/go/current/)
- [MongoDB Best Practices](https://www.mongodb.com/docs/manual/administration/production-notes/)
- [OpenTelemetry MongoDB Instrumentation](https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo)
- [Repository Pattern](https://martinfowler.com/eaaCatalog/repository.html)

## Support

For issues or questions:
1. Check the logs for error messages
2. Verify MongoDB connection and credentials
3. Check that indexes are created
4. Review the repository README for usage examples

