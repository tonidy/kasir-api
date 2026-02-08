package errors

import (
	"fmt"
	"net/http"
)

// ErrorType represents the type of error
type ErrorType string

const (
	ErrorTypeValidation   ErrorType = "validation_error"
	ErrorTypeNotFound     ErrorType = "not_found"
	ErrorTypeConflict     ErrorType = "conflict"
	ErrorTypeInternal     ErrorType = "internal_error"
	ErrorTypeForbidden    ErrorType = "forbidden"
	ErrorTypeUnauthorized ErrorType = "unauthorized"
)

// AppError represents an application error
type AppError struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
	Code    int       `json:"-"` // HTTP status code
	Cause   error     `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Cause
}

// New creates a new AppError
func New(errorType ErrorType, message string) *AppError {
	return &AppError{
		Type:    errorType,
		Message: message,
		Code:    getStatusCode(errorType),
	}
}

// Wrap wraps an existing error with additional context
func Wrap(err error, errorType ErrorType, message string) *AppError {
	return &AppError{
		Type:    errorType,
		Message: message,
		Code:    getStatusCode(errorType),
		Cause:   err,
	}
}

// FromHTTPCode creates an AppError from an HTTP status code
func FromHTTPCode(code int, message string) *AppError {
	var errorType ErrorType
	switch code {
	case http.StatusBadRequest:
		errorType = ErrorTypeValidation
	case http.StatusNotFound:
		errorType = ErrorTypeNotFound
	case http.StatusConflict:
		errorType = ErrorTypeConflict
	case http.StatusForbidden:
		errorType = ErrorTypeForbidden
	case http.StatusUnauthorized:
		errorType = ErrorTypeUnauthorized
	default:
		errorType = ErrorTypeInternal
	}

	return &AppError{
		Type:    errorType,
		Message: message,
		Code:    code,
	}
}

// IsType checks if the error is of a specific type
func IsType(err error, errorType ErrorType) bool {
	var appErr AppError
	if As(err, &appErr) {
		return appErr.Type == errorType
	}
	return false
}

// As is a helper to unwrap errors to AppError
func As(err error, target *AppError) bool {
	for {
		if err == nil {
			return false
		}

		if appErr, ok := err.(*AppError); ok {
			*target = *appErr
			return true
		}

		if unwrapped, ok := err.(interface{ Unwrap() error }); ok {
			err = unwrapped.Unwrap()
		} else {
			return false
		}
	}
}

// Helper functions for common error types
func ValidationError(message string) *AppError {
	return New(ErrorTypeValidation, message)
}

func NotFoundError(message string) *AppError {
	return New(ErrorTypeNotFound, message)
}

func ConflictError(message string) *AppError {
	return New(ErrorTypeConflict, message)
}

func InternalError(message string) *AppError {
	return New(ErrorTypeInternal, message)
}

func ForbiddenError(message string) *AppError {
	return New(ErrorTypeForbidden, message)
}

func UnauthorizedError(message string) *AppError {
	return New(ErrorTypeUnauthorized, message)
}

// getStatusCode returns the HTTP status code for an error type
func getStatusCode(errorType ErrorType) int {
	switch errorType {
	case ErrorTypeValidation:
		return http.StatusBadRequest
	case ErrorTypeNotFound:
		return http.StatusNotFound
	case ErrorTypeConflict:
		return http.StatusConflict
	case ErrorTypeForbidden:
		return http.StatusForbidden
	case ErrorTypeUnauthorized:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
