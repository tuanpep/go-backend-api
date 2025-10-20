package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// PasswordPolicy defines password requirements
type PasswordPolicy struct {
	MinLength        int
	MaxLength        int
	RequireUppercase bool
	RequireLowercase bool
	RequireNumbers   bool
	RequireSpecial   bool
	ForbiddenWords   []string
	MaxConsecutive   int
}

// DefaultPasswordPolicy returns the default password policy
func DefaultPasswordPolicy() *PasswordPolicy {
	return &PasswordPolicy{
		MinLength:        8,
		MaxLength:        128,
		RequireUppercase: true,
		RequireLowercase: true,
		RequireNumbers:   true,
		RequireSpecial:   true,
		ForbiddenWords: []string{
			"password", "123456", "qwerty", "abc123", "password123",
			"admin", "root", "user", "test", "guest", "demo",
		},
		MaxConsecutive: 3,
	}
}

// ValidatePassword validates a password against the policy
func (pp *PasswordPolicy) ValidatePassword(password string) error {
	// Length check
	if len(password) < pp.MinLength {
		return fmt.Errorf("password must be at least %d characters long", pp.MinLength)
	}
	if len(password) > pp.MaxLength {
		return fmt.Errorf("password must be no more than %d characters long", pp.MaxLength)
	}

	// Character type checks
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if pp.RequireUppercase && !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if pp.RequireLowercase && !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if pp.RequireNumbers && !hasNumber {
		return fmt.Errorf("password must contain at least one number")
	}
	if pp.RequireSpecial && !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}

	// Check for forbidden words
	lowerPassword := strings.ToLower(password)
	for _, word := range pp.ForbiddenWords {
		if strings.Contains(lowerPassword, strings.ToLower(word)) {
			return fmt.Errorf("password contains forbidden word: %s", word)
		}
	}

	// Check for consecutive characters
	if pp.MaxConsecutive > 0 {
		consecutive := 1
		for i := 1; i < len(password); i++ {
			if password[i] == password[i-1] {
				consecutive++
				if consecutive > pp.MaxConsecutive {
					return fmt.Errorf("password cannot contain more than %d consecutive identical characters", pp.MaxConsecutive)
				}
			} else {
				consecutive = 1
			}
		}
	}

	// Check for common patterns
	if err := pp.checkCommonPatterns(password); err != nil {
		return err
	}

	return nil
}

// checkCommonPatterns checks for common weak patterns
func (pp *PasswordPolicy) checkCommonPatterns(password string) error {
	// Sequential characters
	sequential := []string{
		"abcdefghijklmnopqrstuvwxyz",
		"zyxwvutsrqponmlkjihgfedcba",
		"0123456789",
		"9876543210",
		"qwertyuiop",
		"poiuytrewq",
		"asdfghjkl",
		"lkjhgfdsa",
		"zxcvbnm",
		"mnbvcxz",
	}

	lowerPassword := strings.ToLower(password)
	for _, seq := range sequential {
		for i := 0; i <= len(seq)-4; i++ {
			if strings.Contains(lowerPassword, seq[i:i+4]) {
				return fmt.Errorf("password contains sequential characters")
			}
		}
	}

	// Repeated patterns
	patternRegex := regexp.MustCompile(`(.{2,})\1`)
	if patternRegex.MatchString(password) {
		return fmt.Errorf("password contains repeated patterns")
	}

	return nil
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	// Use a higher cost for better security
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedBytes), nil
}

// VerifyPassword verifies a password against its hash
func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// GenerateSecurePassword generates a secure random password
func GenerateSecurePassword(length int) (string, error) {
	if length < 8 {
		length = 8
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	password := make([]byte, length)

	for i := range password {
		randomBytes := make([]byte, 1)
		if _, err := rand.Read(randomBytes); err != nil {
			return "", fmt.Errorf("failed to generate random password: %w", err)
		}
		password[i] = charset[randomBytes[0]%byte(len(charset))]
	}

	return string(password), nil
}

// GenerateSalt generates a random salt
func GenerateSalt() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// ConstantTimeCompare compares two strings in constant time
func ConstantTimeCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

// PasswordStrength calculates password strength score (0-100)
func PasswordStrength(password string) int {
	score := 0

	// Length score (max 25 points)
	length := len(password)
	if length >= 8 {
		score += 10
	}
	if length >= 12 {
		score += 10
	}
	if length >= 16 {
		score += 5
	}

	// Character variety score (max 50 points)
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if hasUpper {
		score += 10
	}
	if hasLower {
		score += 10
	}
	if hasNumber {
		score += 10
	}
	if hasSpecial {
		score += 10
	}

	// Complexity score (max 25 points)
	uniqueChars := make(map[rune]bool)
	for _, char := range password {
		uniqueChars[char] = true
	}
	uniqueCount := len(uniqueChars)

	if uniqueCount >= 8 {
		score += 10
	}
	if uniqueCount >= 12 {
		score += 10
	}
	if uniqueCount >= 16 {
		score += 5
	}

	// Penalty for common patterns
	if regexp.MustCompile(`(?i)(password|123456|qwerty)`).MatchString(password) {
		score -= 20
	}

	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}

// AccountLockout represents account lockout information
type AccountLockout struct {
	Attempts    int
	LastAttempt time.Time
	LockedUntil *time.Time
	IsLocked    bool
}

// NewAccountLockout creates a new account lockout tracker
func NewAccountLockout() *AccountLockout {
	return &AccountLockout{
		Attempts:    0,
		LastAttempt: time.Time{},
		LockedUntil: nil,
		IsLocked:    false,
	}
}

// RecordFailedAttempt records a failed login attempt
func (al *AccountLockout) RecordFailedAttempt(maxAttempts int, lockoutDuration time.Duration) {
	al.Attempts++
	al.LastAttempt = time.Now()

	if al.Attempts >= maxAttempts {
		lockoutUntil := time.Now().Add(lockoutDuration)
		al.LockedUntil = &lockoutUntil
		al.IsLocked = true
	}
}

// RecordSuccessfulAttempt resets the lockout counter
func (al *AccountLockout) RecordSuccessfulAttempt() {
	al.Attempts = 0
	al.LastAttempt = time.Time{}
	al.LockedUntil = nil
	al.IsLocked = false
}

// IsCurrentlyLocked checks if the account is currently locked
func (al *AccountLockout) IsCurrentlyLocked() bool {
	if !al.IsLocked {
		return false
	}

	if al.LockedUntil != nil && time.Now().After(*al.LockedUntil) {
		// Lockout period has expired
		al.IsLocked = false
		al.LockedUntil = nil
		al.Attempts = 0
		return false
	}

	return true
}

// GetRemainingLockoutTime returns the remaining lockout time
func (al *AccountLockout) GetRemainingLockoutTime() time.Duration {
	if !al.IsLocked || al.LockedUntil == nil {
		return 0
	}

	remaining := time.Until(*al.LockedUntil)
	if remaining < 0 {
		return 0
	}

	return remaining
}
