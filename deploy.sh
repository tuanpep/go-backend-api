#!/bin/bash
# Manual deployment script
# Usage: ./deploy.sh

set -e

cd /opt/go-backend-api

echo "=== Starting deployment ==="

# Pull latest changes
if [ -d .git ]; then
    echo "Pulling latest changes from git..."
    git fetch origin
    git reset --hard origin/main || git reset --hard origin/master
else
    echo "Warning: Not a git repository. Skipping git pull."
fi

# Pull latest images and restart containers
echo "Pulling latest Docker images..."
docker compose -f docker-compose.prod.yml --env-file .env.production pull

echo "Restarting containers..."
docker compose -f docker-compose.prod.yml --env-file .env.production up -d

# Clean up old images
echo "Cleaning up old Docker images..."
docker image prune -f

echo "=== Deployment completed ==="
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
