package validation

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator wraps the go-playground validator
type Validator struct {
	validator *validator.Validate
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	v := validator.New()

	// Register custom validators
	v.RegisterValidation("username", validateUsername)
	v.RegisterValidation("password", validatePassword)

	return &Validator{validator: v}
}

// Validate validates a struct
func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

// ValidateVar validates a single variable
func (v *Validator) ValidateVar(field interface{}, tag string) error {
	return v.validator.Var(field, tag)
}

// validateUsername validates username format
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()

	// Username should be 3-20 characters, alphanumeric and underscores only
	if len(username) < 3 || len(username) > 20 {
		return false
	}

	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, username)
	return matched
}

// validatePassword validates password strength
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Password should be at least 6 characters
	if len(password) < 6 {
		return false
	}

	// Password should contain at least one letter and one number
	hasLetter, _ := regexp.MatchString(`[a-zA-Z]`, password)
	hasNumber, _ := regexp.MatchString(`[0-9]`, password)

	return hasLetter && hasNumber
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

// GetValidationErrors extracts validation errors
func (v *Validator) GetValidationErrors(err error) []ValidationError {
	var validationErrors []ValidationError

	if validationErr, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErr {
			validationErrors = append(validationErrors, ValidationError{
				Field:   strings.ToLower(e.Field()),
				Tag:     e.Tag(),
				Value:   e.Param(),
				Message: getValidationMessage(e),
			})
		}
	}

	return validationErrors
}

// getValidationMessage returns a human-readable validation message
func getValidationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Must be a valid email address"
	case "min":
		return "Must be at least " + e.Param() + " characters"
	case "max":
		return "Must be no more than " + e.Param() + " characters"
	case "username":
		return "Username must be 3-20 characters, alphanumeric and underscores only"
	case "password":
		return "Password must be at least 6 characters with letters and numbers"
	default:
		return "Invalid value"
	}
}
