# Repository Layer

This directory contains the repository layer implementations for data access in the application. The repository pattern provides an abstraction layer between the domain/business logic and the data source.

## Architecture

### Base Repository

The `BaseRepository` (`base_repository.go`) provides common MongoDB operations that can be embedded and reused across all domain-specific repositories. This promotes code reuse and consistency.

#### Features

- **CRUD Operations**: Create, Read, Update, Delete
- **Query Operations**: Find, FindOne, FindAll with flexible filtering
- **Aggregation**: Support for MongoDB aggregation pipelines
- **Counting**: Count documents matching criteria
- **Existence Checks**: Check if documents exist
- **OpenTelemetry Tracing**: Built-in distributed tracing for all operations
- **Error Handling**: Consistent error handling with custom error types

#### Common Methods

```go
// Finding documents
FindByID(ctx, id, result) error
FindOne(ctx, filter, result) error
Find(ctx, filter, results) error
FindAll(ctx, results) error

// Inserting documents
InsertOne(ctx, document) (string, error)
InsertMany(ctx, documents) ([]string, error)

// Updating documents
UpdateByID(ctx, id, update) error
UpdateOne(ctx, filter, update) error
UpdateMany(ctx, filter, update) (int64, error)

// Deleting documents
DeleteByID(ctx, id) error
DeleteOne(ctx, filter) error
DeleteMany(ctx, filter) (int64, error)

// Utilities
Count(ctx, filter) (int64, error)
Exists(ctx, filter) (bool, error)
Aggregate(ctx, pipeline, results) error
```

#### Error Types

- `ErrNotFound`: Document not found
- `ErrAlreadyExists`: Document with unique constraint already exists
- `ErrInvalidID`: Invalid document ID format
- `ErrInvalidInput`: Invalid input data

### Domain-Specific Repositories

Domain-specific repositories (like `MongoUserRepository`) embed the `BaseRepository` and add domain-specific logic:

```go
type MongoUserRepository struct {
    *BaseRepository
    db *resources.DB
}

func NewMongoUserRepository(db resources.DBResource) UserRepository {
    dbInstance := db.(*resources.DB)
    collection := dbInstance.Collection("users")

    return &MongoUserRepository{
        BaseRepository: NewBaseRepository(collection),
        db:             dbInstance,
    }
}
```

### Repository Interface

All domain repositories implement their respective interfaces (e.g., `UserRepository`):

```go
type UserRepository interface {
    GetByID(id string) (*domain.User, error)
    List() ([]*domain.User, error)
    Create(user *domain.User) error
    Update(user *domain.User) error
    Delete(id string) error
}
```

## MongoDB Implementation

### User Repository

The `MongoUserRepository` (`mongo_user_repository.go`) is a concrete implementation using MongoDB:

**Features:**
- Inherits all base CRUD operations
- Domain-specific methods (e.g., `GetByEmail`, `ExistsById`)
- Automatic index creation
- Document mapping between domain models and MongoDB documents

**Collection:** `users`

**Indexes:**
- Unique index on `email`
- Index on `createdAt` (descending for efficient listing)

**Example Usage:**

```go
// Create repository
userRepo := repository.NewMongoUserRepository(dbResource)

// Ensure indexes are created
if err := userRepo.EnsureIndexes(); err != nil {
    log.Fatal(err)
}

// Use repository
user := &domain.User{
    Name:  "John Doe",
    Email: "john@example.com",
}
err := userRepo.Create(user)
```

## Mock Repository

The `MockUserRepository` (`mock_user_repository.go`) provides an in-memory implementation for testing:

- Thread-safe using `sync.RWMutex`
- No external dependencies
- Fast and deterministic
- Useful for unit tests

## Creating a New Repository

To create a new repository for a domain entity:

1. **Define the repository interface:**

```go
type ProductRepository interface {
    GetByID(id string) (*domain.Product, error)
    List() ([]*domain.Product, error)
    Create(product *domain.Product) error
    Update(product *domain.Product) error
    Delete(id string) error
    // Add domain-specific methods
    GetBySKU(sku string) (*domain.Product, error)
}
```

