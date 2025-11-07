package main

import (
	"go-backend-api/docs"
	"go-backend-api/internal/app/config"
	"go-backend-api/internal/app/logger"
	"go-backend-api/internal/application/handlers"
	"go-backend-api/internal/application/repositories"
	"go-backend-api/internal/application/services"
	"go-backend-api/internal/database"
	"go-backend-api/internal/middleware"
	"go-backend-api/internal/pkg/auth"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Go Backend API
// @version         1.0
// @description     A comprehensive REST API built with Go for learning backend development
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize logger
	logger := logger.NewLogger(cfg.App.LogLevel)

	// Connect to database
	if err := database.Connect(cfg.Database.URL); err != nil {
		logger.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// Note: Run migrations manually using the SQL file
	// psql -h localhost -p 5433 -U go_user -d go_learning_db -f internal/database/migrations_v2.sql

	// Set Gin mode
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(
		cfg.JWT.AccessSecretKey,
		cfg.JWT.RefreshSecretKey,
		cfg.JWT.Issuer,
		cfg.JWT.Audience,
		cfg.JWT.AccessExpiration,
		cfg.JWT.RefreshExpiration,
	)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(database.GetDB())
	postRepo := repositories.NewPostRepository(database.GetDB())

	// Initialize services
	userService := services.NewUserService(userRepo, jwtManager)
	postService := services.NewPostService(postRepo, userRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userService, jwtManager)
	userHandler := handlers.NewUserHandler(userService)
	postHandler := handlers.NewPostHandler(postService)

	// Create Gin router
	router := gin.New()

	// Add middleware
	router.Use(logger.GinLogger())
	router.Use(logger.GinRecovery())
	router.Use(middleware.CORS())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Go Backend API is running",
		})
	})

	// Swagger documentation
	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
		}

		// Protected routes (authentication required)
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(jwtManager))
		{
			// Current user endpoint
			protected.GET("/me", userHandler.GetMe)

			// User routes
			users := protected.Group("/users")
			{
				users.GET("/profile", userHandler.GetProfile)
				users.PUT("/profile", userHandler.UpdateProfile)
				users.DELETE("/profile", userHandler.DeleteProfile)
				users.PUT("/:id/activate", userHandler.ActivateUser)
				users.PUT("/:id/deactivate", userHandler.DeactivateUser)
			}

			// Post routes
			posts := protected.Group("/posts")
			{
				posts.POST("", postHandler.Create)
				posts.GET("", postHandler.GetAll)
				posts.GET("/:id", postHandler.GetByID)
				posts.PUT("/:id", postHandler.Update)
				posts.DELETE("/:id", postHandler.Delete)
			}
		}
	}

	// Start server
	logger.Infof("Server starting on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		logger.Fatal("Failed to start server:", err)
	}
}
