package main

import (
	"go-backend-api/api"
	"go-backend-api/internal/config"
	"go-backend-api/internal/database"
	"go-backend-api/internal/handlers"
	"go-backend-api/internal/logger"
	"go-backend-api/internal/middleware"
	"go-backend-api/internal/pkg/auth"
	"go-backend-api/internal/repositories"
	"go-backend-api/internal/services"

	"github.com/gin-gonic/gin"
)

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
	refreshTokenRepo := repositories.NewRefreshTokenRepository(database.GetDB())

	// Initialize services
	userService := services.NewUserService(userRepo, refreshTokenRepo, jwtManager)
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

	// OpenAPI documentation and specification endpoints (at root level)
	router.GET("/docs", api.ServeOpenAPIDocs)
	router.GET("/api-docs", api.ServeOpenAPIDocs) // Alias for /docs
	router.GET("/openapi.yaml", api.ServeOpenAPISpec)
	router.GET("/openapi.json", api.ServeOpenAPISpec)

	// API routes with /api/v1 prefix
	api := router.Group("/api/v1")
	{
		// Health check endpoint (under /api/v1 for consistency)
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"message": "Go Backend API is running",
			})
		})

		// Public routes (no authentication required)
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/refresh", authHandler.Refresh)
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
				users.POST("/logout", userHandler.Logout)
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
