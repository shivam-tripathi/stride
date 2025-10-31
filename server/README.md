# Go Project Template (Gin + Wire)

A template for Go projects using Gin Gonic for HTTP handling and Wire for dependency injection.

## Project Structure

```
.
├── cmd/                # Entry points for different applications
│   └── api/            # Main API service entry point
├── internal/           # Private application code
│   ├── api/            # API handlers
│   ├── config/         # Application configuration
│   ├── domain/         # Domain models
│   ├── repository/     # Data access layer
│   └── service/        # Business logic implementation
├── pkg/                # Public libraries that can be used by external applications
│   └── middleware/     # Reusable middleware
└── wire/               # Dependency injection configuration
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Git

### Setup and Running

1. Clone the repository:
```bash
git clone https://github.com/yourusername/project-template-go-v2.git
cd project-template-go-v2
```

2. Install dependencies:
```bash
make deps
```

3. Generate dependency injection code:
```bash
make wire
```

4. Run the application:
```bash
make run
```

#### Other helpful commands

1. To watch and restart server on changes:
```
brew install watchexec
make watch
```

2. To wire, build the binary and run the server:
```
make dev
```

## Development

### Adding a New Dependency

1. Add the dependency to `go.mod` or use:
```bash
go get example.com/package
```

2. Update the wire provider sets in the `wire` directory as needed.

3. Regenerate the dependency injection code:
```bash
make wire
```

## License

[MIT](LICENSE)

### Makefile targets

- **Build**: `make build` - compiles the API binary to `server` in the repository.
- **Run**: `make run` - runs the API server.
- **All**: `make all` - runs `wire` then `build`.
- **Tests**:
  - **Unit tests**: `make test-unit` (default test target is `make test` which runs unit tests).
  - **Integration tests**: `make test-integration`.
  - **Test coverage**: `make test-coverage` (produces `coverage.out` and `coverage.html`).
  - **Race tests**: `make test-race`.
  - **Short tests**: `make test-short`.
- **Lint**: `make lint`.
- **Clean**: `make clean` - removes the binary and coverage files.
- **Dev**: `make dev` - runs `wire`, `build`, and then `run`.
- **Watch**: `make watch` - runs watchexec on `make dev`
- **Wire**: `make wire` - regenerates dependency injection wiring.
- **Docker**:
  - `make docker-build` - builds the Docker image.
  - `make docker-run` - runs the Docker container.
  - `make docker-stop` - stops and removes the Docker container.
- **Deps**: `make deps` - updates dependencies and tidies `go.mod`.

