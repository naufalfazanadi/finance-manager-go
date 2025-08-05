# Go Finance Manager - Development Commands

# Variables
BINARY_NAME=finance-manager
BUILD_DIR=./tmp
MAIN_PATH=./cmd/server

# Development commands
.PHONY: dev build run clean test deps air swagger

# Run with air (live reload)
dev:
	air -c .air.toml

# Generate Swagger documentation
swagger:
	swag init -g cmd/server/main.go -o docs --parseDependency

# Build the application
build: swagger
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Run the built binary
run: build
	$(BUILD_DIR)/$(BINARY_NAME)

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
	go clean

# Run tests
test:
	go test -v ./...

# Download dependencies
deps:
	go mod download
	go mod verify

# Install air for development
install-air:
	go install github.com/air-verse/air@latest

# Run without live reload
start:
	go run $(MAIN_PATH)

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Tidy dependencies
tidy:
	go mod tidy

# Show help
help:
	@echo "Available commands:"
	@echo "  dev          - Run with air (live reload)"
	@echo "  build        - Build the application"
	@echo "  run          - Build and run the application"
	@echo "  start        - Run without live reload"
	@echo "  clean        - Clean build artifacts"
	@echo "  test         - Run tests"
	@echo "  deps         - Download dependencies"
	@echo "  install-air  - Install air for development"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code"
	@echo "  tidy         - Tidy dependencies"
	@echo "  help         - Show this help"
