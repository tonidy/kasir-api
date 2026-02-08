package model

import (
	"reflect"

	errorsPkg "kasir-api/pkg/errors"
	"kasir-api/pkg/validation"
)

type Category struct {
	ID          int    `json:"id" validate:"omitempty,min=1"`
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"omitempty,max=500"`
}

func (c Category) Validate() error {
	validator := validation.NewValidator()

	// Use reflection to check if this is a new category (ID is 0) or an update
	v := reflect.ValueOf(c)
	idField := v.FieldByName("ID")
	_ = idField.IsValid() && idField.Int() == 0

	// For updates, we might want to allow some fields to be optional
	// For now, we'll validate the struct as is
	if err := validator.ValidateStruct(c); err != nil {
		// Convert validation error to our custom error type
		return errorsPkg.ValidationError(err.Error())
	}

	return nil
}
