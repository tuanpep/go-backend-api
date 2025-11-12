package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"go-backend-api/internal/models"

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
func (j *JWTManager) GenerateTokenPair(user *models.User) (*TokenPair, error) {
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
func (j *JWTManager) generateAccessToken(user *models.User, tokenID string) (string, error) {
	claims := &models.TokenClaims{
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
func (j *JWTManager) generateRefreshToken(user *models.User, tokenID string) (string, error) {
	claims := &models.TokenClaims{
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
func (j *JWTManager) ValidateAccessToken(tokenString string) (*models.TokenClaims, error) {
	return j.validateToken(tokenString, j.accessSecretKey, "access")
}

// ValidateRefreshToken validates a refresh token
func (j *JWTManager) ValidateRefreshToken(tokenString string) (*models.TokenClaims, error) {
	return j.validateToken(tokenString, j.refreshSecretKey, "refresh")
}

// validateToken validates a JWT token
func (j *JWTManager) validateToken(tokenString, secretKey, expectedType string) (*models.TokenClaims, error) {
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

	return &models.TokenClaims{
		UserID:   userID,
		Username: username,
		TokenID:  tokenID,
		Type:     tokenType,
	}, nil
}

// generateTokenID generates a cryptographically secure random token ID
func generateTokenID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GetRefreshDuration returns the refresh token duration
func (j *JWTManager) GetRefreshDuration() time.Duration {
	return j.refreshDuration
}

// HashRefreshToken hashes a refresh token using SHA256
func HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
