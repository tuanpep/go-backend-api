package security

import (
	"html"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// InputValidator provides enhanced input validation and sanitization
type InputValidator struct {
	validator *validator.Validate
}

// NewInputValidator creates a new input validator
func NewInputValidator() *InputValidator {
	v := validator.New()

	// Register custom validators
	if err := v.RegisterValidation("username", validateUsername); err != nil {
		panic("failed to register username validator: " + err.Error())
	}
	if err := v.RegisterValidation("password", validatePassword); err != nil {
		panic("failed to register password validator: " + err.Error())
	}
	if err := v.RegisterValidation("email", validateEmail); err != nil {
		panic("failed to register email validator: " + err.Error())
	}
	if err := v.RegisterValidation("no_sql_injection", validateNoSQLInjection); err != nil {
		panic("failed to register no_sql_injection validator: " + err.Error())
	}
	if err := v.RegisterValidation("no_xss", validateNoXSS); err != nil {
		panic("failed to register no_xss validator: " + err.Error())
	}

	return &InputValidator{validator: v}
}

// Validate validates a struct
func (iv *InputValidator) Validate(i interface{}) error {
	return iv.validator.Struct(i)
}

// SanitizeString sanitizes a string input
func (iv *InputValidator) SanitizeString(input string) string {
	// HTML escape
	sanitized := html.EscapeString(input)

	// Remove null bytes
	sanitized = strings.ReplaceAll(sanitized, "\x00", "")

	// Trim whitespace
	sanitized = strings.TrimSpace(sanitized)

	return sanitized
}

// SanitizeHTML sanitizes HTML input (more permissive than SanitizeString)
func (iv *InputValidator) SanitizeHTML(input string) string {
	// Remove script tags and their content
	scriptRegex := regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	sanitized := scriptRegex.ReplaceAllString(input, "")

	// Remove javascript: protocols
	jsRegex := regexp.MustCompile(`(?i)javascript:`)
	sanitized = jsRegex.ReplaceAllString(sanitized, "")

	// Remove on* event handlers
	eventRegex := regexp.MustCompile(`(?i)\s*on\w+\s*=\s*"[^"]*"`)
	sanitized = eventRegex.ReplaceAllString(sanitized, "")

	// HTML escape remaining content
	sanitized = html.EscapeString(sanitized)

	return sanitized
}

// validateUsername validates username format and security
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()

	// Length check
	if len(username) < 3 || len(username) > 20 {
		return false
	}

	// Only alphanumeric and underscores
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, username)
	if !matched {
		return false
	}

	// Check for reserved usernames
	reserved := []string{"admin", "root", "administrator", "api", "www", "mail", "ftp", "test"}
	for _, reserved := range reserved {
		if strings.EqualFold(username, reserved) {
			return false
		}
	}

	return true
}

// validatePassword validates password strength
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Length check
	if len(password) < 8 {
		return false
	}

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

	// Require at least 3 of 4 character types
	types := 0
	if hasUpper {
		types++
	}
	if hasLower {
		types++
	}
	if hasNumber {
		types++
	}
	if hasSpecial {
		types++
	}

	return types >= 3
}

// validateEmail validates email format
func validateEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()

	// Basic email regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return false
	}

	// Check for suspicious patterns
	suspiciousPatterns := []string{
		"..", "++", "--", "__",
	}

	for _, pattern := range suspiciousPatterns {
		if strings.Contains(email, pattern) {
			return false
		}
	}

	return true
}

