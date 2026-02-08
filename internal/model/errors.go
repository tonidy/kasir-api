package model

import (
	errorsPkg "kasir-api/pkg/errors"
)

var (
	ErrValidation = errorsPkg.ValidationError("validation error")
	ErrNotFound   = errorsPkg.NotFoundError("not found")
	ErrConflict   = errorsPkg.ConflictError("conflict")
)

// IsValidationError checks if the error is a validation error
func IsValidationError(err error) bool {
	return errorsPkg.IsType(err, errorsPkg.ErrorTypeValidation)
}

// IsNotFoundError checks if the error is a not found error
func IsNotFoundError(err error) bool {
	return errorsPkg.IsType(err, errorsPkg.ErrorTypeNotFound)
}

// IsConflictError checks if the error is a conflict error
func IsConflictError(err error) bool {
	return errorsPkg.IsType(err, errorsPkg.ErrorTypeConflict)
}
