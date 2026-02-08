package model

import (
	"reflect"

	errorsPkg "kasir-api/pkg/errors"
	"kasir-api/pkg/validation"
)

type Product struct {
	ID         int       `json:"id" validate:"omitempty,min=1"`
	Name       string    `json:"name" validate:"required,min=1,max=255"`
	Price      int       `json:"price" validate:"min=0"`
	Stock      int       `json:"stock" validate:"min=0"`
	Active     bool      `json:"active"`
	CategoryID *int      `json:"category_id,omitempty" validate:"omitempty,min=1"`
	Category   *Category `json:"category,omitempty"`
}

func (p Product) Validate() error {
	validator := validation.NewValidator()

	// Use reflection to check if this is a new product (ID is 0) or an update
	v := reflect.ValueOf(p)
	idField := v.FieldByName("ID")
	_ = idField.IsValid() && idField.Int() == 0

	// For updates, we might want to allow some fields to be optional
	// For now, we'll validate the struct as is
	if err := validator.ValidateStruct(p); err != nil {
		// Convert validation error to our custom error type
		return errorsPkg.ValidationError(err.Error())
	}

	return nil
}
