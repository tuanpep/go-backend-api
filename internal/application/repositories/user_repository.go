package repositories

import (
	"database/sql"
	"time"

	"go-backend-api/internal/domain/entities"
	"go-backend-api/internal/pkg/errors"

	"github.com/google/uuid"
)

// userRepository implements UserRepository interface
type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) entities.UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user
func (r *userRepository) Create(user *entities.User) error {
	query := `INSERT INTO users (username, email, password, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := r.db.QueryRow(query, user.Username, user.Email, user.Password, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
	if err != nil {
		return errors.WrapError(err, "Failed to create user")
	}

	return nil
}

// GetByID gets a user by ID
func (r *userRepository) GetByID(id uuid.UUID) (*entities.User, error) {
	user := &entities.User{}
	query := `SELECT id, username, email, password, created_at, updated_at FROM users WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.WrapError(err, "Failed to get user by ID")
	}

	return user, nil
}

// GetByEmail gets a user by email
func (r *userRepository) GetByEmail(email string) (*entities.User, error) {
	user := &entities.User{}
	query := `SELECT id, username, email, password, created_at, updated_at FROM users WHERE email = $1`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.WrapError(err, "Failed to get user by email")
	}

	return user, nil
}

// GetByUsername gets a user by username
func (r *userRepository) GetByUsername(username string) (*entities.User, error) {
	user := &entities.User{}
	query := `SELECT id, username, email, password, created_at, updated_at FROM users WHERE username = $1`

	err := r.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.WrapError(err, "Failed to get user by username")
	}

	return user, nil
}

// Update updates a user
func (r *userRepository) Update(user *entities.User) error {
	query := `UPDATE users SET username = $1, email = $2, is_active = $3, last_login = $4, updated_at = $5 WHERE id = $6`

	_, err := r.db.Exec(query, user.Username, user.Email, user.IsActive, user.LastLogin, user.UpdatedAt, user.ID)
	if err != nil {
		return errors.WrapError(err, "Failed to update user")
	}

	return nil
}

// Delete deletes a user
func (r *userRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return errors.WrapError(err, "Failed to delete user")
	}

	return nil
}

// UpdateLastLogin updates the last login time for a user
func (r *userRepository) UpdateLastLogin(id uuid.UUID) error {
	query := `UPDATE users SET last_login = $1, updated_at = $2 WHERE id = $3`
	now := time.Now()

	_, err := r.db.Exec(query, now, now, id)
	if err != nil {
		return errors.WrapError(err, "Failed to update last login")
	}

	return nil
}

// Deactivate deactivates a user account
func (r *userRepository) Deactivate(id uuid.UUID) error {
	query := `UPDATE users SET is_active = false, updated_at = $1 WHERE id = $2`
	now := time.Now()

	_, err := r.db.Exec(query, now, id)
	if err != nil {
		return errors.WrapError(err, "Failed to deactivate user")
	}

	return nil
}

// ExistsByEmail checks if a user exists with the given email
func (r *userRepository) ExistsByEmail(email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	err := r.db.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return false, errors.WrapError(err, "Failed to check user existence by email")
	}

	return exists, nil
}

// ExistsByUsername checks if a user exists with the given username
func (r *userRepository) ExistsByUsername(username string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`

	err := r.db.QueryRow(query, username).Scan(&exists)
	if err != nil {
		return false, errors.WrapError(err, "Failed to check user existence by username")
	}

	return exists, nil
}
