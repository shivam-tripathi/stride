package repository

import (
	"errors"
	"sync"

	"quizizz.com/internal/domain"
)

// Common errors
var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)

// MockUserRepository is an in-memory implementation of UserRepository for testing
type MockUserRepository struct {
	users map[string]*domain.User
	mutex sync.RWMutex
}

// NewMockUserRepository creates a new MockUserRepository
func NewMockUserRepository() UserRepository {
	return &MockUserRepository{
		users: make(map[string]*domain.User),
	}
}

// GetByID returns a user by ID
func (r *MockUserRepository) GetByID(id string) (*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, nil // Return nil without error to indicate user not found
	}

	return user, nil
}

// List returns all users
func (r *MockUserRepository) List() ([]*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	users := make([]*domain.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	return users, nil
}

// Create adds a new user
func (r *MockUserRepository) Create(user *domain.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check if user already exists
	if _, exists := r.users[user.ID]; exists {
		return ErrUserExists
	}

	// Make a copy to avoid external modifications
	userCopy := *user
	r.users[user.ID] = &userCopy

	return nil
}

// Update updates an existing user
func (r *MockUserRepository) Update(user *domain.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check if user exists
	if _, exists := r.users[user.ID]; !exists {
		return ErrUserNotFound
	}

	// Make a copy to avoid external modifications
	userCopy := *user
	r.users[user.ID] = &userCopy

	return nil
}

// Delete removes a user
func (r *MockUserRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check if user exists
	if _, exists := r.users[id]; !exists {
		return ErrUserNotFound
	}

	delete(r.users, id)

	return nil
}
