package user_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"quizizz.com/internal/api/response"
	"quizizz.com/internal/domain"
	"quizizz.com/internal/service"
	"quizizz.com/internal/testutil/integration"
)

func TestIntegration_UserAPI(t *testing.T) {
	// Test creating a new user
	t.Run("Create and retrieve user", func(t *testing.T) {
		// Setup test environment for this specific test
		env := integration.Setup(t)
		defer env.Cleanup()
		// Create a new user
		userJSON := `{
			"name": "Integration Test User",
			"email": "integration@example.com"
		}`

		// POST request to create user
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/users", strings.NewReader(userJSON))
		req.Header.Set("Content-Type", "application/json")
		env.Router.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusCreated, w.Code)

		// Parse response
		var createResp response.Response
		err := json.Unmarshal(w.Body.Bytes(), &createResp)
		require.NoError(t, err)

		// Extract user ID
		userData, ok := createResp.Data.(map[string]interface{})
		require.True(t, ok)
		userID, ok := userData["id"].(string)
		require.True(t, ok)
		assert.NotEmpty(t, userID)

		// GET request to retrieve the user
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/api/v1/users/"+userID, nil)
		env.Router.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusOK, w.Code)

		// Parse response
		var getResp response.Response
		err = json.Unmarshal(w.Body.Bytes(), &getResp)
		require.NoError(t, err)

		// Check user data
		retrievedUser, ok := getResp.Data.(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, userID, retrievedUser["id"])
		assert.Equal(t, "Integration Test User", retrievedUser["name"])
		assert.Equal(t, "integration@example.com", retrievedUser["email"])
	})

	// Test listing users
	t.Run("List users", func(t *testing.T) {
		// Setup test environment for this specific test
		env := integration.Setup(t)
		defer env.Cleanup()

		// First, create a user to ensure we have at least one
		user := domain.NewUser("List Test User", "list@example.com")
		err := env.UserService.Create(context.Background(), user)
		require.NoError(t, err)

		// GET request to list users
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/users", nil)
		env.Router.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusOK, w.Code)

		// Parse response
		var listResp response.Response
		err = json.Unmarshal(w.Body.Bytes(), &listResp)
		require.NoError(t, err)

		// Check that we have users
		data, ok := listResp.Data.(map[string]interface{})
		require.True(t, ok)

		users, ok := data["users"].([]interface{})
		require.True(t, ok)
		assert.NotEmpty(t, users)

		count, ok := data["count"].(float64)
		require.True(t, ok)
		assert.True(t, count > 0)
	})

	// Test updating a user
	t.Run("Update user", func(t *testing.T) {
		// Setup test environment for this specific test
		env := integration.Setup(t)
		defer env.Cleanup()

		// First, create a user
		user := domain.NewUser("Update Test User", "update@example.com")
		err := env.UserService.Create(context.Background(), user)
		require.NoError(t, err)

		// Update JSON
		updateJSON := `{
			"name": "Updated User",
			"email": "updated@example.com"
		}`

		// PUT request to update user
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/users/"+user.ID, strings.NewReader(updateJSON))
		req.Header.Set("Content-Type", "application/json")
		env.Router.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusOK, w.Code)

		// Get the updated user
		updatedUser, err := env.UserService.GetByID(context.Background(), user.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated User", updatedUser.Name)
		assert.Equal(t, "updated@example.com", updatedUser.Email)
	})

	// Test deleting a user
	t.Run("Delete user", func(t *testing.T) {
		// Setup test environment for this specific test
		env := integration.Setup(t)
		defer env.Cleanup()

		// First, create a user
		user := domain.NewUser("Delete Test User", "delete@example.com")
		err := env.UserService.Create(context.Background(), user)
		require.NoError(t, err)

		// DELETE request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/users/"+user.ID, nil)
		env.Router.ServeHTTP(w, req)

		// Check status code
		assert.Equal(t, http.StatusNoContent, w.Code)

		// Verify user is deleted
		deletedUser, err := env.UserService.GetByID(context.Background(), user.ID)
		assert.Equal(t, service.ErrUserNotFound, err)
		assert.Nil(t, deletedUser)
	})
}
