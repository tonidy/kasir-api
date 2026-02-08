package model

import "fmt"

type Product struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Price      int       `json:"price"`
	Stock      int       `json:"stock"`
	Active     bool      `json:"active"`
	CategoryID *int      `json:"category_id,omitempty"`
	Category   *Category `json:"category,omitempty"`
}

func (p Product) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("%w: name is required", ErrValidation)
	}
	if p.Price < 0 {
		return fmt.Errorf("%w: price must be positive", ErrValidation)
	}
	if p.Stock < 0 {
		return fmt.Errorf("%w: stock must be positive", ErrValidation)
	}
	return nil
}
