package services

import (
	"time"

	"go-backend-api/internal/models"
	"go-backend-api/internal/pkg/auth"
	"go-backend-api/internal/pkg/errors"
	"go-backend-api/internal/pkg/validation"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// userService implements UserService interface
type userService struct {
	userRepo  models.UserRepository
	jwtMgr    *auth.JWTManager
	validator *validation.Validator
}

// NewUserService creates a new user service
func NewUserService(userRepo models.UserRepository, jwtMgr *auth.JWTManager) models.UserService {
	return &userService{
		userRepo:  userRepo,
		jwtMgr:    jwtMgr,
		validator: validation.NewValidator(),
	}
}

// CreateUser creates a new user
func (s *userService) CreateUser(req *models.CreateUserRequest) (*models.User, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, errors.WrapErrorWithCode(err, 400, "Validation failed")
	}

	// Check if user already exists
	exists, err := s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, errors.WrapError(err, "Failed to check user existence")
	}
	if exists {
		return nil, errors.ErrUserExists
	}

	exists, err = s.userRepo.ExistsByUsername(req.Username)
	if err != nil {
		return nil, errors.WrapError(err, "Failed to check username existence")
	}
	if exists {
		return nil, errors.NewAppErrorWithDetails(409, "Username already taken", "Username must be unique", nil)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.WrapError(err, "Failed to hash password")
	}

	// Create user (active by default)
	user := &models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.WrapError(err, "Failed to create user")
	}

	// Clear password from response
	user.Password = ""

	return user, nil
}

// GetUserByID gets a user by ID
func (s *userService) GetUserByID(id uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, errors.WrapError(err, "Failed to get user")
	}
	if user == nil {
		return nil, errors.ErrUserNotFound
	}

	// Clear password from response
	user.Password = ""

	return user, nil
}

// GetUserByEmail gets a user by email
func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, errors.WrapError(err, "Failed to get user")
	}
	if user == nil {
		return nil, errors.ErrUserNotFound
	}

	return user, nil
}

// UpdateUser updates a user
func (s *userService) UpdateUser(id uuid.UUID, req *models.UpdateUserRequest) (*models.User, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, errors.WrapErrorWithCode(err, 400, "Validation failed")
	}

	// Get existing user
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, errors.WrapError(err, "Failed to get user")
	}
	if user == nil {
		return nil, errors.ErrUserNotFound
	}

	// Update fields if provided
	if req.Username != "" {
		// Check if username is already taken by another user
		exists, err := s.userRepo.ExistsByUsername(req.Username)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to check username existence")
		}
		if exists && user.Username != req.Username {
			return nil, errors.NewAppErrorWithDetails(409, "Username already taken", "Username must be unique", nil)
		}
		user.Username = req.Username
	}

	if req.Email != "" {
		// Check if email is already taken by another user
		exists, err := s.userRepo.ExistsByEmail(req.Email)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to check email existence")
		}
		if exists && user.Email != req.Email {
			return nil, errors.NewAppErrorWithDetails(409, "Email already taken", "Email must be unique", nil)
		}
		user.Email = req.Email
	}

	user.UpdatedAt = time.Now()

	// Update user
	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.WrapError(err, "Failed to update user")
	}

	// Clear password from response
	user.Password = ""

	return user, nil
}

// DeleteUser deletes a user
func (s *userService) DeleteUser(id uuid.UUID) error {
	// Check if user exists
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return errors.WrapError(err, "Failed to get user")
	}
	if user == nil {
		return errors.ErrUserNotFound
	}

	// Delete user
	if err := s.userRepo.Delete(id); err != nil {
		return errors.WrapError(err, "Failed to delete user")
	}

	return nil
}

// RefreshToken refreshes an access token using a refresh token
func (s *userService) RefreshToken(req *models.RefreshTokenRequest) (*models.LoginResponse, error) {
	// Validate refresh token
	claims, err := s.jwtMgr.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.WrapError(err, "Invalid refresh token")
	}

	// Get user
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, errors.WrapError(err, "Failed to get user")
	}
	if user == nil {
		return nil, errors.NewErrorWithCode(404, "User not found")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.NewErrorWithCode(403, "Account is deactivated")
	}

	// Generate new token pair
	tokenPair, err := s.jwtMgr.GenerateTokenPair(user)
	if err != nil {
		return nil, errors.WrapError(err, "Failed to generate token")
	}

	return &models.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
		User:         *user,
	}, nil
}

// Logout logs out a user (in a real implementation, you might want to blacklist the token)
func (s *userService) Logout(userID uuid.UUID, tokenID string) error {
	// In a real implementation, you would:
	// 1. Add the token to a blacklist
	// 2. Remove refresh tokens from database
	// 3. Log the logout event

	// For now, we'll just log the event
	// In a production system, you'd want to implement proper token revocation
	return nil
}

// ValidateUser validates a user entity
func (s *userService) ValidateUser(user *models.User) error {
	return s.validator.Validate(user)
}

// ActivateUser activates a user account
func (s *userService) ActivateUser(id uuid.UUID) error {
	// Check if user exists
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return errors.WrapError(err, "Failed to get user")
	}
	if user == nil {
		return errors.ErrUserNotFound
	}

	// Activate user
	if err := s.userRepo.Activate(id); err != nil {
		return errors.WrapError(err, "Failed to activate user")
	}

	return nil
}

// DeactivateUser deactivates a user account
func (s *userService) DeactivateUser(id uuid.UUID) error {
	// Check if user exists
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return errors.WrapError(err, "Failed to get user")
	}
	if user == nil {
		return errors.ErrUserNotFound
	}

	// Deactivate user
	if err := s.userRepo.Deactivate(id); err != nil {
		return errors.WrapError(err, "Failed to deactivate user")
	}

	return nil
}

// AuthenticateUser authenticates a user with email and password
func (s *userService) AuthenticateUser(req *models.LoginRequest) (*models.LoginResponse, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, errors.WrapErrorWithCode(err, 400, "Validation failed")
	}

	// Get user by email
	user, err := s.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	// Get user with password for authentication
	userWithPassword, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.WrapError(err, "Failed to get user")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(userWithPassword.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.ErrUnauthorized
	}

	// Check if user is active
	if !userWithPassword.IsActive {
		return nil, errors.NewErrorWithCode(403, "Account is deactivated")
	}

	// Generate JWT token pair
	tokenPair, err := s.jwtMgr.GenerateTokenPair(user)
	if err != nil {
		return nil, errors.WrapError(err, "Failed to generate token")
	}

	return &models.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
		User:         *user,
	}, nil
}
