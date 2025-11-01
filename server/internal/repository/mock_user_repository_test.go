package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"quizizz.com/internal/domain"
)

func TestMockUserRepository_GetByID(t *testing.T) {
	// Setup
	repo := NewMockUserRepository()
	user := &domain.User{
		ID:        "test-id",
		Name:      "Test User",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Add the user to the repository
	err := repo.Create(context.Background(), user)
	require.NoError(t, err)

	// Test successful retrieval
	t.Run("Existing user", func(t *testing.T) {
		foundUser, err := repo.GetByID(context.Background(), "test-id")
		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, user.ID, foundUser.ID)
		assert.Equal(t, user.Name, foundUser.Name)
		assert.Equal(t, user.Email, foundUser.Email)
	})

	// Test user not found
	t.Run("Non-existent user", func(t *testing.T) {
		foundUser, err := repo.GetByID(context.Background(), "non-existent-id")
		assert.NoError(t, err) // No error, just nil user
		assert.Nil(t, foundUser)
	})
}

func TestMockUserRepository_List(t *testing.T) {
	// Setup
	repo := NewMockUserRepository()
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

	// Add users to the repository
	for _, user := range users {
		err := repo.Create(context.Background(), user)
		require.NoError(t, err)
	}

	// Test list users
	t.Run("List all users", func(t *testing.T) {
		foundUsers, err := repo.List(context.Background())
		assert.NoError(t, err)
		assert.Len(t, foundUsers, len(users))

		// Check that all users are present
		foundIDs := make(map[string]bool)
		for _, user := range foundUsers {
			foundIDs[user.ID] = true
		}

		for _, user := range users {
			assert.True(t, foundIDs[user.ID], "User with ID %s not found", user.ID)
		}
	})

	// Test empty repository
	t.Run("Empty repository", func(t *testing.T) {
		emptyRepo := NewMockUserRepository()
		foundUsers, err := emptyRepo.List(context.Background())
		assert.NoError(t, err)
		assert.Empty(t, foundUsers)
	})
}

func TestMockUserRepository_Create(t *testing.T) {
	// Setup
	repo := NewMockUserRepository()
	user := &domain.User{
		ID:        "test-id",
		Name:      "Test User",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Test successful creation
	t.Run("Create new user", func(t *testing.T) {
		err := repo.Create(context.Background(), user)
		assert.NoError(t, err)

		// Verify user was created
		foundUser, err := repo.GetByID(context.Background(), user.ID)
		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, user.ID, foundUser.ID)
	})

	// Test duplicate user
	t.Run("Create duplicate user", func(t *testing.T) {
		err := repo.Create(context.Background(), user)
		assert.Error(t, err)
		assert.Equal(t, ErrUserExists, err)
	})
}

func TestMockUserRepository_Update(t *testing.T) {
	// Setup
	repo := NewMockUserRepository()
	user := &domain.User{
		ID:        "test-id",
		Name:      "Test User",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Add the user to the repository
	err := repo.Create(context.Background(), user)
	require.NoError(t, err)

	// Test successful update
	t.Run("Update existing user", func(t *testing.T) {
		// Update user
		updatedUser := &domain.User{
			ID:        user.ID,
			Name:      "Updated Name",
			Email:     "updated@example.com",
			CreatedAt: user.CreatedAt,
			UpdatedAt: time.Now(),
		}

		err := repo.Update(context.Background(), updatedUser)
		assert.NoError(t, err)

		// Verify user was updated
		foundUser, err := repo.GetByID(context.Background(), user.ID)
		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, updatedUser.Name, foundUser.Name)
		assert.Equal(t, updatedUser.Email, foundUser.Email)
	})

	// Test update non-existent user
	t.Run("Update non-existent user", func(t *testing.T) {
		nonExistentUser := &domain.User{
			ID:        "non-existent-id",
			Name:      "Non-existent User",
			Email:     "nonexistent@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Update(context.Background(), nonExistentUser)
		assert.Error(t, err)
		assert.Equal(t, ErrUserNotFound, err)
	})
}

func TestMockUserRepository_Delete(t *testing.T) {
	// Setup
	repo := NewMockUserRepository()
	user := &domain.User{
		ID:        "test-id",
		Name:      "Test User",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Add the user to the repository
	err := repo.Create(context.Background(), user)
	require.NoError(t, err)

	// Test successful deletion
	t.Run("Delete existing user", func(t *testing.T) {
		err := repo.Delete(context.Background(), user.ID)
		assert.NoError(t, err)

		// Verify user was deleted
		foundUser, err := repo.GetByID(context.Background(), user.ID)
		assert.NoError(t, err)
		assert.Nil(t, foundUser)
	})

	// Test delete non-existent user
	t.Run("Delete non-existent user", func(t *testing.T) {
		err := repo.Delete(context.Background(), "non-existent-id")
		assert.Error(t, err)
		assert.Equal(t, ErrUserNotFound, err)
	})
}
