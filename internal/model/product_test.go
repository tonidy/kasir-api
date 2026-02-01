package domain

import (
	"testing"
)

func TestProduct_Validate(t *testing.T) {
	tests := []struct {
		name    string
		product Product
		wantErr bool
	}{
		{
			name: "valid product",
			product: Product{
				Name:  "Indomie",
				Price: 3500,
				Stock: 100,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			product: Product{
				Name:  "",
				Price: 3500,
				Stock: 100,
			},
			wantErr: true,
		},
		{
			name: "negative price",
			product: Product{
				Name:  "Indomie",
				Price: -100,
				Stock: 100,
			},
			wantErr: true,
		},
		{
			name: "negative stock",
			product: Product{
				Name:  "Indomie",
				Price: 3500,
				Stock: -10,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.product.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Product.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
