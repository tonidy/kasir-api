package model

import (
	"fmt"

	errorsPkg "kasir-api/pkg/errors"
	"kasir-api/pkg/validation"
	"time"
)

type Transaction struct {
	ID          int                 `json:"id"`
	TotalAmount int                 `json:"total_amount"`
	CreatedAt   time.Time           `json:"created_at"`
	Details     []TransactionDetail `json:"details"`
}

type TransactionDetail struct {
	ID            int    `json:"id"`
	TransactionID int    `json:"transaction_id"`
	ProductID     int    `json:"product_id"`
	ProductName   string `json:"product_name,omitempty"`
	Quantity      int    `json:"quantity"`
	Subtotal      int    `json:"subtotal"`
}

type CheckoutItem struct {
	ProductID int `json:"product_id" validate:"min=1"`
	Quantity  int `json:"quantity" validate:"min=1"`
}

type CheckoutRequest struct {
	Items []CheckoutItem `json:"items" validate:"required,min=1,dive"`
}

func (c CheckoutRequest) Validate() error {
	validator := validation.NewValidator()

	// Validate the struct using the validator
	if err := validator.ValidateStruct(c); err != nil {
		// Convert validation error to our custom error type
		return errorsPkg.ValidationError(err.Error())
	}

	// Additional custom validation if needed
	if len(c.Items) == 0 {
		return errorsPkg.ValidationError("items cannot be empty")
	}

	for i, item := range c.Items {
		if item.ProductID <= 0 {
			return errorsPkg.ValidationError(fmt.Sprintf("item[%d] product_id must be positive", i))
		}
		if item.Quantity <= 0 {
			return errorsPkg.ValidationError(fmt.Sprintf("item[%d] quantity must be positive", i))
		}
	}

	return nil
}
