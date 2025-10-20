# Go Backend API Learning Project

A comprehensive REST API built with Go for learning backend development. This project demonstrates common patterns, best practices, and features you'll encounter in real-world Go backend applications.

## 🚀 Features

- **User Authentication & Authorization** - JWT-based auth system
- **RESTful API Design** - Clean, consistent API endpoints
- **Database Integration** - PostgreSQL with migrations
- **Middleware** - CORS, logging, authentication
- **Error Handling** - Structured error responses
- **Data Validation** - Request validation and sanitization
- **Pagination** - Efficient data pagination
- **Security** - Password hashing, JWT tokens

## 📁 Project Structure

```
go-backend-api/
├── cmd/
│   └── main.go                    # Application entry point
├── internal/
│   ├── app/
│   │   └── config/
│   │       └── config.go          # Configuration management
│   ├── application/
│   │   ├── handlers/              # HTTP handlers
│   │   │   ├── auth_handler.go
│   │   │   ├── post_handler.go
│   │   │   └── user_handler.go
│   │   ├── repositories/          # Data access layer
│   │   │   ├── post_repository.go
│   │   │   └── user_repository.go
│   │   └── services/              # Business logic layer
│   │       ├── post_service.go
│   │       └── user_service.go
│   ├── database/
│   │   ├── database.go            # Database connection
│   │   └── migrations_v2.sql      # Database migrations
│   ├── domain/
│   │   └── entities/              # Domain entities
│   │       ├── post.go
│   │       └── user.go
│   ├── middleware/                # HTTP middleware
│   │   ├── auth.go
│   │   ├── cors.go
│   │   └── logger.go
│   └── pkg/                       # Shared packages
│       ├── auth/                  # JWT authentication
│       ├── errors/                # Error handling
│       ├── response/              # HTTP responses
│       ├── security/              # Security utilities
│       └── validation/            # Input validation
├── go.mod                         # Go module dependencies
└── README.md                     # This file
```

## 🛠️ Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- Git

> **Note**: You can also use a local PostgreSQL installation, but Docker is recommended for easy setup.

## 🚀 Getting Started

### Quick Start (Recommended)

```bash
# Navigate to your project directory
cd /home/tuanbt/Learning/go-backend-api

# Complete setup (database + dependencies + build)
make setup

# Run the application
make run

# Test the API
make test-api
```

The setup script will:
- Start PostgreSQL with Docker
- Download Go dependencies
- Build the application
- Create environment configuration
- Set up the database with sample data

### Manual Setup

#### 1. Database Setup with Docker

```bash
# Start PostgreSQL with Docker Compose
docker-compose up -d postgres

# Check if PostgreSQL is running
docker-compose ps postgres
```

#### 2. Application Setup

```bash
# Download dependencies
go mod tidy

# Build the application
go build -o bin/main cmd/main.go

# Or run directly
go run cmd/main.go
```

#### 3. Environment Configuration

The application will use default configuration, but you can create a `.env` file:

```bash
# Server Configuration
PORT=8080
ENVIRONMENT=development

# Database Configuration (Docker PostgreSQL)
DATABASE_URL=postgres://go_user:go_password@localhost:5432/go_learning_db?sslmode=disable

# JWT Configuration
JWT_SECRET=your-secret-key-change-this-in-production
```

### Alternative: Full Docker Setup

To run everything in Docker containers:

```bash
# Build and start all services (API + PostgreSQL + pgAdmin)
docker-compose -f docker-compose.full.yml up --build

# Or run in background
docker-compose -f docker-compose.full.yml up -d --build
```

## 🐳 Docker Services

- **PostgreSQL**: `localhost:5432`
- **API**: `localhost:8080`
- **pgAdmin**: `localhost:5050` (admin@example.com / admin123)

## 📚 API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login user

### Users (Protected)
- `GET /api/v1/users/profile` - Get current user profile
- `PUT /api/v1/users/profile` - Update current user profile
- `DELETE /api/v1/users/profile` - Delete current user account

### Posts (Protected)
- `POST /api/v1/posts` - Create a new post
- `GET /api/v1/posts` - Get all posts (with pagination)
- `GET /api/v1/posts/:id` - Get a specific post
- `PUT /api/v1/posts/:id` - Update a post (author only)
- `DELETE /api/v1/posts/:id` - Delete a post (author only)

### Health Check
- `GET /health` - Health check endpoint

## 🧪 Testing the API

### 1. Register a User
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 2. Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 3. Create a Post (with JWT token)
```bash
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "My First Post",
    "content": "This is the content of my first blog post!"
  }'
```

### 4. Get All Posts
```bash
curl -X GET "http://localhost:8080/api/v1/posts?page=1&per_page=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 🎓 Learning Concepts Demonstrated

### 1. **Go Project Structure**
- Standard Go project layout
- Separation of concerns
- Internal vs external packages

### 2. **HTTP Server with Gin**
- RESTful API design
- Route grouping and middleware
- Request/response handling

### 3. **Database Operations**
- PostgreSQL integration
- SQL queries and prepared statements
- Database migrations
- Connection management

### 4. **Authentication & Security**
- JWT token generation and validation
- Password hashing with bcrypt
- Middleware for route protection

### 5. **Error Handling**
- Structured error responses
- HTTP status codes
- Graceful error handling

### 6. **Data Validation**
- Request validation with Gin
- Input sanitization
- Custom validation rules

### 7. **Configuration Management**
- Environment variables
- Configuration structs
- Environment-specific settings

## 🔧 Development

### Available Commands
```bash
# Build and run
make build          # Build the application
make run            # Build and run
make dev            # Development mode

# Database
make db-up          # Start PostgreSQL
make db-down        # Stop PostgreSQL
make db-logs        # View database logs
make migrate        # Run migrations

# Testing and quality
make test           # Run tests
make test-api       # Test API endpoints
make fmt            # Format code

# Setup
make setup          # Complete project setup
make deps           # Download dependencies
make clean          # Clean build artifacts

# Help
make help           # Show all commands
```

### Adding New Features
1. Define entities in `internal/domain/entities/`
2. Create repositories in `internal/application/repositories/`
3. Create services in `internal/application/services/`
4. Create handlers in `internal/application/handlers/`
5. Add routes in `cmd/main.go`
6. Update database schema in `internal/database/migrations_v2.sql`

### Database Migrations
- Add new migration queries to `internal/database/migrations_v2.sql`
- Run migrations using Docker: `docker exec -i go-learning-postgres psql -U go_user -d go_learning_db < internal/database/migrations_v2.sql`
- Or install psql locally: `sudo apt install postgresql-client-common`
- Always backup data before schema changes

### Testing
- Use tools like Postman or curl for API testing
- Test both success and error scenarios
- Verify authentication and authorization

## 🚀 Next Steps for Learning

1. **Add Unit Tests** - Learn Go testing with `testing` package
2. **Add Integration Tests** - Test database operations
3. **Add Docker Support** - Containerize your application
4. **Add Logging** - Implement structured logging
5. **Add Caching** - Implement Redis caching
6. **Add Rate Limiting** - Protect against abuse
7. **Add API Documentation** - Use Swagger/OpenAPI
8. **Add Monitoring** - Add health checks and metrics

## 📖 Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Gin Framework](https://gin-gonic.com/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [JWT.io](https://jwt.io/) - JWT token debugging
- [REST API Best Practices](https://restfulapi.net/)

## 🤝 Contributing

This is a learning project! Feel free to:
- Add new features
- Improve existing code
- Fix bugs
- Add tests
- Improve documentation

Happy coding! 🎉
