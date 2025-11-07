package entities

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user entity
type User struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	Username  string     `json:"username" db:"username"`
	Email     string     `json:"email" db:"email"`
	Password  string     `json:"-" db:"password"` // Hidden from JSON output
	IsActive  bool       `json:"is_active" db:"is_active"`
	LastLogin *time.Time `json:"last_login,omitempty" db:"last_login"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(user *User) error
	GetByID(id uuid.UUID) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByUsername(username string) (*User, error)
	Update(user *User) error
	Delete(id uuid.UUID) error
	ExistsByEmail(email string) (bool, error)
	ExistsByUsername(username string) (bool, error)
	UpdateLastLogin(id uuid.UUID) error
	Activate(id uuid.UUID) error
	Deactivate(id uuid.UUID) error
}

// UserService defines the interface for user business logic
type UserService interface {
	CreateUser(req *CreateUserRequest) (*User, error)
	GetUserByID(id uuid.UUID) (*User, error)
	GetUserByEmail(email string) (*User, error)
	UpdateUser(id uuid.UUID, req *UpdateUserRequest) (*User, error)
	DeleteUser(id uuid.UUID) error
	ValidateUser(user *User) error
	AuthenticateUser(req *LoginRequest) (*LoginResponse, error)
	RefreshToken(req *RefreshTokenRequest) (*LoginResponse, error)
	Logout(userID uuid.UUID, tokenID string) error
	ActivateUser(id uuid.UUID) error
	DeactivateUser(id uuid.UUID) error
}

// CreateUserRequest represents the request to create a user
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,username"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	Username string `json:"username,omitempty" validate:"omitempty,username"`
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
}

// LoginRequest represents the request to login a user
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RefreshTokenRequest represents the request to refresh a token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// LoginResponse represents the response after successful login
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	User         User   `json:"user"`
}

// TokenClaims represents JWT token claims
type TokenClaims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	TokenID  string    `json:"token_id"`
	Type     string    `json:"type"` // "access" or "refresh"
}
