package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"quizizz.com/internal/api/handlers"
	"quizizz.com/internal/api/response"
	"quizizz.com/internal/domain"
	"quizizz.com/internal/service"
)

// Mock implementations
type MockAppService struct {
	mock.Mock
}

func (m *MockAppService) GetPingMessage() string {
	args := m.Called()
	return args.String(0)
}

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) List(ctx context.Context) ([]*domain.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *MockUserService) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserService) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserService) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Setup test function
func setupUserHandler() (*Handler, *MockAppService, *MockUserService) {
	gin.SetMode(gin.TestMode)

	mockAppService := new(MockAppService)
	mockUserService := new(MockUserService)

	baseHandler := handlers.NewBaseHandler(mockAppService)
	handler := NewHandler(baseHandler, mockUserService)

	return handler, mockAppService, mockUserService
}

// Helper functions
func createTestRouter(handler *Handler) *gin.Engine {
	router := gin.New()

	// Add the logger for testing
	router.Use(func(c *gin.Context) {
		c.Set("requestID", "test-request-id")
		c.Next()
	})

	// Setup user routes
	users := router.Group("/api/v1/users")
	{
		users.GET("", handler.ListUsers)
		users.POST("", handler.CreateUser)
		users.GET("/:id", handler.GetUser)
		users.PUT("/:id", handler.UpdateUser)
		users.DELETE("/:id", handler.DeleteUser)
	}

	return router
}

// Test function to parse response body
func parseResponse(t *testing.T, w *httptest.ResponseRecorder, target interface{}) {
	require.NotNil(t, w.Body)
	err := json.Unmarshal(w.Body.Bytes(), target)
	require.NoError(t, err, "Failed to parse response body")
}

// Tests
func TestHandler_ListUsers(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Setup
		handler, _, mockUserService := setupUserHandler()
		router := createTestRouter(handler)

		// Mock data
		domainUsers := []*domain.User{
			{
				ID:    "user-1",
				Name:  "User 1",
				Email: "user1@example.com",
			},
			{
				ID:    "user-2",
				Name:  "User 2",
				Email: "user2@example.com",
			},
		}

		// Set expectations
		mockUserService.On("List", mock.Anything).Return(domainUsers, nil)

		// Perform request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/users", nil)
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		// Parse response
		var responseObj response.Response
		parseResponse(t, w, &responseObj)

		// Check response structure
		assert.True(t, responseObj.Success)
		assert.Nil(t, responseObj.Error)

		// Convert and check data
		data, ok := responseObj.Data.(map[string]interface{})
		require.True(t, ok, "Data is not a map")

		users, ok := data["users"].([]interface{})
		require.True(t, ok, "Users is not an array")
		assert.Len(t, users, 2)

		count, ok := data["count"].(float64)
		require.True(t, ok, "Count is not a number")
		assert.Equal(t, float64(2), count)

		// Verify mock expectations
		mockUserService.AssertExpectations(t)
	})

	t.Run("Service error", func(t *testing.T) {
		// Setup
		handler, _, mockUserService := setupUserHandler()
		router := createTestRouter(handler)

		// Set expectations
		mockUserService.On("List", mock.Anything).Return(nil, errors.New("service error"))

		// Perform request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/users", nil)
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// Parse response
		var responseObj response.Response
		parseResponse(t, w, &responseObj)

		// Check response structure
		assert.False(t, responseObj.Success)
		assert.NotNil(t, responseObj.Error)
		assert.Equal(t, "Failed to list users", responseObj.Error.Message)

		// Verify mock expectations
		mockUserService.AssertExpectations(t)
	})
}

