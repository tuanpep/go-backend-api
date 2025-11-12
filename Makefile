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

.PHONY: all build clean test deps fmt run run-once dev watch setup stop logs help openapi

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

# Generate OpenAPI types and validate setup
openapi:
	@echo "$(BLUE)Setting up OpenAPI...$(NC)"
	@echo "$(BLUE)1. Checking oapi-codegen installation...$(NC)"
	@if ! command -v oapi-codegen > /dev/null 2>&1 && [ ! -f $$(go env GOPATH)/bin/oapi-codegen ]; then \
		echo "$(YELLOW)   Installing oapi-codegen...$(NC)"; \
		go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest; \
	fi
	@echo "$(BLUE)2. Validating OpenAPI specification...$(NC)"
	@if [ ! -f api/openapi.yaml ]; then \
		echo "$(RED)   Error: api/openapi.yaml not found!$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)   ✓ OpenAPI spec found$(NC)"
	@echo "$(BLUE)3. Generating Go types from OpenAPI spec...$(NC)"
	@if command -v oapi-codegen > /dev/null 2>&1; then \
		oapi-codegen -generate types -package api api/openapi.yaml > api/types.go; \
	elif [ -f $$(go env GOPATH)/bin/oapi-codegen ]; then \
		$$(go env GOPATH)/bin/oapi-codegen -generate types -package api api/openapi.yaml > api/types.go; \
	else \
		echo "$(RED)   Error: oapi-codegen not found after installation$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)   ✓ Types generated in api/types.go$(NC)"
	@echo "$(BLUE)4. Verifying generated files...$(NC)"
	@if [ ! -f api/types.go ]; then \
		echo "$(RED)   Error: api/types.go was not generated!$(NC)"; \
		exit 1; \
	fi
	@if [ ! -f api/serve.go ]; then \
		echo "$(YELLOW)   Warning: api/serve.go not found (should be created manually)$(NC)"; \
	fi
	@echo "$(GREEN)   ✓ All files verified$(NC)"
	@echo "$(GREEN)OpenAPI setup complete!$(NC)"
	@echo "$(YELLOW)Note:$(NC)"
	@echo "  - OpenAPI spec: api/openapi.yaml (edit this file to update API)"
	@echo "  - Generated types: api/types.go (auto-generated, do not edit)"
	@echo "  - Server handler: api/serve.go (serves the OpenAPI spec)"
	@echo "  - View docs: http://localhost:8080/docs (HTML documentation)"
	@echo "  - Raw spec: http://localhost:8080/openapi.yaml (download YAML)"

# Show help
help:
	@echo "$(BLUE)Go Backend API - Available Commands:$(NC)"
	@echo ""
	@echo "$(GREEN)Build Commands:$(NC)"
	@echo "  build     - Build the application"
	@echo "  clean     - Clean build artifacts"
	@echo "  deps      - Download and update dependencies"
	@echo "  fmt       - Format code"
	@echo "  openapi   - Generate OpenAPI types and validate setup"
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
	@echo "  make openapi  # Generate OpenAPI types and validate setup"
	@echo "  make test-api # Test the API"