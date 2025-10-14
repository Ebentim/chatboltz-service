package utils

import (
	"regexp"
	"strings"

	appErrors "github.com/alpinesboltltd/boltz-ai/internal/errors"
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
