#!/bin/bash

# Database Migration Script
echo "🗄️  Running database migrations..."

# Check if PostgreSQL container is running
if ! docker ps | grep -q go-learning-postgres; then
    echo "❌ PostgreSQL container is not running. Please start it first:"
    echo "   docker-compose up -d postgres"
    exit 1
fi

# Run migrations
docker exec -i go-learning-postgres psql -U go_user -d go_learning_db < internal/database/migrations_v2.sql

if [ $? -eq 0 ]; then
    echo "✅ Database migrations completed successfully!"
else
    echo "❌ Failed to run database migrations"
    exit 1
fi
