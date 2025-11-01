package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"quizizz.com/internal/domain"
	"quizizz.com/internal/resources"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	GetByID(ctx context.Context, id string) (*domain.User, error)
	List(ctx context.Context) ([]*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error
}

// userRepositoryImpl is the MongoDB implementation of UserRepository
type userRepositoryImpl struct {
	*BaseRepository[userDocument]
	db *resources.DB
}

// userDocument represents the MongoDB document structure for users
type userDocument struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Email     string             `bson:"email"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db resources.DBResource) UserRepository {
	dbInstance := db.(*resources.DB)
	collection := dbInstance.Collection("users")

	return &userRepositoryImpl{
		BaseRepository: NewBaseRepositoryWithConfig[userDocument](BaseRepositoryConfig{
			Collection: collection,
			EntityName: "user",
		}),
		db: dbInstance,
	}
}

// GetByID returns a user by ID
func (r *userRepositoryImpl) GetByID(ctx context.Context, id string) (*domain.User, error) {
	doc, err := r.FindByID(ctx, id)
	if err != nil {
		if err == ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	return toUser(doc), nil
}

// List returns all users
func (r *userRepositoryImpl) List(ctx context.Context) ([]*domain.User, error) {
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})

	docs, err := r.FindAll(ctx, opts)
	if err != nil {
		return nil, err
	}

	return toUsers(docs), nil
}

// Create adds a new user
func (r *userRepositoryImpl) Create(ctx context.Context, user *domain.User) error {
	if exists, _ := r.Exists(ctx, bson.M{"email": user.Email}); exists {
		return ErrUserExists
	}

	doc := toDocument(user)
	doc.CreatedAt = time.Now()
	doc.UpdatedAt = time.Now()

	id, err := r.InsertOne(ctx, &doc)
	if err != nil {
		return err
	}

	user.ID = id
	user.CreatedAt = doc.CreatedAt
	user.UpdatedAt = doc.UpdatedAt

	return nil
}

// Update updates an existing user
func (r *userRepositoryImpl) Update(ctx context.Context, user *domain.User) error {
	update := bson.M{
		"name":      user.Name,
		"email":     user.Email,
		"updatedAt": time.Now(),
	}

	if err := r.UpdateByID(ctx, user.ID, update); err != nil {
		if err == ErrNotFound {
			return ErrUserNotFound
		}
		return err
	}

	user.UpdatedAt = time.Now()
	return nil
}

// Delete removes a user
func (r *userRepositoryImpl) Delete(ctx context.Context, id string) error {
	if err := r.DeleteByID(ctx, id); err != nil {
		if err == ErrNotFound {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}

// EnsureIndexes creates necessary indexes for the users collection
func (r *userRepositoryImpl) EnsureIndexes() error {
	ctx := context.Background()

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "createdAt", Value: -1}},
		},
	}

	return r.db.EnsureIndexes(ctx, "users", indexes)
}

// Conversion helpers

func toUser(doc *userDocument) *domain.User {
	return &domain.User{
		ID:        doc.ID.Hex(),
		Name:      doc.Name,
		Email:     doc.Email,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
	}
}

func toUsers(docs []userDocument) []*domain.User {
	users := make([]*domain.User, len(docs))
	for i := range docs {
		users[i] = toUser(&docs[i])
	}
	return users
}

func toDocument(user *domain.User) userDocument {
	doc := userDocument{
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if user.ID != "" {
		if objectID, err := primitive.ObjectIDFromHex(user.ID); err == nil {
			doc.ID = objectID
		}
	}

	return doc
}