// validateNoSQLInjection validates against SQL injection patterns
func validateNoSQLInjection(fl validator.FieldLevel) bool {
	input := fl.Field().String()

	// Common SQL injection patterns
	sqlPatterns := []string{
		"(?i)(union|select|insert|update|delete|drop|create|alter|exec|execute)",
		"(?i)(or|and)\\s+\\d+\\s*=\\s*\\d+",
		"(?i)(or|and)\\s+'.*'\\s*=\\s*'.*'",
		"(?i)(or|and)\\s+\".*\"\\s*=\\s*\".*\"",
		"(?i)(or|and)\\s+[a-zA-Z_][a-zA-Z0-9_]*\\s*=\\s*[a-zA-Z_][a-zA-Z0-9_]*",
		"(?i)(or|and)\\s+[a-zA-Z_][a-zA-Z0-9_]*\\s*like\\s+",
		"(?i)(or|and)\\s+[a-zA-Z_][a-zA-Z0-9_]*\\s*in\\s*\\(",
		"(?i)(or|and)\\s+[a-zA-Z_][a-zA-Z0-9_]*\\s*between\\s+",
		"(?i)(or|and)\\s+[a-zA-Z_][a-zA-Z0-9_]*\\s*is\\s+null",
		"(?i)(or|and)\\s+[a-zA-Z_][a-zA-Z0-9_]*\\s*is\\s+not\\s+null",
		"(?i)(or|and)\\s+[a-zA-Z_][a-zA-Z0-9_]*\\s*exists\\s*\\(",
		"(?i)(or|and)\\s+[a-zA-Z_][a-zA-Z0-9_]*\\s*not\\s+exists\\s*\\(",
		"(?i)(or|and)\\s+[a-zA-Z_][a-zA-Z0-9_]*\\s*in\\s*\\(",
		"(?i)(or|and)\\s+[a-zA-Z_][a-zA-Z0-9_]*\\s*not\\s+in\\s*\\(",
		"(?i)(or|and)\\s+[a-zA-Z_][a-zA-Z0-9_]*\\s*like\\s+",
		"(?i)(or|and)\\s+[a-zA-Z_][a-zA-Z0-9_]*\\s*not\\s+like\\s+",
		"(?i)(or|and)\\s+[a-zA-Z_][a-zA-Z0-9_]*\\s*regexp",
		"(?i)(or|and)\\s+[a-zA-Z_][a-zA-Z0-9_]*\\s*not\\s+regexp",
		"(?i)(or|and)\\s+[a-zA-Z_][a-zA-Z0-9_]*\\s*similar\\s+to",
		"(?i)(or|and)\\s+[a-zA-Z_][a-zA-Z0-9_]*\\s*not\\s+similar\\s+to",
	}

	for _, pattern := range sqlPatterns {
		matched, _ := regexp.MatchString(pattern, input)
		if matched {
			return false
		}
	}

	return true
}

// validateNoXSS validates against XSS patterns
func validateNoXSS(fl validator.FieldLevel) bool {
	input := fl.Field().String()

	// Common XSS patterns
	xssPatterns := []string{
		"(?i)<script[^>]*>.*?</script>",
		"(?i)<iframe[^>]*>.*?</iframe>",
		"(?i)<object[^>]*>.*?</object>",
		"(?i)<embed[^>]*>.*?</embed>",
		"(?i)<applet[^>]*>.*?</applet>",
		"(?i)<form[^>]*>.*?</form>",
		"(?i)<input[^>]*>",
		"(?i)<textarea[^>]*>.*?</textarea>",
		"(?i)<select[^>]*>.*?</select>",
		"(?i)<option[^>]*>.*?</option>",
		"(?i)<button[^>]*>.*?</button>",
		"(?i)<link[^>]*>",
		"(?i)<meta[^>]*>",
		"(?i)<style[^>]*>.*?</style>",
		"(?i)<link[^>]*>",
		"(?i)javascript:",
		"(?i)vbscript:",
		"(?i)onload\\s*=",
		"(?i)onerror\\s*=",
		"(?i)onclick\\s*=",
		"(?i)onmouseover\\s*=",
		"(?i)onfocus\\s*=",
		"(?i)onblur\\s*=",
		"(?i)onchange\\s*=",
		"(?i)onsubmit\\s*=",
		"(?i)onreset\\s*=",
		"(?i)onselect\\s*=",
		"(?i)onkeydown\\s*=",
		"(?i)onkeyup\\s*=",
		"(?i)onkeypress\\s*=",
		"(?i)onmousedown\\s*=",
		"(?i)onmouseup\\s*=",
		"(?i)onmousemove\\s*=",
		"(?i)onmouseout\\s*=",
		"(?i)onmouseenter\\s*=",
		"(?i)onmouseleave\\s*=",
		"(?i)oncontextmenu\\s*=",
		"(?i)ondblclick\\s*=",
		"(?i)onwheel\\s*=",
		"(?i)onabort\\s*=",
		"(?i)oncanplay\\s*=",
		"(?i)oncanplaythrough\\s*=",
		"(?i)ondurationchange\\s*=",
		"(?i)onemptied\\s*=",
		"(?i)onended\\s*=",
		"(?i)onerror\\s*=",
		"(?i)onloadeddata\\s*=",
		"(?i)onloadedmetadata\\s*=",
		"(?i)onloadstart\\s*=",
		"(?i)onpause\\s*=",
		"(?i)onplay\\s*=",
		"(?i)onplaying\\s*=",
		"(?i)onprogress\\s*=",
		"(?i)onratechange\\s*=",
		"(?i)onseeked\\s*=",
		"(?i)onseeking\\s*=",
		"(?i)onstalled\\s*=",
		"(?i)onsuspend\\s*=",
		"(?i)ontimeupdate\\s*=",
		"(?i)onvolumechange\\s*=",
		"(?i)onwaiting\\s*=",
	}

	for _, pattern := range xssPatterns {
		matched, _ := regexp.MatchString(pattern, input)
		if matched {
			return false
		}
	}

	return true
}
