package utils

import (
	"crypto/rand"
	"fmt"
	"regexp"
	"strings"
	"time"

	appErrors "github.com/alpinesboltltd/boltz-ai/internal/errors"
	"golang.org/x/crypto/bcrypt"
)

func ValidateEmail(email string) error {
	if email == "" {
		return appErrors.NewValidationError("Email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return appErrors.NewValidationError("Invalid email format")
	}

	return nil
}

func ValidatePassword(password string) error {
	if password == "" {
		return appErrors.NewValidationError("Password is required")
	}

	if len(password) < 6 {
		return appErrors.NewValidationError("Password must be at least 6 characters long")
	}

	return nil
}

func ValidateRequired(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return appErrors.NewValidationError(fieldName + " is required")
	}
	return nil
}

func ValidateTokenPurpose(purpose string) error {
	switch purpose {
	case "login", "register", "password_reset", "2fa":
		return nil
	default:
		return appErrors.NewValidationError("invalid token purpose")
	}
}

// GenerateOTP generates a cryptographically secure numeric OTP of the specified length (default 6 if invalid)
func GenerateOTP(length int) (string, error) {
	if length < 4 || length > 10 { // sanity bounds
		length = 6
	}
	// Each digit needs ~3.32 bits; we will generate length bytes and map to digits for uniformity
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	digits := make([]byte, length)
	for i := 0; i < length; i++ {
		digits[i] = '0' + (b[i] % 10)
	}
	return string(digits), nil
}

// OTPConfig simple configuration for OTP validation
type OTPConfig struct {
	MaxAttempts int
	TTL         time.Duration
}

// ValidateOTPInput validates common OTP request parameters
func ValidateOTPInput(email, purpose, otp string) error {
	if err := ValidateEmail(email); err != nil {
		return err
	}
	if err := ValidateTokenPurpose(purpose); err != nil {
		return err
	}
	if strings.TrimSpace(otp) == "" {
		return appErrors.NewValidationError("otp is required")
	}
	if len(otp) < 4 || len(otp) > 10 {
		return appErrors.NewValidationError("invalid otp length")
	}
	if !regexp.MustCompile(`^[0-9]+$`).MatchString(otp) {
		return appErrors.NewValidationError("otp must be numeric")
	}
	return nil
}

// FormatOTPPurpose helper to namespace OTP purposes (e.g., adds prefix)
func FormatOTPPurpose(base string) string { return fmt.Sprintf("otp_%s", base) }

func CreateHash(token []byte) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword(token, bcrypt.DefaultCost)
	if err != nil {
	}
	return string(bytes), err
}

func ValidateHash(token, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(token))
}
