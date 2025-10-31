.PHONY: all build run test test-unit test-integration test-coverage test-race clean wire docker-build docker-run docker-stop lint

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOLINT=golangci-lint
WIRE=wire
GOTOOLS=$(GOPATH)/bin
DOCKER=docker

# Binary name
BINARY_NAME=server

# Test parameters
COVERPROFILE=coverage.out
COVERHTML=coverage.html

# Dependency injection
wire:
	cd wire && $(WIRE)

# Build
build:
	$(GOBUILD) -o $(BINARY_NAME) ./cmd/server

run:
	$(GORUN) ./cmd/server

dev: wire build run

watch:
	@clear
	watchexec -r -e go -- 'make dev'

all: wire build

dev: wire build run

watch:
	@clear
	watchexec -r -e go -- 'make dev'

# Testing
test: test-unit

test-unit:
	$(GOTEST) -v ./...

test-integration:
	$(GOTEST) -v -tags=integration ./...

test-coverage:
	$(GOTEST) -v -coverprofile=$(COVERPROFILE) ./...
	$(GOCMD) tool cover -html=$(COVERPROFILE) -o $(COVERHTML)
	$(GOCMD) tool cover -func=$(COVERPROFILE)

test-race:
	$(GOTEST) -v -race ./...

test-short:
	$(GOTEST) -v -short ./...

# Linting
lint:
	$(GOLINT) run

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(COVERPROFILE)
	rm -f $(COVERHTML)

# Docker commands
docker-build:
	docker build -t $(BINARY_NAME) .

docker-run:
	docker run -p 8080:8080 --name $(BINARY_NAME) -d $(BINARY_NAME)

docker-stop:
	$(DOCKER) stop $(BINARY_NAME)
	$(DOCKER) rm $(BINARY_NAME)

# Dependencies
deps:
	$(GOGET) -u ./...
	go mod tidy