func TestHandler_GetUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Setup
		handler, _, mockUserService := setupUserHandler()
		router := createTestRouter(handler)

		// Mock data
		user := &domain.User{
			ID:    "user-1",
			Name:  "User 1",
			Email: "user1@example.com",
		}

		// Set expectations
		mockUserService.On("GetByID", mock.Anything, "user-1").Return(user, nil)

		// Perform request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/users/user-1", nil)
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		// Parse response
		var responseObj response.Response
		parseResponse(t, w, &responseObj)

		// Check response structure
		assert.True(t, responseObj.Success)
		assert.Nil(t, responseObj.Error)

		// Convert and check data
		userData, ok := responseObj.Data.(map[string]interface{})
		require.True(t, ok, "Data is not a map")

		assert.Equal(t, "user-1", userData["id"])
		assert.Equal(t, "User 1", userData["name"])
		assert.Equal(t, "user1@example.com", userData["email"])

		// Verify mock expectations
		mockUserService.AssertExpectations(t)
	})

	t.Run("User not found", func(t *testing.T) {
		// Setup
		handler, _, mockUserService := setupUserHandler()
		router := createTestRouter(handler)

		// Set expectations
		mockUserService.On("GetByID", mock.Anything, "non-existent").Return(nil, service.ErrUserNotFound)

		// Perform request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/users/non-existent", nil)
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusNotFound, w.Code)

		// Parse response
		var responseObj response.Response
		parseResponse(t, w, &responseObj)

		// Check response structure
		assert.False(t, responseObj.Success)
		assert.NotNil(t, responseObj.Error)
		assert.Equal(t, "User not found", responseObj.Error.Message)

		// Verify mock expectations
		mockUserService.AssertExpectations(t)
	})

	t.Run("Service error", func(t *testing.T) {
		// Setup
		handler, _, mockUserService := setupUserHandler()
		router := createTestRouter(handler)

		// Set expectations
		mockUserService.On("GetByID", mock.Anything, "user-1").Return(nil, errors.New("service error"))

		// Perform request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/users/user-1", nil)
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// Parse response
		var responseObj response.Response
		parseResponse(t, w, &responseObj)

		// Check response structure
		assert.False(t, responseObj.Success)
		assert.NotNil(t, responseObj.Error)
		assert.Equal(t, "Failed to get user", responseObj.Error.Message)

		// Verify mock expectations
		mockUserService.AssertExpectations(t)
	})
}

func TestHandler_CreateUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Setup
		handler, _, mockUserService := setupUserHandler()
		router := createTestRouter(handler)

		// Mock behavior - service will create a user
		mockUserService.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).
			Run(func(args mock.Arguments) {
				user := args.Get(1).(*domain.User)
				user.ID = "new-user-id" // Set ID in the domain user
			}).
			Return(nil)

		// Create request body
		requestBody := `{"name":"New User","email":"newuser@example.com"}`

		// Perform request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/users", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusCreated, w.Code)

		// Parse response
		var responseObj response.Response
		parseResponse(t, w, &responseObj)

		// Check response structure
		assert.True(t, responseObj.Success)
		assert.Nil(t, responseObj.Error)

		// Convert and check data
		userData, ok := responseObj.Data.(map[string]interface{})
		require.True(t, ok, "Data is not a map")

		assert.Equal(t, "new-user-id", userData["id"])
		assert.Equal(t, "New User", userData["name"])
		assert.Equal(t, "newuser@example.com", userData["email"])

		// Verify mock expectations
		mockUserService.AssertExpectations(t)
	})

	t.Run("Invalid request body", func(t *testing.T) {
		// Setup
		handler, _, _ := setupUserHandler()
		router := createTestRouter(handler)

		// Create invalid request body
		requestBody := `{invalid_json`

		// Perform request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/users", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Parse response
		var responseObj response.Response
		parseResponse(t, w, &responseObj)

		// Check response structure
		assert.False(t, responseObj.Success)
		assert.NotNil(t, responseObj.Error)
		assert.Equal(t, "Invalid request body", responseObj.Error.Message)
	})

	t.Run("Missing name", func(t *testing.T) {
		// Setup
		handler, _, _ := setupUserHandler()
		router := createTestRouter(handler)

		// Create request body with missing name
		requestBody := `{"email":"newuser@example.com"}`

		// Perform request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/users", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Parse response
		var responseObj response.Response
		parseResponse(t, w, &responseObj)

		// Check response structure
		assert.False(t, responseObj.Success)
		assert.NotNil(t, responseObj.Error)
		assert.Equal(t, "Name is required", responseObj.Error.Message)
	})

	t.Run("Service error", func(t *testing.T) {
		// Setup
		handler, _, mockUserService := setupUserHandler()
		router := createTestRouter(handler)

		// Mock behavior - service will return an error
		mockUserService.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).
			Return(errors.New("service error"))

		// Create request body
		requestBody := `{"name":"New User","email":"newuser@example.com"}`

		// Perform request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/users", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// Parse response
		var responseObj response.Response
		parseResponse(t, w, &responseObj)

		// Check response structure
		assert.False(t, responseObj.Success)
		assert.NotNil(t, responseObj.Error)
		assert.Equal(t, "Failed to create user", responseObj.Error.Message)

		// Verify mock expectations
		mockUserService.AssertExpectations(t)
	})
}

