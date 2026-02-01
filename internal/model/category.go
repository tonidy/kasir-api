package model

import "fmt"

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c Category) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("%w: name is required", ErrValidation)
	}
	return nil
}
