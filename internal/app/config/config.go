package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for our application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Security SecurityConfig
	App      AppConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         string
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	AccessSecretKey   string
	RefreshSecretKey  string
	AccessExpiration  time.Duration
	RefreshExpiration time.Duration
	Issuer            string
	Audience          string
}

// SecurityConfig holds security configuration
type SecurityConfig struct {
	RateLimitRequests      int
	RateLimitWindow        time.Duration
	MaxLoginAttempts       int
	AccountLockoutTime     time.Duration
	PasswordMinLength      int
	PasswordRequireUpper   bool
	PasswordRequireLower   bool
	PasswordRequireNumber  bool
	PasswordRequireSpecial bool
	SessionTimeout         time.Duration
	RefreshTokenCleanup    time.Duration
}

// AppConfig holds application configuration
type AppConfig struct {
	Environment string
	Debug       bool
	LogLevel    string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			Host:         getEnv("HOST", "0.0.0.0"),
			ReadTimeout:  getDurationEnv("READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getDurationEnv("WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:  getDurationEnv("IDLE_TIMEOUT", 120*time.Second),
		},
		Database: DatabaseConfig{
			URL:             getEnv("DATABASE_URL", "postgres://go_user:go_password@localhost:5433/go_learning_db?sslmode=disable"),
			MaxOpenConns:    getIntEnv("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS", 25),
			ConnMaxLifetime: getDurationEnv("DB_CONN_MAX_LIFETIME", 5*time.Minute),
			ConnMaxIdleTime: getDurationEnv("DB_CONN_MAX_IDLE_TIME", 5*time.Minute),
		},
		JWT: JWTConfig{
			AccessSecretKey:   getEnv("JWT_ACCESS_SECRET", "your-access-secret-key-change-this-in-production"),
			RefreshSecretKey:  getEnv("JWT_REFRESH_SECRET", "your-refresh-secret-key-change-this-in-production"),
			AccessExpiration:  getDurationEnv("JWT_ACCESS_EXPIRATION", 15*time.Minute),
			RefreshExpiration: getDurationEnv("JWT_REFRESH_EXPIRATION", 7*24*time.Hour),
			Issuer:            getEnv("JWT_ISSUER", "go-backend-api"),
			Audience:          getEnv("JWT_AUDIENCE", "go-backend-api-users"),
		},
		Security: SecurityConfig{
			RateLimitRequests:      getIntEnv("RATE_LIMIT_REQUESTS", 100),
			RateLimitWindow:        getDurationEnv("RATE_LIMIT_WINDOW", time.Minute),
			MaxLoginAttempts:       getIntEnv("MAX_LOGIN_ATTEMPTS", 5),
			AccountLockoutTime:     getDurationEnv("ACCOUNT_LOCKOUT_TIME", 15*time.Minute),
			PasswordMinLength:      getIntEnv("PASSWORD_MIN_LENGTH", 8),
			PasswordRequireUpper:   getBoolEnv("PASSWORD_REQUIRE_UPPER", true),
			PasswordRequireLower:   getBoolEnv("PASSWORD_REQUIRE_LOWER", true),
			PasswordRequireNumber:  getBoolEnv("PASSWORD_REQUIRE_NUMBER", true),
			PasswordRequireSpecial: getBoolEnv("PASSWORD_REQUIRE_SPECIAL", true),
			SessionTimeout:         getDurationEnv("SESSION_TIMEOUT", 24*time.Hour),
			RefreshTokenCleanup:    getDurationEnv("REFRESH_TOKEN_CLEANUP", time.Hour),
		},
		App: AppConfig{
			Environment: getEnv("ENVIRONMENT", "development"),
			Debug:       getBoolEnv("DEBUG", true),
			LogLevel:    getEnv("LOG_LEVEL", "info"),
		},
	}
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getIntEnv gets an integer environment variable with a fallback value
func getIntEnv(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}

// getBoolEnv gets a boolean environment variable with a fallback value
func getBoolEnv(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return fallback
}

// getDurationEnv gets a duration environment variable with a fallback value
func getDurationEnv(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return fallback
}

// IsProduction returns true if the environment is production
func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

// IsDevelopment returns true if the environment is development
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}
