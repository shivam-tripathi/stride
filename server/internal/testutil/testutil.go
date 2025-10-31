// Package testutil provides utilities and helpers for testing
package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Setup initializes test environment
func Setup() {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
}

// CreateTestServer creates a test HTTP server using Gin
func CreateTestServer() *gin.Engine {
	return gin.New()
}

// MakeTestRequest is a helper function to make HTTP requests in tests
func MakeTestRequest(t *testing.T, router *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	// Create request body if provided
	var reqBody io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		require.NoError(t, err, "Failed to marshal request body")
		reqBody = bytes.NewBuffer(jsonBytes)
	}

	// Create HTTP request
	req, err := http.NewRequest(method, path, reqBody)
	require.NoError(t, err, "Failed to create request")
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(w, req)
	return w
}

// ParseResponse parses a JSON response into the provided target struct
func ParseResponse(t *testing.T, w *httptest.ResponseRecorder, target interface{}) {
	err := json.Unmarshal(w.Body.Bytes(), target)
	require.NoError(t, err, "Failed to parse response body")
}

// AssertStatusCode asserts the HTTP status code
func AssertStatusCode(t *testing.T, w *httptest.ResponseRecorder, expectedCode int) {
	assert.Equal(t, expectedCode, w.Code, "Expected status code %d but got %d", expectedCode, w.Code)
}

// AssertJSONContains asserts that the JSON response contains specific key-value pairs
func AssertJSONContains(t *testing.T, w *httptest.ResponseRecorder, expected map[string]interface{}) {
	var actual map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &actual)
	require.NoError(t, err, "Failed to parse JSON response")

	for key, value := range expected {
		assert.Contains(t, actual, key, "Response doesn't contain key: %s", key)
		assert.Equal(t, value, actual[key], "Value mismatch for key %s", key)
	}
}

// GetProjectRoot returns the absolute path to the project root
func GetProjectRoot() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(b), "../..")
}

// LoadFixture loads a test fixture file from the testdata directory
func LoadFixture(t *testing.T, name string) []byte {
	path := filepath.Join(GetProjectRoot(), "testdata", name)
	data, err := os.ReadFile(path)
	require.NoError(t, err, "Failed to load test fixture: %s", name)
	return data
}

// LoadFixtureInto loads a JSON fixture into the provided target struct
func LoadFixtureInto(t *testing.T, name string, target interface{}) {
	data := LoadFixture(t, name)
	err := json.Unmarshal(data, target)
	require.NoError(t, err, "Failed to unmarshal test fixture: %s", name)
}

// CreateTempFile creates a temporary file with the given content
func CreateTempFile(t *testing.T, content []byte) (string, func()) {
	tmpFile, err := os.CreateTemp("", "test-")
	require.NoError(t, err, "Failed to create temp file")

	_, err = tmpFile.Write(content)
	require.NoError(t, err, "Failed to write to temp file")

	err = tmpFile.Close()
	require.NoError(t, err, "Failed to close temp file")

	return tmpFile.Name(), func() {
		os.Remove(tmpFile.Name())
	}
}

// SkipIfShort skips a test if the -short flag is used
func SkipIfShort(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}
}

// PrintJSON pretty prints a JSON value for debugging
func PrintJSON(t *testing.T, v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	require.NoError(t, err, "Failed to marshal JSON")
	fmt.Println(string(b))
}
