package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"go-backend-api/internal/domain/entities"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTManager handles JWT operations with enhanced security
type JWTManager struct {
	accessSecretKey  string
	refreshSecretKey string
	accessDuration   time.Duration
	refreshDuration  time.Duration
	issuer           string
	audience         string
}

// TokenPair represents access and refresh token pair
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// NewJWTManager creates a new JWT manager with enhanced security
func NewJWTManager(accessSecret, refreshSecret, issuer, audience string, accessDuration, refreshDuration time.Duration) *JWTManager {
	return &JWTManager{
		accessSecretKey:  accessSecret,
		refreshSecretKey: refreshSecret,
		accessDuration:   accessDuration,
		refreshDuration:  refreshDuration,
		issuer:           issuer,
		audience:         audience,
	}
}

// GenerateTokenPair generates both access and refresh tokens
func (j *JWTManager) GenerateTokenPair(user *entities.User) (*TokenPair, error) {
	// Generate unique token ID for tracking
	tokenID, err := generateTokenID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token ID: %w", err)
	}

	// Generate access token
	accessToken, err := j.generateAccessToken(user, tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := j.generateRefreshToken(user, tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(j.accessDuration.Seconds()),
	}, nil
}

// generateAccessToken creates an access token
func (j *JWTManager) generateAccessToken(user *entities.User, tokenID string) (string, error) {
	claims := &entities.TokenClaims{
		UserID:   user.ID,
		Username: user.Username,
		TokenID:  tokenID,
		Type:     "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  claims.UserID.String(),
		"username": claims.Username,
		"token_id": claims.TokenID,
		"type":     claims.Type,
		"iss":      j.issuer,
		"aud":      j.audience,
		"exp":      time.Now().Add(j.accessDuration).Unix(),
		"iat":      time.Now().Unix(),
		"nbf":      time.Now().Unix(),
	})

	return token.SignedString([]byte(j.accessSecretKey))
}

// generateRefreshToken creates a refresh token
func (j *JWTManager) generateRefreshToken(user *entities.User, tokenID string) (string, error) {
	claims := &entities.TokenClaims{
		UserID:   user.ID,
		Username: user.Username,
		TokenID:  tokenID,
		Type:     "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  claims.UserID.String(),
		"username": claims.Username,
		"token_id": claims.TokenID,
		"type":     claims.Type,
		"iss":      j.issuer,
		"aud":      j.audience,
		"exp":      time.Now().Add(j.refreshDuration).Unix(),
		"iat":      time.Now().Unix(),
		"nbf":      time.Now().Unix(),
	})

	return token.SignedString([]byte(j.refreshSecretKey))
}

// ValidateAccessToken validates an access token
func (j *JWTManager) ValidateAccessToken(tokenString string) (*entities.TokenClaims, error) {
	return j.validateToken(tokenString, j.accessSecretKey, "access")
}

// ValidateRefreshToken validates a refresh token
func (j *JWTManager) ValidateRefreshToken(tokenString string) (*entities.TokenClaims, error) {
	return j.validateToken(tokenString, j.refreshSecretKey, "refresh")
}

// validateToken validates a JWT token
func (j *JWTManager) validateToken(tokenString, secretKey, expectedType string) (*entities.TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Validate token type
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != expectedType {
		return nil, fmt.Errorf("invalid token type")
	}

	// Validate issuer
	if iss, ok := claims["iss"].(string); !ok || iss != j.issuer {
		return nil, fmt.Errorf("invalid issuer")
	}

	// Validate audience
	if aud, ok := claims["aud"].(string); !ok || aud != j.audience {
		return nil, fmt.Errorf("invalid audience")
	}

	// Extract user ID
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid user_id in token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id format: %w", err)
	}

	// Extract username
	username, ok := claims["username"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid username in token")
	}

	// Extract token ID
	tokenID, ok := claims["token_id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid token_id in token")
	}

	return &entities.TokenClaims{
		UserID:   userID,
		Username: username,
		TokenID:  tokenID,
		Type:     tokenType,
	}, nil
}

// RefreshAccessToken generates a new access token from a refresh token
func (j *JWTManager) RefreshAccessToken(refreshTokenString string, user *entities.User) (*TokenPair, error) {
	// Validate refresh token
	claims, err := j.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Verify user ID matches
	if claims.UserID != user.ID {
		return nil, fmt.Errorf("token user mismatch")
	}

	// Generate new token pair with same token ID
	return j.GenerateTokenPair(user)
}

// generateTokenID generates a cryptographically secure random token ID
func generateTokenID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// ExtractTokenID extracts token ID from context (for logout tracking)
func ExtractTokenID(claims jwt.MapClaims) (string, error) {
	tokenID, ok := claims["token_id"].(string)
	if !ok {
		return "", fmt.Errorf("token_id not found in claims")
	}
	return tokenID, nil
}
