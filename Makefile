# Go Backend API Makefile

# Variables
BINARY_NAME=go-backend-api
MAIN_PATH=./cmd/main.go
DOCKER_COMPOSE=docker-compose

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

.PHONY: all build clean test deps fmt run run-once dev watch setup stop logs help

# Default target
all: clean deps fmt test build

# Build the application
build:
	@echo "$(BLUE)Building $(BINARY_NAME)...$(NC)"
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)Build completed!$(NC)"

# Clean build artifacts
clean:
	@echo "$(YELLOW)Cleaning...$(NC)"
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	@echo "$(GREEN)Clean completed!$(NC)"

# Run tests
test:
	@echo "$(BLUE)Running tests...$(NC)"
	$(GOTEST) -v ./...
	@echo "$(GREEN)Tests completed!$(NC)"

# Download dependencies
deps:
	@echo "$(BLUE)Downloading dependencies...$(NC)"
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "$(GREEN)Dependencies updated!$(NC)"

# Format code
fmt:
	@echo "$(BLUE)Formatting code...$(NC)"
	$(GOFMT) ./...
	@echo "$(GREEN)Code formatted!$(NC)"

# Run the application with hot-reload (auto-restart on file changes)
run:
	@echo "$(BLUE)Starting $(BINARY_NAME) with hot-reload...$(NC)"
	@if command -v air > /dev/null 2>&1; then \
		air; \
	elif [ -f $$(go env GOPATH)/bin/air ]; then \
		$$(go env GOPATH)/bin/air; \
	else \
		echo "$(YELLOW)Air not found. Installing...$(NC)"; \
		go install github.com/air-verse/air@latest; \
		$$(go env GOPATH)/bin/air; \
	fi

# Run the application once (no hot-reload)
run-once: build
	@echo "$(BLUE)Starting $(BINARY_NAME)...$(NC)"
	./$(BINARY_NAME)

# Development mode with hot-reload (alias for run)
dev: run

# Watch mode (alias for run)
watch: run

# Setup the project
setup:
	@echo "$(BLUE)Setting up the project...$(NC)"
	chmod +x scripts/setup.sh
	./scripts/setup.sh

# Start database
db-up:
	@echo "$(BLUE)Starting PostgreSQL database...$(NC)"
	$(DOCKER_COMPOSE) up -d postgres
	@echo "$(GREEN)Database started!$(NC)"

# Stop database
db-down:
	@echo "$(YELLOW)Stopping PostgreSQL database...$(NC)"
	$(DOCKER_COMPOSE) down
	@echo "$(GREEN)Database stopped!$(NC)"

# View database logs
db-logs:
	@echo "$(BLUE)Showing database logs...$(NC)"
	$(DOCKER_COMPOSE) logs -f postgres

# Run database migrations
migrate:
	@echo "$(BLUE)Running database migrations...$(NC)"
	chmod +x scripts/migrate.sh
	./scripts/migrate.sh

# Test the API
test-api:
	@echo "$(BLUE)Testing API endpoints...$(NC)"
	chmod +x scripts/test_api.sh
	./scripts/test_api.sh

# Show help
help:
	@echo "$(BLUE)Go Backend API - Available Commands:$(NC)"
	@echo ""
	@echo "$(GREEN)Build Commands:$(NC)"
	@echo "  build     - Build the application"
	@echo "  clean     - Clean build artifacts"
	@echo "  deps      - Download and update dependencies"
	@echo "  fmt       - Format code"
	@echo ""
	@echo "$(GREEN)Run Commands:$(NC)"
	@echo "  run       - Run with hot-reload (auto-restart on file changes)"
	@echo "  run-once  - Build and run once (no hot-reload)"
	@echo "  dev       - Alias for run (hot-reload mode)"
	@echo "  watch     - Alias for run (hot-reload mode)"
	@echo ""
	@echo "$(GREEN)Database Commands:$(NC)"
	@echo "  db-up     - Start PostgreSQL database"
	@echo "  db-down   - Stop PostgreSQL database"
	@echo "  db-logs   - View database logs"
	@echo "  migrate   - Run database migrations"
	@echo ""
	@echo "$(GREEN)Setup Commands:$(NC)"
	@echo "  setup     - Complete project setup"
	@echo "  test-api  - Test API endpoints"
	@echo ""
	@echo "$(GREEN)Other Commands:$(NC)"
	@echo "  test      - Run tests"
	@echo "  help      - Show this help message"
	@echo ""
	@echo "$(YELLOW)Quick Start:$(NC)"
	@echo "  make setup    # Complete setup"
	@echo "  make run      # Run with hot-reload (recommended for development)"
	@echo "  make run-once # Build and run once"
	@echo "  make test-api # Test the API"