2. **Create MongoDB implementation:**

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

3. **Implement interface methods:**

```go
func (r *MongoProductRepository) GetByID(id string) (*domain.Product, error) {
    ctx := context.Background()
    var doc ProductDocument

    err := r.FindByID(ctx, id, &doc)
    if err != nil {
        if err == ErrNotFound {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to get product: %w", err)
    }

    return r.documentToProduct(&doc), nil
}

// Implement other methods...
```

4. **Create indexes:**

```go
func (r *MongoProductRepository) EnsureIndexes() error {
    ctx := context.Background()

    indexes := []mongo.IndexModel{
        {
            Keys:    bson.D{{Key: "sku", Value: 1}},
            Options: options.Index().SetUnique(true),
        },
        {
            Keys: bson.D{{Key: "category", Value: 1}},
        },
    }

    return r.db.EnsureIndexes(ctx, "products", indexes)
}
```

5. **Create mock implementation:**

```go
type MockProductRepository struct {
    products map[string]*domain.Product
    mutex    sync.RWMutex
}

func NewMockProductRepository() ProductRepository {
    return &MockProductRepository{
        products: make(map[string]*domain.Product),
    }
}

// Implement interface methods...
```

6. **Wire it up in `wire/wire.go`:**

```go
func provideProductRepository(db resources.DBResource) repository.ProductRepository {
    return repository.NewMongoProductRepository(db)
}

var RepositorySet = wire.NewSet(
    provideUserRepository,
    provideProductRepository,
)
```

## Benefits of This Approach

1. **Code Reuse**: Common operations are implemented once in `BaseRepository`
2. **Consistency**: All repositories follow the same patterns
3. **Maintainability**: Changes to common logic only need to be made once
4. **Testability**: Easy to mock and test with `MockRepository` implementations
5. **Tracing**: Built-in OpenTelemetry tracing for observability
6. **Type Safety**: Compile-time type checking with Go's type system
7. **Flexibility**: Domain-specific methods can be added easily

## Best Practices

1. **Always use context**: Pass context through all repository methods for cancellation and tracing
2. **Handle errors consistently**: Use the predefined error types (`ErrNotFound`, `ErrAlreadyExists`, etc.)
3. **Create indexes**: Call `EnsureIndexes()` during application startup
4. **Document mapping**: Keep domain models and database documents separate with clear conversion functions
5. **Thread safety**: Make sure mock repositories are thread-safe for testing
6. **Transactions**: Use `db.WithTransaction()` for operations that need atomicity
7. **Logging**: Use the logger from the context for consistent logging

## Migration from PostgreSQL

This codebase was migrated from PostgreSQL to MongoDB. Key changes:

- **SQL → BSON**: Queries now use BSON documents instead of SQL
- **Tables → Collections**: PostgreSQL tables are now MongoDB collections
- **Joins → Embedded/Referenced**: Denormalization patterns instead of SQL joins
- **Transactions**: MongoDB supports multi-document transactions similar to PostgreSQL
- **Indexes**: Created using MongoDB's index API

## Environment Variables

Configure MongoDB connection via environment variables:

```bash
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=app
MONGODB_MAX_POOL_SIZE=100
MONGODB_MIN_POOL_SIZE=10
MONGODB_CONNECT_TIMEOUT=10s
MONGODB_TIMEOUT=5s
```

## Testing

### Unit Tests

Use `MockRepository` for unit tests:

```go
func TestUserService_Create(t *testing.T) {
    repo := repository.NewMockUserRepository()
    service := service.NewUserService(repo)

    user := &domain.User{
        Name:  "Test User",
        Email: "test@example.com",
    }

    err := service.Create(context.Background(), user)
    assert.NoError(t, err)
}
```

### Integration Tests

Use the actual MongoDB repository with a test database:

```go
func TestMongoUserRepository_Integration(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)

    repo := repository.NewMongoUserRepository(db)

    // Test repository operations
    user := &domain.User{
        Name:  "Test User",
        Email: "test@example.com",
    }

    err := repo.Create(user)
    assert.NoError(t, err)
    assert.NotEmpty(t, user.ID)
}
```

