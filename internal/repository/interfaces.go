package repository

import (
	"context"
	"kasir-api/internal/model"
)

// ProductReader defines read operations for products
type ProductReader interface {
	FindByID(ctx context.Context, id int) (*domain.ProductWithCategory, error)
	FindAll(ctx context.Context) ([]domain.ProductWithCategory, error)
}

// ProductWriter defines write operations for products
type ProductWriter interface {
	Create(ctx context.Context, p domain.Product) (*domain.Product, error)
	Update(ctx context.Context, id int, p domain.Product) (*domain.Product, error)
	Delete(ctx context.Context, id int) error
}

// CategoryReader defines read operations for categories
type CategoryReader interface {
	FindByID(ctx context.Context, id int) (*domain.Category, error)
	FindAll(ctx context.Context) ([]domain.Category, error)
}

// CategoryWriter defines write operations for categories
type CategoryWriter interface {
	Create(ctx context.Context, c domain.Category) (*domain.Category, error)
	Update(ctx context.Context, id int, c domain.Category) (*domain.Category, error)
	Delete(ctx context.Context, id int) error
}
