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

.PHONY: all build clean test deps fmt run dev setup stop logs help

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

# Run the application
run: build
	@echo "$(BLUE)Starting $(BINARY_NAME)...$(NC)"
	./$(BINARY_NAME)

# Development mode (build and run)
dev: build
	@echo "$(BLUE)Starting in development mode...$(NC)"
	./$(BINARY_NAME)

# Setup the project
setup:
	@echo "$(BLUE)Setting up the project...$(NC)"
	chmod +x setup.sh
	./setup.sh

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
	chmod +x migrate.sh
	./migrate.sh

# Test the API
test-api:
	@echo "$(BLUE)Testing API endpoints...$(NC)"
	chmod +x test_api.sh
	./test_api.sh

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
	@echo "  run       - Build and run the application"
	@echo "  dev       - Development mode (build and run)"
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
	@echo "  make run      # Build and run"
	@echo "  make test-api # Test the API"