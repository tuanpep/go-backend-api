#!/bin/bash

# Go Backend API Setup Script
echo "🚀 Setting up Go Backend API with Docker PostgreSQL..."

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    echo "   Visit: https://docs.docker.com/get-docker/"
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    echo "   Visit: https://docs.docker.com/compose/install/"
    exit 1
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    echo "   Visit: https://golang.org/doc/install"
    exit 1
fi

echo "✅ All prerequisites are installed!"

# Start PostgreSQL with Docker Compose
echo "🐘 Starting PostgreSQL with Docker Compose..."
docker-compose up -d postgres

# Wait for PostgreSQL to be ready
echo "⏳ Waiting for PostgreSQL to be ready..."
sleep 10

# Check if PostgreSQL is running
if ! docker-compose ps postgres | grep -q "Up"; then
    echo "❌ Failed to start PostgreSQL. Check the logs:"
    docker-compose logs postgres
    exit 1
fi

echo "✅ PostgreSQL is running!"

# Wait a moment for PostgreSQL to fully start
echo "⏳ Waiting for PostgreSQL to be ready..."
sleep 5

# Run database migrations
echo "🗄️  Running database migrations..."
docker exec -i go-learning-postgres psql -U go_user -d go_learning_db < internal/database/migrations_v2.sql
if [ $? -eq 0 ]; then
    echo "✅ Database migrations completed successfully!"
else
    echo "❌ Failed to run database migrations"
    exit 1
fi

# Download Go dependencies
echo "📦 Downloading Go dependencies..."
go mod tidy

# Build the application
echo "🔨 Building the application..."
go build -o bin/main cmd/main.go

if [ $? -eq 0 ]; then
    echo "✅ Application built successfully!"
else
    echo "❌ Failed to build application"
    exit 1
fi

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo "📝 Creating .env file..."
    cat > .env << EOF
# Server Configuration
PORT=8080
ENVIRONMENT=development

# Database Configuration (Docker PostgreSQL)
DATABASE_URL=postgres://go_user:go_password@localhost:5433/go_learning_db?sslmode=disable

# JWT Configuration
JWT_ACCESS_SECRET=your-super-secure-access-secret-key-here
JWT_REFRESH_SECRET=your-super-secure-refresh-secret-key-here
JWT_ACCESS_EXPIRATION=15m
JWT_REFRESH_EXPIRATION=168h
JWT_ISSUER=go-backend-api
JWT_AUDIENCE=go-backend-api-users

# Security Configuration
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=1m
MAX_LOGIN_ATTEMPTS=5
ACCOUNT_LOCKOUT_TIME=15m
PASSWORD_MIN_LENGTH=8
PASSWORD_REQUIRE_UPPER=true
PASSWORD_REQUIRE_LOWER=true
PASSWORD_REQUIRE_NUMBER=true
PASSWORD_REQUIRE_SPECIAL=true
EOF
    echo "✅ .env file created!"
else
    echo "ℹ️  .env file already exists, skipping creation"
fi

echo ""
echo "🎉 Setup completed successfully!"
echo ""
echo "📋 Next steps:"
echo "1. Start the application: ./bin/main"
echo "2. Or run directly: go run cmd/main.go"
echo "3. Test the API: ./test_api.sh"
echo "4. Access pgAdmin: http://localhost:5050 (admin@example.com / admin123)"
echo ""
echo "🔗 API Endpoints:"
echo "   Health: http://localhost:8080/health"
echo "   API: http://localhost:8080/api/v1/"
echo ""
echo "📚 Database:"
echo "   Host: localhost"
echo "   Port: 5433"
echo "   Database: go_learning_db"
echo "   Username: go_user"
echo "   Password: go_password"
