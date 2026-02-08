package model

import (
	"fmt"
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
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type CheckoutRequest struct {
	Items []CheckoutItem `json:"items"`
}

func (c CheckoutRequest) Validate() error {
	if len(c.Items) == 0 {
		return fmt.Errorf("%w: items cannot be empty", ErrValidation)
	}
	for i, item := range c.Items {
		if item.ProductID <= 0 {
			return fmt.Errorf("%w: item[%d] product_id must be positive", ErrValidation, i)
		}
		if item.Quantity <= 0 {
			return fmt.Errorf("%w: item[%d] quantity must be positive", ErrValidation, i)
		}
	}
	return nil
}
