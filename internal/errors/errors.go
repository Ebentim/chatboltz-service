package errors

import (
	"fmt"
	"log"
	"net/http"
)

type ErrorType string

const (
	ValidationError     ErrorType = "VALIDATION_ERROR"
	AuthenticationError ErrorType = "AUTHENTICATION_ERROR"
	AuthorizationError  ErrorType = "AUTHORIZATION_ERROR"
	NotFoundError       ErrorType = "NOT_FOUND_ERROR"
	ConflictError       ErrorType = "CONFLICT_ERROR"
	DatabaseError       ErrorType = "DATABASE_ERROR"
	ExternalAPIError    ErrorType = "EXTERNAL_API_ERROR"
	InternalError       ErrorType = "INTERNAL_ERROR"
)

type AppError struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
	Code    int       `json:"code"`
	Details string    `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewValidationError(message string) *AppError {
	return &AppError{
		Type:    ValidationError,
		Message: message,
		Code:    http.StatusBadRequest,
	}
}

func NewAuthenticationError(message string) *AppError {
	return &AppError{
		Type:    AuthenticationError,
		Message: message,
		Code:    http.StatusUnauthorized,
	}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{
		Type:    NotFoundError,
		Message: message,
		Code:    http.StatusNotFound,
	}
}

func NewConflictError(message string) *AppError {
	return &AppError{
		Type:    ConflictError,
		Message: message,
		Code:    http.StatusConflict,
	}
}

func NewDatabaseError(message string, details string) *AppError {
	return &AppError{
		Type:    DatabaseError,
		Message: message,
		Code:    http.StatusInternalServerError,
		Details: details,
	}
}

func NewExternalAPIError(message string, details string) *AppError {
	return &AppError{
		Type:    ExternalAPIError,
		Message: message,
		Code:    http.StatusBadGateway,
		Details: details,
	}
}

func NewInternalError(message string, details string) *AppError {
	return &AppError{
		Type:    InternalError,
		Message: message,
		Code:    http.StatusInternalServerError,
		Details: details,
	}
}

func LogError(err error, context string) {
	if appErr, ok := err.(*AppError); ok {
		log.Printf("[ERROR] %s: %s (Type: %s, Code: %d)", context, appErr.Message, appErr.Type, appErr.Code)
		if appErr.Details != "" {
			log.Printf("[ERROR] Details: %s", appErr.Details)
		}
	} else {
		log.Printf("[ERROR] %s: %s", context, err.Error())
	}
}

func WrapDatabaseError(err error, operation string) *AppError {
	details := fmt.Sprintf("Database operation '%s' failed: %s", operation, err.Error())
	LogError(err, fmt.Sprintf("Database Error - %s", operation))
	return NewDatabaseError("Database operation failed", details)
}

func WrapExternalAPIError(err error, service string) *AppError {
	details := fmt.Sprintf("External API '%s' failed: %s", service, err.Error())
	LogError(err, fmt.Sprintf("External API Error - %s", service))
	return NewExternalAPIError(fmt.Sprintf("%s service unavailable", service), details)
}