func TestHandler_UpdateUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Setup
		handler, _, mockUserService := setupUserHandler()
		router := createTestRouter(handler)

		// Mock data
		existingUser := &domain.User{
			ID:    "user-1",
			Name:  "Original Name",
			Email: "original@example.com",
		}

		// Set expectations for getting and updating the user
		mockUserService.On("GetByID", mock.Anything, "user-1").Return(existingUser, nil)
		mockUserService.On("Update", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

		// Create request body
		requestBody := `{"name":"Updated Name","email":"updated@example.com"}`

		// Perform request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/users/user-1", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		// Parse response
		var responseObj response.Response
		parseResponse(t, w, &responseObj)

		// Check response structure
		assert.True(t, responseObj.Success)
		assert.Nil(t, responseObj.Error)

		// Convert and check data
		userData, ok := responseObj.Data.(map[string]interface{})
		require.True(t, ok, "Data is not a map")

		assert.Equal(t, "user-1", userData["id"])
		assert.Equal(t, "Updated Name", userData["name"])
		assert.Equal(t, "updated@example.com", userData["email"])

		// Verify mock expectations
		mockUserService.AssertExpectations(t)
	})

	t.Run("User not found", func(t *testing.T) {
		// Setup
		handler, _, mockUserService := setupUserHandler()
		router := createTestRouter(handler)

		// Set expectations
		mockUserService.On("GetByID", mock.Anything, "non-existent").Return(nil, service.ErrUserNotFound)

		// Create request body
		requestBody := `{"name":"Updated Name","email":"updated@example.com"}`

		// Perform request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/users/non-existent", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusNotFound, w.Code)

		// Parse response
		var responseObj response.Response
		parseResponse(t, w, &responseObj)

		// Check response structure
		assert.False(t, responseObj.Success)
		assert.NotNil(t, responseObj.Error)
		assert.Equal(t, "User not found", responseObj.Error.Message)

		// Verify mock expectations
		mockUserService.AssertExpectations(t)
	})
}

func TestHandler_DeleteUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Setup
		handler, _, mockUserService := setupUserHandler()
		router := createTestRouter(handler)

		// Set expectations
		mockUserService.On("Delete", mock.Anything, "user-1").Return(nil)

		// Perform request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/users/user-1", nil)
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusNoContent, w.Code)

		// Verify mock expectations
		mockUserService.AssertExpectations(t)
	})

	t.Run("User not found", func(t *testing.T) {
		// Setup
		handler, _, mockUserService := setupUserHandler()
		router := createTestRouter(handler)

		// Set expectations
		mockUserService.On("Delete", mock.Anything, "non-existent").Return(service.ErrUserNotFound)

		// Perform request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/users/non-existent", nil)
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusNotFound, w.Code)

		// Parse response
		var responseObj response.Response
		parseResponse(t, w, &responseObj)

		// Check response structure
		assert.False(t, responseObj.Success)
		assert.NotNil(t, responseObj.Error)
		assert.Equal(t, "User not found", responseObj.Error.Message)

		// Verify mock expectations
		mockUserService.AssertExpectations(t)
	})
}
