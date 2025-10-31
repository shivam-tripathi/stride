# Testing Guide

This document provides a comprehensive guide to all the tests in the project and how to run them.

## Test Types

The project contains several types of tests:

1. **Unit Tests**: Test individual components in isolation using mocks
2. **Integration Tests**: Test interactions between multiple components using mock external dependencies
3. **Benchmark Tests**: Measure performance of specific operations
4. **Handler Tests**: Test API endpoints with mocked services
5. **End-to-End Tests**: (Future) Test complete workflows with real external dependencies

## Running Tests

### Basic Test Commands

```bash
# Run all tests
make test

# Run only unit tests
make test-unit

# Run only integration tests
make test-integration

# Run tests with race detection
make test-race

# Run tests and generate coverage report
make test-coverage

# Run only short tests (useful for quick feedback)
make test-short
```

### Running Specific Tests

```bash
# Run tests in a specific package
go test ./internal/service

# Run a specific test
go test -run=TestUserService_GetByID ./internal/service

# Run tests with verbose output
go test -v ./...
```

### Running Benchmark Tests

```bash
# Run all benchmarks
go test -bench=. ./...

# Run benchmarks in a specific package
go test -bench=. ./internal/service

# Run a specific benchmark
go test -bench=BenchmarkUserService_GetByID ./internal/service

# Run benchmarks with memory allocation statistics
go test -bench=. -benchmem ./internal/service

# Run benchmarks with less verbosity
go test -bench=. -benchmem -v=0 ./internal/service

# Run only benchmarks (no regular tests)
go test -bench=. -benchmem -run=^$ ./internal/service

# Run benchmarks with reduced logging
QUIZIZZ_LOG_LEVEL=error go test -bench=. -benchmem ./internal/service
```

## Test Inventory

### Unit Tests

Located in `*_test.go` files alongside the code they test:

1. **Repository Tests** (`internal/repository/mock_user_repository_test.go`):
   - Tests for the mock user repository implementation
   - Tests CRUD operations on users

2. **Service Tests** (`internal/service/user_service_test.go`):
   - Tests for the user service layer
   - Tests validation and business logic for user operations

3. **Handler Tests** (`internal/api/handlers/user/user_handler_test.go`):
   - Tests for API handlers
   - Uses mocked services to isolate handler logic

### Integration Tests

Located in files with `integration` build tag:

1. **API Integration Tests** (`internal/api/handlers/user/user_handler_integration_test.go`):
   - Tests the full API request/response cycle
   - Verifies correct interaction between handlers, services, and repositories
   - Uses **mock external dependencies** (database, Redis) to avoid requiring real services
   - Tests the complete request flow from HTTP input to response

**Important**: Integration tests use mock implementations of external resources (`MockDB`, `MockRedis`) rather than real connections. This approach:
- Makes tests faster and more reliable
- Eliminates external dependencies in test environments
- Focuses on testing the integration between application components rather than external service connectivity

Run with:
```bash
make test-integration
```
or
```bash
go test -tags=integration ./...
```

### Benchmark Tests

Located in `*_benchmark_test.go` files:

1. **Service Benchmarks** (`internal/service/user_service_benchmark_test.go`):
   - `BenchmarkUserService_GetByID`: Measures performance of retrieving a user by ID
   - `BenchmarkUserService_List`: Measures performance of listing all users
   - `BenchmarkUserService_Create`: Measures performance of creating a new user
   - `BenchmarkUserService_Update`: Measures performance of updating an existing user
   - `BenchmarkUserService_Delete`: Measures performance of deleting users (both single and batch)

Run with:
```bash
go test -bench=. -benchmem ./internal/service
```

For clearer benchmark results with less logging:
```bash
go test -bench=. -benchmem -run=^$ ./internal/service
```

## Test Utilities

The project includes several test utilities:

1. **General Test Utilities** (`internal/testutil/testutil.go`):
   - `Setup`: Initializes the test environment
   - `CreateTestServer`: Creates a test HTTP server
   - `MakeTestRequest`: Helper for making HTTP requests in tests
   - `ParseResponse`: Parses JSON responses
   - `AssertStatusCode`: Asserts the HTTP status code
   - `LoadFixture`: Loads test data from fixtures

2. **Integration Test Utilities** (`internal/testutil/integration/integration.go`):
   - `Setup`: Creates a complete test environment for integration testing
   - Initializes all necessary components (router, services, repositories)

3. **Benchmark Utilities** (`internal/service/benchmark_test_helper.go`):
   - `DisableLoggingForBenchmark`: Temporarily disables logging during benchmarks

## Test Data

Test fixtures are stored in the `testdata/` directory:

1. **User Fixtures** (`testdata/users.json`):
   - Sample user data for testing

## Continuous Integration

Tests are automatically run in CI pipelines. A PR cannot be merged if tests are failing.

## Testing Philosophy and Best Practices

### Resource Management in Tests

The project follows a layered testing approach:

1. **Unit Tests**: Use mocks for all external dependencies (repositories, services, HTTP clients)
2. **Integration Tests**: Use mock external resources (DB, Redis) but real application components
3. **End-to-End Tests**: (Future) Use real external dependencies in controlled environments

### Why Mock External Resources in Integration Tests?

Integration tests focus on testing the **integration between application components**, not the connectivity to external services. Using mock resources provides:

- **Reliability**: Tests don't fail due to external service unavailability
- **Speed**: No network calls or database transactions
- **Consistency**: Predictable behavior across different environments
- **Isolation**: Tests can run in parallel without conflicts
- **Simplicity**: No setup/teardown of external services required

### When to Use Real External Dependencies

Consider real external dependencies for:
- **End-to-End Tests**: Testing complete user journeys
- **Performance Testing**: Measuring real-world performance
- **Database Migration Tests**: Validating schema changes
- **Service Compatibility Tests**: Ensuring compatibility with external APIs

## Troubleshooting

If benchmark tests produce too much log output, you can:

1. Set the log level to error:
   ```bash
   QUIZIZZ_LOG_LEVEL=error go test -bench=. -benchmem ./internal/service
   ```

2. Use the `-run=^$` flag to skip regular tests:
   ```bash
   go test -bench=. -benchmem -run=^$ ./internal/service
   ```

3. Use the `DisableLoggingForBenchmark` helper function in benchmark tests:
   ```go
   func BenchmarkMyFunction(b *testing.B) {
       DisableLoggingForBenchmark(b)
       // benchmark code
   }
   ```

### Integration Test Issues

If integration tests fail:

1. **Check build tags**: Ensure files have `//go:build integration` tag
2. **Run with integration tag**: Use `go test -tags=integration ./...`
3. **Verify mock implementations**: Ensure mock resources implement required interfaces