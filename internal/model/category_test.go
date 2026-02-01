package model

import (
	"testing"
)

func TestCategory_Validate(t *testing.T) {
	tests := []struct {
		name     string
		category Category
		wantErr  bool
	}{
		{
			name: "valid category",
			category: Category{
				Name:        "Food",
				Description: "Food items",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			category: Category{
				Name:        "",
				Description: "Food items",
			},
			wantErr: true,
		},
		{
			name: "valid category without description",
			category: Category{
				Name:        "Food",
				Description: "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.category.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Category.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
