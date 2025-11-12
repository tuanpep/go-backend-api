package repositories

import (
	"database/sql"
	"time"

	"go-backend-api/internal/models"
	"go-backend-api/internal/pkg/errors"

	"github.com/google/uuid"
)

// refreshTokenRepository implements RefreshTokenRepository interface
type refreshTokenRepository struct {
	db *sql.DB
}

// NewRefreshTokenRepository creates a new refresh token repository
func NewRefreshTokenRepository(db *sql.DB) models.RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

// Create creates a new refresh token record
func (r *refreshTokenRepository) Create(tokenID, tokenHash string, userID uuid.UUID, expiresAt time.Time) error {
	query := `INSERT INTO refresh_tokens (user_id, token_id, token_hash, expires_at, is_revoked, created_at) 
			  VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.Exec(query, userID, tokenID, tokenHash, expiresAt, false, time.Now())
	if err != nil {
		return errors.WrapError(err, "Failed to create refresh token")
	}

	return nil
}

// GetByTokenID gets a refresh token by token_id
func (r *refreshTokenRepository) GetByTokenID(tokenID string) (*models.RefreshToken, error) {
	token := &models.RefreshToken{}
	query := `SELECT id, user_id, token_id, token_hash, expires_at, is_revoked, created_at, revoked_at 
			  FROM refresh_tokens WHERE token_id = $1`

	err := r.db.QueryRow(query, tokenID).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.IsRevoked,
		&token.CreatedAt,
		&token.RevokedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.WrapError(err, "Failed to get refresh token by token_id")
	}

	return token, nil
}

// Revoke revokes a refresh token by token_id
func (r *refreshTokenRepository) Revoke(tokenID string) error {
	query := `UPDATE refresh_tokens SET is_revoked = true, revoked_at = $1 WHERE token_id = $2`

	_, err := r.db.Exec(query, time.Now(), tokenID)
	if err != nil {
		return errors.WrapError(err, "Failed to revoke refresh token")
	}

	return nil
}

// RevokeAllForUser revokes all refresh tokens for a user
func (r *refreshTokenRepository) RevokeAllForUser(userID uuid.UUID) error {
	query := `UPDATE refresh_tokens SET is_revoked = true, revoked_at = $1 WHERE user_id = $2 AND is_revoked = false`

	_, err := r.db.Exec(query, time.Now(), userID)
	if err != nil {
		return errors.WrapError(err, "Failed to revoke all refresh tokens for user")
	}

	return nil
}

// IsValid checks if a refresh token is valid (exists, not revoked, not expired)
func (r *refreshTokenRepository) IsValid(tokenID string) (bool, error) {
	var isValid bool
	query := `SELECT EXISTS(
		SELECT 1 FROM refresh_tokens 
		WHERE token_id = $1 
		AND is_revoked = false 
		AND expires_at > NOW()
	)`

	err := r.db.QueryRow(query, tokenID).Scan(&isValid)
	if err != nil {
		return false, errors.WrapError(err, "Failed to check refresh token validity")
	}

	return isValid, nil
}

// IsValidWithLock checks if a refresh token is valid with row-level locking to prevent race conditions
func (r *refreshTokenRepository) IsValidWithLock(tokenID string) (bool, error) {
	var isValid bool
	query := `SELECT EXISTS(
		SELECT 1 FROM refresh_tokens 
		WHERE token_id = $1 
		AND is_revoked = false 
		AND expires_at > NOW()
		FOR UPDATE
	)`

	err := r.db.QueryRow(query, tokenID).Scan(&isValid)
	if err != nil {
		return false, errors.WrapError(err, "Failed to check refresh token validity")
	}

	return isValid, nil
}

// RotateToken atomically creates a new refresh token and revokes the old one in a transaction
func (r *refreshTokenRepository) RotateToken(oldTokenID, newTokenID, newTokenHash string, userID uuid.UUID, expiresAt time.Time) error {
	tx, err := r.db.Begin()
	if err != nil {
		return errors.WrapError(err, "Failed to begin transaction")
	}
	defer tx.Rollback()

	// First, validate and lock the old token row
	// Use SELECT FOR UPDATE to lock the row and prevent concurrent access
	var tokenID string
	var isRevoked bool
	var expiresAtDB time.Time
	checkQuery := `SELECT token_id, is_revoked, expires_at 
					FROM refresh_tokens 
					WHERE token_id = $1 
					FOR UPDATE`

	err = tx.QueryRow(checkQuery, oldTokenID).Scan(&tokenID, &isRevoked, &expiresAtDB)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.NewErrorWithCode(401, "Invalid refresh token")
		}
		return errors.WrapError(err, "Failed to validate old token")
	}

	// Check if token is valid (not revoked and not expired)
	if isRevoked || expiresAtDB.Before(time.Now()) {
		return errors.NewErrorWithCode(401, "Invalid refresh token")
	}

	// Create new token
	createQuery := `INSERT INTO refresh_tokens (user_id, token_id, token_hash, expires_at, is_revoked, created_at) 
					VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.Exec(createQuery, userID, newTokenID, newTokenHash, expiresAt, false, time.Now())
	if err != nil {
		return errors.WrapError(err, "Failed to create new refresh token")
	}

	// Revoke old token
	revokeQuery := `UPDATE refresh_tokens SET is_revoked = true, revoked_at = $1 WHERE token_id = $2`
	_, err = tx.Exec(revokeQuery, time.Now(), oldTokenID)
	if err != nil {
		return errors.WrapError(err, "Failed to revoke old refresh token")
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return errors.WrapError(err, "Failed to commit transaction")
	}

	return nil
}

// DeleteExpired deletes expired refresh tokens
func (r *refreshTokenRepository) DeleteExpired() error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < NOW() OR (is_revoked = true AND revoked_at < NOW() - INTERVAL '7 days')`

	_, err := r.db.Exec(query)
	if err != nil {
		return errors.WrapError(err, "Failed to delete expired refresh tokens")
	}

	return nil
}
