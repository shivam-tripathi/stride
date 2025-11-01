package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"quizizz.com/internal/domain"
)

// MockUserRepo is a mock implementation of the UserRepository for testing
type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepo) List(ctx context.Context) ([]*domain.User, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *MockUserRepo) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestUserService_GetByID(t *testing.T) {
	// Create test context
	ctx := context.Background()

	t.Run("Valid user", func(t *testing.T) {
		// Setup mock
		mockRepo := new(MockUserRepo)
		user := &domain.User{
			ID:        "test-id",
			Name:      "Test User",
			Email:     "test@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Set expectations
		mockRepo.On("GetByID", ctx, "test-id").Return(user, nil)

		// Create service with mock
		service := NewUserService(mockRepo)

		// Call service
		result, err := service.GetByID(ctx, "test-id")

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, user, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Empty ID", func(t *testing.T) {
		// Setup mock
		mockRepo := new(MockUserRepo)

		// Create service with mock
		service := NewUserService(mockRepo)

		// Call service
		result, err := service.GetByID(ctx, "")

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidUser, err)
		assert.Nil(t, result)
		mockRepo.AssertNotCalled(t, "GetByID")
	})

	t.Run("User not found", func(t *testing.T) {
		// Setup mock
		mockRepo := new(MockUserRepo)

		// Set expectations
		mockRepo.On("GetByID", ctx, "non-existent").Return(nil, nil)

		// Create service with mock
		service := NewUserService(mockRepo)

		// Call service
		result, err := service.GetByID(ctx, "non-existent")

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, ErrUserNotFound, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository error", func(t *testing.T) {
		// Setup mock
		mockRepo := new(MockUserRepo)
		repoErr := errors.New("repository error")

		// Set expectations
		mockRepo.On("GetByID", ctx, "test-id").Return(nil, repoErr)

		// Create service with mock
		service := NewUserService(mockRepo)

		// Call service
		result, err := service.GetByID(ctx, "test-id")

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, repoErr, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_List(t *testing.T) {
	// Create test context
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// Setup mock
		mockRepo := new(MockUserRepo)
		users := []*domain.User{
			{
				ID:        "test-id-1",
				Name:      "Test User 1",
				Email:     "test1@example.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        "test-id-2",
				Name:      "Test User 2",
				Email:     "test2@example.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		// Set expectations
		mockRepo.On("List", ctx).Return(users, nil)

		// Create service with mock
		service := NewUserService(mockRepo)

		// Call service
		result, err := service.List(ctx)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, users, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Empty list", func(t *testing.T) {
		// Setup mock
		mockRepo := new(MockUserRepo)
		users := []*domain.User{}

		// Set expectations
		mockRepo.On("List", ctx).Return(users, nil)

		// Create service with mock
		service := NewUserService(mockRepo)

		// Call service
		result, err := service.List(ctx)

		// Assertions
		assert.NoError(t, err)
		assert.Empty(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository error", func(t *testing.T) {
		// Setup mock
		mockRepo := new(MockUserRepo)
		repoErr := errors.New("repository error")

		// Set expectations
		mockRepo.On("List", ctx).Return(nil, repoErr)

		// Create service with mock
		service := NewUserService(mockRepo)

		// Call service
		result, err := service.List(ctx)

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, repoErr, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Create(t *testing.T) {
	// Create test context
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// Setup mock
		mockRepo := new(MockUserRepo)
		user := &domain.User{
			ID:        "test-id",
			Name:      "Test User",
			Email:     "test@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Set expectations
		mockRepo.On("Create", ctx, user).Return(nil)

		// Create service with mock
		service := NewUserService(mockRepo)

		// Call service
		err := service.Create(ctx, user)

		// Assertions
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Missing name", func(t *testing.T) {
		// Setup mock
		mockRepo := new(MockUserRepo)
		user := &domain.User{
			ID:        "test-id",
			Email:     "test@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Create service with mock
		service := NewUserService(mockRepo)

		// Call service
		err := service.Create(ctx, user)

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidUser, err)
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Missing email", func(t *testing.T) {
		// Setup mock
		mockRepo := new(MockUserRepo)
		user := &domain.User{
			ID:        "test-id",
			Name:      "Test User",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Create service with mock
		service := NewUserService(mockRepo)

		// Call service
		err := service.Create(ctx, user)

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidUser, err)
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Repository error", func(t *testing.T) {
		// Setup mock
		mockRepo := new(MockUserRepo)
		user := &domain.User{
			ID:        "test-id",
			Name:      "Test User",
			Email:     "test@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		repoErr := errors.New("repository error")

		// Set expectations
		mockRepo.On("Create", ctx, user).Return(repoErr)

		// Create service with mock
		service := NewUserService(mockRepo)

		// Call service
		err := service.Create(ctx, user)

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, repoErr, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Update(t *testing.T) {
	// Create test context
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// Setup mock
		mockRepo := new(MockUserRepo)
		user := &domain.User{
			ID:        "test-id",
			Name:      "Updated User",
			Email:     "updated@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Set expectations
		mockRepo.On("GetByID", ctx, "test-id").Return(user, nil)
		mockRepo.On("Update", ctx, user).Return(nil)

		// Create service with mock
		service := NewUserService(mockRepo)

		// Call service
		err := service.Update(ctx, user)

		// Assertions
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Empty ID", func(t *testing.T) {
		// Setup mock
		mockRepo := new(MockUserRepo)
		user := &domain.User{
			Name:      "Updated User",
			Email:     "updated@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Create service with mock
		service := NewUserService(mockRepo)

		// Call service
		err := service.Update(ctx, user)

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidUser, err)
		mockRepo.AssertNotCalled(t, "GetByID")
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("User not found", func(t *testing.T) {
		// Setup mock
		mockRepo := new(MockUserRepo)
		user := &domain.User{
			ID:        "test-id",
			Name:      "Updated User",
			Email:     "updated@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Set expectations
		mockRepo.On("GetByID", ctx, "test-id").Return(nil, nil)

		// Create service with mock
		service := NewUserService(mockRepo)

		// Call service
		err := service.Update(ctx, user)

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, ErrUserNotFound, err)
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("Repository error during get", func(t *testing.T) {
		// Setup mock
		mockRepo := new(MockUserRepo)
		user := &domain.User{
			ID:        "test-id",
			Name:      "Updated User",
			Email:     "updated@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		repoErr := errors.New("repository error")

		// Set expectations
		mockRepo.On("GetByID", ctx, "test-id").Return(nil, repoErr)

		// Create service with mock
		service := NewUserService(mockRepo)

		// Call service
		err := service.Update(ctx, user)

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, repoErr, err)
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("Repository error during update", func(t *testing.T) {
		// Setup mock
		mockRepo := new(MockUserRepo)
		user := &domain.User{
			ID:        "test-id",
			Name:      "Updated User",
			Email:     "updated@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		repoErr := errors.New("repository error")

		// Set expectations
		mockRepo.On("GetByID", ctx, "test-id").Return(user, nil)
		mockRepo.On("Update", ctx, user).Return(repoErr)

		// Create service with mock
		service := NewUserService(mockRepo)

		// Call service
		err := service.Update(ctx, user)

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, repoErr, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Delete(t *testing.T) {
	// Create test context
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// Setup mock
		mockRepo := new(MockUserRepo)
		user := &domain.User{
			ID:        "test-id",
			Name:      "Test User",
			Email:     "test@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Set expectations
		mockRepo.On("GetByID", ctx, "test-id").Return(user, nil)
		mockRepo.On("Delete", ctx, "test-id").Return(nil)

		// Create service with mock
		service := NewUserService(mockRepo)

		// Call service
		err := service.Delete(ctx, "test-id")

		// Assertions
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Empty ID", func(t *testing.T) {
		// Setup mock
		mockRepo := new(MockUserRepo)

		// Create service with mock
		service := NewUserService(mockRepo)

		// Call service
		err := service.Delete(ctx, "")

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidUser, err)
		mockRepo.AssertNotCalled(t, "GetByID")
		mockRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("User not found", func(t *testing.T) {
		// Setup mock
		mockRepo := new(MockUserRepo)

		// Set expectations
		mockRepo.On("GetByID", ctx, "test-id").Return(nil, nil)

		// Create service with mock
		service := NewUserService(mockRepo)

		// Call service
		err := service.Delete(ctx, "test-id")

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, ErrUserNotFound, err)
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("Repository error during get", func(t *testing.T) {
		// Setup mock
		mockRepo := new(MockUserRepo)
		repoErr := errors.New("repository error")

		// Set expectations
		mockRepo.On("GetByID", ctx, "test-id").Return(nil, repoErr)

		// Create service with mock
		service := NewUserService(mockRepo)

		// Call service
		err := service.Delete(ctx, "test-id")

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, repoErr, err)
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("Repository error during delete", func(t *testing.T) {
		// Setup mock
		mockRepo := new(MockUserRepo)
		user := &domain.User{
			ID:        "test-id",
			Name:      "Test User",
			Email:     "test@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		repoErr := errors.New("repository error")

		// Set expectations
		mockRepo.On("GetByID", ctx, "test-id").Return(user, nil)
		mockRepo.On("Delete", ctx, "test-id").Return(repoErr)

		// Create service with mock
		service := NewUserService(mockRepo)

		// Call service
		err := service.Delete(ctx, "test-id")

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, repoErr, err)
		mockRepo.AssertExpectations(t)
	})
}
