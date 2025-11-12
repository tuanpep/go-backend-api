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
	userRepo         models.UserRepository
	refreshTokenRepo models.RefreshTokenRepository
	jwtMgr           *auth.JWTManager
	validator        *validation.Validator
}

// NewUserService creates a new user service
func NewUserService(userRepo models.UserRepository, refreshTokenRepo models.RefreshTokenRepository, jwtMgr *auth.JWTManager) models.UserService {
	return &userService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtMgr:           jwtMgr,
		validator:        validation.NewValidator(),
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

// RefreshToken refreshes an access token using a refresh token with rotation
func (s *userService) RefreshToken(req *models.RefreshTokenRequest) (*models.LoginResponse, error) {
	// Step 1: Validate refresh token JWT signature and claims
	claims, err := s.jwtMgr.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		// Generic error message - don't reveal token state
		return nil, errors.NewErrorWithCode(401, "Invalid refresh token")
	}

	// Step 2: Get user
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, errors.WrapError(err, "Failed to get user")
	}
	if user == nil {
		// Generic error message - don't reveal user existence
		return nil, errors.NewErrorWithCode(401, "Invalid refresh token")
	}

	// Step 3: Check if user is active
	if !user.IsActive {
		return nil, errors.NewErrorWithCode(403, "Account is deactivated")
	}

	// Step 4: Generate new token pair (with new token_id)
	tokenPair, err := s.jwtMgr.GenerateTokenPair(user)
	if err != nil {
		return nil, errors.WrapError(err, "Failed to generate token")
	}

	// Step 5: Extract token_id from new refresh token
	newRefreshClaims, err := s.jwtMgr.ValidateRefreshToken(tokenPair.RefreshToken)
	if err != nil {
		return nil, errors.WrapError(err, "Failed to validate generated refresh token")
	}

	// Step 6: Hash new refresh token
	tokenHash := auth.HashRefreshToken(tokenPair.RefreshToken)
	expiresAt := time.Now().Add(s.jwtMgr.GetRefreshDuration())

	// Step 7: Atomically rotate token (validate old token with lock, create new, revoke old)
	// This prevents race conditions and ensures atomicity
	err = s.refreshTokenRepo.RotateToken(claims.TokenID, newRefreshClaims.TokenID, tokenHash, user.ID, expiresAt)
	if err != nil {
		// Generic error message - don't reveal why token is invalid
		return nil, errors.NewErrorWithCode(401, "Invalid refresh token")
	}

	return &models.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
		User:         *user,
	}, nil
}

// Logout logs out a user by revoking the refresh token
func (s *userService) Logout(userID uuid.UUID, tokenID string) error {
	// Revoke the refresh token associated with this token_id
	if err := s.refreshTokenRepo.Revoke(tokenID); err != nil {
		return errors.WrapError(err, "Failed to revoke refresh token")
	}

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

	// Extract token_id from refresh token to store in database
	refreshClaims, err := s.jwtMgr.ValidateRefreshToken(tokenPair.RefreshToken)
	if err != nil {
		return nil, errors.WrapError(err, "Failed to validate generated refresh token")
	}

	// Hash refresh token for storage
	tokenHash := auth.HashRefreshToken(tokenPair.RefreshToken)

	// Calculate expiration time from refresh token duration
	expiresAt := time.Now().Add(s.jwtMgr.GetRefreshDuration())

	// Store refresh token in database
	if err := s.refreshTokenRepo.Create(refreshClaims.TokenID, tokenHash, user.ID, expiresAt); err != nil {
		return nil, errors.WrapError(err, "Failed to store refresh token")
	}

	return &models.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
		User:         *user,
	}, nil
}
