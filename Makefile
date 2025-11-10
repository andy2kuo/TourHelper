.PHONY: help build run test clean docker-build docker-up docker-down lint fmt

# Default target
help:
	@echo "Available targets:"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-up    - Start Docker containers"
	@echo "  docker-down  - Stop Docker containers"
	@echo "  lint         - Run linter"
	@echo "  fmt          - Format code"

# Build the application
build:
	@echo "Building TourHelper..."
	go build -o bin/api ./cmd/api

# Run the application
run:
	@echo "Running TourHelper..."
	go run ./cmd/api/main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v -cover ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	go clean

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t tourhelper:latest .

# Start Docker containers
docker-up:
	@echo "Starting Docker containers..."
	docker-compose up -d

# Stop Docker containers
docker-down:
	@echo "Stopping Docker containers..."
	docker-compose down

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run ./...

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w .

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Run database migrations (placeholder for future)
migrate-up:
	@echo "Running migrations..."
	# Add migration command here

# Run database rollback (placeholder for future)
migrate-down:
	@echo "Rolling back migrations..."
	# Add rollback command here
