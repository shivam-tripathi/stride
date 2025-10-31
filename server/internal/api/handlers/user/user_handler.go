// Package user provides user-related handlers
package user

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"quizizz.com/internal/api/handlers"
	"quizizz.com/internal/api/response"
	"quizizz.com/internal/domain"
	"quizizz.com/internal/errors"
	"quizizz.com/internal/service"
)

// User represents a user in the API
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
}

// Handler handles user-related requests
type Handler struct {
	*handlers.BaseHandler
	userService service.UserService
}

// NewHandler creates a new user handler
func NewHandler(base *handlers.BaseHandler, userService service.UserService) *Handler {
	return &Handler{
		BaseHandler: base,
		userService: userService,
	}
}

// ListUsers returns a list of users
func (h *Handler) ListUsers(c *gin.Context) {
	logger := h.GetRequestLogger(c)
	logger.Debug("Listing users")

	// Use service to get users
	domainUsers, err := h.userService.List(context.Background())
	if err != nil {
		logger.Error("Failed to list users", zap.Error(err))
		response.InternalServerError(c, "Failed to list users")
		return
	}

	// Convert domain users to API users
	users := make([]User, 0, len(domainUsers))
	for _, domainUser := range domainUsers {
		users = append(users, User{
			ID:    domainUser.ID,
			Name:  domainUser.Name,
			Email: domainUser.Email,
		})
	}

	response.Success(c, gin.H{
		"users": users,
		"count": len(users),
	})
}

// GetUser returns a user by ID
func (h *Handler) GetUser(c *gin.Context) {
	id := c.Param("id")
	logger := h.GetRequestLogger(c).With(zap.String("userId", id))
	logger.Debug("Getting user by ID")

	// Example of error handling
	if id == "" {
		logger.Warn("User ID is empty")
		response.BadRequest(c, "User ID is required")
		return
	}

	// Use service to get user
	domainUser, err := h.userService.GetByID(context.Background(), id)
	if err != nil {
		// Handle different types of errors
		if err == service.ErrUserNotFound {
			logger.Warn("User not found")
			response.NotFound(c, "User not found")
			return
		}
		logger.Error("Failed to get user", zap.Error(err))
		response.InternalServerError(c, "Failed to get user")
		return
	}

	// Convert domain user to API user
	user := User{
		ID:    domainUser.ID,
		Name:  domainUser.Name,
		Email: domainUser.Email,
	}

	response.Success(c, user)
}

// CreateUser creates a new user
func (h *Handler) CreateUser(c *gin.Context) {
	logger := h.GetRequestLogger(c)
	logger.Debug("Creating new user")

	var userRequest User
	if !h.ShouldBindJSON(c, &userRequest) {
		logger.Warn("Invalid request body")
		response.BadRequest(c, "Invalid request body")
		return
	}

	// Validate user input
	if userRequest.Name == "" {
		err := &errors.AppError{
			StatusCode: http.StatusBadRequest,
			Message:    "Name is required",
		}
		err.WithContext("field", "name")
		response.Fail(c, err)
		return
	}

	// Convert API user to domain user
	domainUser := domain.NewUser(userRequest.Name, userRequest.Email)

	// Use service to create user
	err := h.userService.Create(context.Background(), domainUser)
	if err != nil {
		logger.Error("Failed to create user", zap.Error(err))
		response.InternalServerError(c, "Failed to create user")
		return
	}

	// Return created user
	userRequest.ID = domainUser.ID
	logger.Info("User created", zap.String("userId", userRequest.ID))
	response.Created(c, userRequest)
}

// UpdateUser updates an existing user
func (h *Handler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	logger := h.GetRequestLogger(c).With(zap.String("userId", id))
	logger.Debug("Updating user")

	// Example of error handling
	if id == "" {
		logger.Warn("User ID is empty")
		response.BadRequest(c, "User ID is required")
		return
	}

	var userRequest User
	if !h.ShouldBindJSON(c, &userRequest) {
		logger.Warn("Invalid request body")
		response.BadRequest(c, "Invalid request body")
		return
	}

	// Set the ID from the path parameter
	userRequest.ID = id

	// Get existing user
	existingUser, err := h.userService.GetByID(context.Background(), id)
	if err != nil {
		if err == service.ErrUserNotFound {
			logger.Warn("User not found for update")
			response.NotFound(c, "User not found")
			return
		}
		logger.Error("Failed to get user for update", zap.Error(err))
		response.InternalServerError(c, "Failed to update user")
		return
	}

	// Update user fields
	existingUser.Name = userRequest.Name
	if userRequest.Email != "" {
		existingUser.Email = userRequest.Email
	}

	// Use service to update user
	err = h.userService.Update(context.Background(), existingUser)
	if err != nil {
		logger.Error("Failed to update user", zap.Error(err))
		response.InternalServerError(c, "Failed to update user")
		return
	}

	logger.Info("User updated", zap.String("userId", userRequest.ID))
	response.Success(c, userRequest)
}

// DeleteUser deletes a user
func (h *Handler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	logger := h.GetRequestLogger(c).With(zap.String("userId", id))
	logger.Debug("Deleting user")

	// Example of error handling
	if id == "" {
		logger.Warn("User ID is empty")
		response.BadRequest(c, "User ID is required")
		return
	}

	// Use service to delete user
	err := h.userService.Delete(context.Background(), id)
	if err != nil {
		if err == service.ErrUserNotFound {
			logger.Warn("User not found for deletion")
			response.NotFound(c, "User not found")
			return
		}
		logger.Error("Failed to delete user", zap.Error(err))
		response.InternalServerError(c, "Failed to delete user")
		return
	}

	logger.Info("User deleted", zap.String("userId", id))
	response.NoContent(c)
}
