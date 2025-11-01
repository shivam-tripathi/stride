# Generic BaseRepository Refactoring

## Summary

Refactored `BaseRepository` to use Go generics (Go 1.18+) for type-safe MongoDB operations, eliminating manual type casting and reducing boilerplate code.

## Changes

### Before (Interface-based)

```go
type BaseRepository struct {
    collection *mongo.Collection
}

// Methods required manual type casting
func (r *BaseRepository) FindByID(ctx context.Context, id string, result interface{}) error
func (r *BaseRepository) FindAll(ctx context.Context, results interface{}, opts ...*options.FindOptions) error
```

**Usage:**
```go
var doc UserDocument
err := r.FindByID(ctx, id, &doc)  // Manual type assertion required
```

### After (Generic-based)

```go
type BaseRepository[T any] struct {
    collection *mongo.Collection
}

// Methods return typed results
func (r *BaseRepository[T]) FindByID(ctx context.Context, id string) (*T, error)
func (r *BaseRepository[T]) FindAll(ctx context.Context, opts ...*options.FindOptions) ([]T, error)
```

**Usage:**
```go
doc, err := r.FindByID(ctx, id)  // Type-safe, no casting needed!
```

## Benefits

### 1. **Type Safety at Compile Time**
- ❌ Before: Runtime errors from incorrect type casting
- ✅ After: Compile-time type checking

### 2. **Cleaner Code**
```go
// Before (user_repository.go)
func (r *userRepositoryImpl) List(ctx context.Context) ([]*domain.User, error) {
    var docs []userDocument
    opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})

    if err := r.FindAll(ctx, &docs, opts); err != nil {  // Pass pointer
        return nil, err
    }
    return toUsers(docs), nil
}

// After (user_repository.go)
func (r *userRepositoryImpl) List(ctx context.Context) ([]*domain.User, error) {
    opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})

    docs, err := r.FindAll(ctx, opts)  // Direct return, no pointer gymnastics
    if err != nil {
        return nil, err
    }
    return toUsers(docs), nil
}
```

### 3. **Better IntelliSense/Autocomplete**
- IDE knows the exact return type
- Better code suggestions and error detection

### 4. **Reduced Boilerplate**
- No more `interface{}` parameters
- No manual type assertions
- Fewer lines of code overall

## Implementation Details

### Generic BaseRepository Declaration

```go
// Type parameter T represents the document type
type BaseRepository[T any] struct {
    collection *mongo.Collection
    tracer     trace.Tracer
    entityName string
}

// Constructor with generic type
func NewBaseRepositoryWithConfig[T any](cfg BaseRepositoryConfig) *BaseRepository[T] {
    return &BaseRepository[T]{
        collection: cfg.Collection,
        tracer:     otel.Tracer("repository"),
        entityName: cfg.EntityName,
    }
}
```

### Concrete Repository Implementation

```go
type userRepositoryImpl struct {
    *BaseRepository[userDocument]  // Specify concrete type
    db *resources.DB
}

func NewUserRepository(db resources.DBResource) UserRepository {
    return &userRepositoryImpl{
        BaseRepository: NewBaseRepositoryWithConfig[userDocument](  // Type parameter
            BaseRepositoryConfig{
                Collection: dbInstance.Collection("users"),
                EntityName: "user",
            }),
        db: dbInstance,
    }
}
```

## Method Signatures Comparison

### Read Operations

| Operation | Before | After |
|-----------|--------|-------|
| `FindByID` | `(ctx, id, result interface{}) error` | `(ctx, id) (*T, error)` |
| `FindOne` | `(ctx, filter, result interface{}, opts) error` | `(ctx, filter, opts) (*T, error)` |
| `Find` | `(ctx, filter, results interface{}, opts) error` | `(ctx, filter, opts) ([]T, error)` |
| `FindAll` | `(ctx, results interface{}, opts) error` | `(ctx, opts) ([]T, error)` |
| `Aggregate` | `(ctx, pipeline, results interface{}, opts) error` | `(ctx, pipeline, opts) ([]T, error)` |

### Write Operations

| Operation | Before | After |
|-----------|--------|-------|
| `InsertOne` | `(ctx, document interface{}) (string, error)` | `(ctx, document *T) (string, error)` |
| `InsertMany` | `(ctx, documents []interface{}) ([]string, error)` | `(ctx, documents []*T) ([]string, error)` |

## Real-World Usage Example

### Product Repository (hypothetical)

```go
type productDocument struct {
    ID    primitive.ObjectID `bson:"_id,omitempty"`
    Name  string            `bson:"name"`
    Price float64           `bson:"price"`
}

type productRepositoryImpl struct {
    *BaseRepository[productDocument]  // Reuse with type safety!
}

func (r *productRepositoryImpl) GetExpensiveProducts(ctx context.Context) ([]productDocument, error) {
    filter := bson.M{"price": bson.M{"$gt": 100}}
    return r.Find(ctx, filter)  // Type-safe, returns []productDocument
}
```

## Migration Notes

### Breaking Changes
- All repository implementations need to specify the document type parameter
- Method signatures changed from accepting `interface{}` to typed parameters/returns
- Tests need to be updated to work with new signatures

### Non-Breaking Changes
- Error handling logic remains the same
- MongoDB operations are unchanged
- Existing business logic is unaffected

## Performance

- **Zero runtime overhead**: Generics are compiled away at build time
- **Same performance**: Generated code is equivalent to hand-written type-specific code
- **Binary size**: Minimal increase (one instantiation per document type)

## Future Improvements

1. **Generic Query Builder**: Type-safe query construction
2. **Projection Support**: Type-safe field selection
3. **Nested Document Handling**: Better support for embedded documents
4. **Validation**: Generic validation interface

## Testing

All tests passing ✅:
- `internal/repository`: Mock repository tests
- `internal/service`: Service layer tests with context support
- Integration tests: Full stack testing

## Conclusion

The generic BaseRepository provides:
- ✅ Type safety without runtime overhead
- ✅ Cleaner, more maintainable code
- ✅ Better developer experience
- ✅ Easier to extend for new entity types
- ✅ Modern Go idioms (Go 1.18+)

This refactoring maintains backward compatibility at the business logic level while modernizing the internal implementation.

