package errors

import (
	"fmt"
	"net/http"
)

// AppError represents an application error
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
	Err     error  `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new application error
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// NewAppErrorWithDetails creates a new application error with details
func NewAppErrorWithDetails(code int, message, details string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
		Err:     err,
	}
}

// Predefined errors
var (
	// Authentication errors
	ErrUnauthorized = NewAppError(http.StatusUnauthorized, "Unauthorized", nil)
	ErrForbidden    = NewAppError(http.StatusForbidden, "Forbidden", nil)
	ErrInvalidToken = NewAppError(http.StatusUnauthorized, "Invalid token", nil)
	ErrTokenExpired = NewAppError(http.StatusUnauthorized, "Token expired", nil)

	// Validation errors
	ErrInvalidInput = NewAppError(http.StatusBadRequest, "Invalid input", nil)
	ErrValidation   = NewAppError(http.StatusBadRequest, "Validation failed", nil)

	// Not found errors
	ErrNotFound     = NewAppError(http.StatusNotFound, "Resource not found", nil)
	ErrUserNotFound = NewAppError(http.StatusNotFound, "User not found", nil)
	ErrPostNotFound = NewAppError(http.StatusNotFound, "Post not found", nil)

	// Conflict errors
	ErrConflict   = NewAppError(http.StatusConflict, "Resource already exists", nil)
	ErrUserExists = NewAppError(http.StatusConflict, "User already exists", nil)

	// Internal errors
	ErrInternal = NewAppError(http.StatusInternalServerError, "Internal server error", nil)
	ErrDatabase = NewAppError(http.StatusInternalServerError, "Database error", nil)
)

// WrapError wraps an existing error with additional context
func WrapError(err error, message string) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return NewAppError(http.StatusInternalServerError, message, err)
}

// WrapErrorWithCode wraps an existing error with additional context and custom code
func WrapErrorWithCode(err error, code int, message string) *AppError {
	if appErr, ok := err.(*AppError); ok {
		appErr.Code = code
		appErr.Message = message
		return appErr
	}
	return NewAppError(code, message, err)
}

// NewErrorWithCode creates a new error with a specific HTTP status code
func NewErrorWithCode(code int, message string) *AppError {
	return NewAppError(code, message, nil)
}
