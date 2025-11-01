package service

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"quizizz.com/internal/domain"
	"quizizz.com/internal/logger"
	"quizizz.com/internal/repository"
)

// Common errors
var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidUser  = errors.New("invalid user data")
)

// UserService defines the interface for user-related business logic
type UserService interface {
	GetByID(ctx context.Context, id string) (*domain.User, error)
	List(ctx context.Context) ([]*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error
}

// userService implements the UserService interface
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new UserService
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// GetByID retrieves a user by ID
func (s *userService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	logger.Debug("Getting user by ID", zap.String("userId", id))

	if id == "" {
		return nil, ErrInvalidUser
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to get user", zap.String("userId", id), zap.Error(err))
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// List retrieves all users
func (s *userService) List(ctx context.Context) ([]*domain.User, error) {
	logger.Debug("Listing users")

	users, err := s.userRepo.List(ctx)
	if err != nil {
		logger.Error("Failed to list users", zap.Error(err))
		return nil, err
	}

	return users, nil
}

// Create creates a new user
func (s *userService) Create(ctx context.Context, user *domain.User) error {
	logger.Debug("Creating user", zap.String("userName", user.Name))

	if user.Name == "" || user.Email == "" {
		return ErrInvalidUser
	}

	err := s.userRepo.Create(ctx, user)
	if err != nil {
		logger.Error("Failed to create user", zap.Error(err))
		return err
	}

	logger.Info("User created", zap.String("userId", user.ID), zap.String("userName", user.Name))
	return nil
}

// Update updates an existing user
func (s *userService) Update(ctx context.Context, user *domain.User) error {
	logger.Debug("Updating user", zap.String("userId", user.ID))

	if user.ID == "" {
		return ErrInvalidUser
	}

	// Check if user exists
	existingUser, err := s.userRepo.GetByID(ctx, user.ID)
	if err != nil {
		logger.Error("Failed to get user for update", zap.String("userId", user.ID), zap.Error(err))
		return err
	}

	if existingUser == nil {
		return ErrUserNotFound
	}

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		logger.Error("Failed to update user", zap.String("userId", user.ID), zap.Error(err))
		return err
	}

	logger.Info("User updated", zap.String("userId", user.ID))
	return nil
}

// Delete deletes a user
func (s *userService) Delete(ctx context.Context, id string) error {
	logger.Debug("Deleting user", zap.String("userId", id))

	if id == "" {
		return ErrInvalidUser
	}

	// Check if user exists
	existingUser, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to get user for deletion", zap.String("userId", id), zap.Error(err))
		return err
	}

	if existingUser == nil {
		return ErrUserNotFound
	}

	err = s.userRepo.Delete(ctx, id)
	if err != nil {
		logger.Error("Failed to delete user", zap.String("userId", id), zap.Error(err))
		return err
	}

	logger.Info("User deleted", zap.String("userId", id))
	return nil
}
