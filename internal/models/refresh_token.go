package models

import (
	"time"

	"github.com/google/uuid"
)

// RefreshToken represents a refresh token entity
type RefreshToken struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	TokenID   string     `json:"token_id" db:"token_id"`
	TokenHash string     `json:"-" db:"token_hash"`
	ExpiresAt time.Time  `json:"expires_at" db:"expires_at"`
	IsRevoked bool       `json:"is_revoked" db:"is_revoked"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty" db:"revoked_at"`
}

// RefreshTokenRepository defines the interface for refresh token data operations
type RefreshTokenRepository interface {
	Create(tokenID, tokenHash string, userID uuid.UUID, expiresAt time.Time) error
	GetByTokenID(tokenID string) (*RefreshToken, error)
	Revoke(tokenID string) error
	RevokeAllForUser(userID uuid.UUID) error
	IsValid(tokenID string) (bool, error)
	IsValidWithLock(tokenID string) (bool, error)
	RotateToken(oldTokenID, newTokenID, newTokenHash string, userID uuid.UUID, expiresAt time.Time) error
	DeleteExpired() error
}
