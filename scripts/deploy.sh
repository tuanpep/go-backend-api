#!/bin/bash

# Production Deployment Script for Go Backend API
# Usage: ./scripts/deploy.sh [--pull] [--build] [--logs]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
COMPOSE_FILE="docker-compose.prod.yml"
ENV_FILE=".env.production"
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

cd "$PROJECT_DIR"

# Check if .env.production exists
if [ ! -f "$ENV_FILE" ]; then
    echo -e "${YELLOW}Warning: $ENV_FILE not found!${NC}"
    echo -e "${BLUE}Creating from template...${NC}"
    if [ -f "env.production.example" ]; then
        cp env.production.example "$ENV_FILE"
        echo -e "${YELLOW}Please edit $ENV_FILE with your production values before deploying!${NC}"
        echo -e "${RED}Exiting...${NC}"
        exit 1
    else
        echo -e "${RED}Error: env.production.example not found!${NC}"
        exit 1
    fi
fi

# Parse arguments
PULL=false
BUILD=false
LOGS=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --pull)
            PULL=true
            shift
            ;;
        --build)
            BUILD=true
            shift
            ;;
        --logs)
            LOGS=true
            shift
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            echo "Usage: $0 [--pull] [--build] [--logs]"
            exit 1
            ;;
    esac
done

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Go Backend API - Production Deployment${NC}"
echo -e "${BLUE}========================================${NC}"

# Pull latest code if requested
if [ "$PULL" = true ]; then
    echo -e "${BLUE}Pulling latest code from git...${NC}"
    git pull || echo -e "${YELLOW}Warning: git pull failed or not in a git repository${NC}"
fi

# Build flag for docker-compose
BUILD_FLAG=""
if [ "$BUILD" = true ]; then
    BUILD_FLAG="--build"
    echo -e "${BLUE}Building Docker images...${NC}"
fi

# Stop existing containers
echo -e "${BLUE}Stopping existing containers...${NC}"
docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" down

# Start services
echo -e "${BLUE}Starting services...${NC}"
docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" up -d $BUILD_FLAG

# Wait for services to be healthy
echo -e "${BLUE}Waiting for services to be healthy...${NC}"
sleep 5

# Check service status
echo -e "${BLUE}Service status:${NC}"
docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" ps

# Run database migrations if needed
echo -e "${BLUE}Checking database migrations...${NC}"
if docker exec go-api-postgres psql -U go_user -d go_learning_db -c "\dt" > /dev/null 2>&1; then
    echo -e "${GREEN}Database is accessible${NC}"
    # You can add migration commands here if needed
    # docker exec go-api-postgres psql -U go_user -d go_learning_db -f /docker-entrypoint-initdb.d/init.sql
else
    echo -e "${YELLOW}Warning: Could not connect to database${NC}"
fi

# Health check
echo -e "${BLUE}Performing health check...${NC}"
sleep 3
if curl -f http://localhost:${API_PORT:-8080}/api/v1/health > /dev/null 2>&1; then
    echo -e "${GREEN}✓ API is healthy${NC}"
else
    echo -e "${YELLOW}Warning: Health check failed, but services are starting...${NC}"
fi

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Deployment completed!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${BLUE}Useful commands:${NC}"
echo "  View logs:     docker-compose -f $COMPOSE_FILE --env-file $ENV_FILE logs -f"
echo "  Stop services: docker-compose -f $COMPOSE_FILE --env-file $ENV_FILE down"
echo "  Restart:       docker-compose -f $COMPOSE_FILE --env-file $ENV_FILE restart"
echo "  Status:        docker-compose -f $COMPOSE_FILE --env-file $ENV_FILE ps"
echo ""

# Show logs if requested
if [ "$LOGS" = true ]; then
    echo -e "${BLUE}Showing logs (Ctrl+C to exit)...${NC}"
    docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" logs -f
fi

