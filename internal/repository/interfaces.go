package repository

import (
	"context"
	"kasir-api/internal/model"
)

// ProductReader defines read operations for products
type ProductReader interface {
	FindByID(ctx context.Context, id int) (*model.Product, error)
	FindAll(ctx context.Context) ([]model.Product, error)
	FindByFilters(ctx context.Context, name string, active *bool) ([]model.Product, error)
}

// ProductWriter defines write operations for products
type ProductWriter interface {
	Create(ctx context.Context, p model.Product) (*model.Product, error)
	Update(ctx context.Context, id int, p model.Product) (*model.Product, error)
	Delete(ctx context.Context, id int) error
}

// CategoryReader defines read operations for categories
type CategoryReader interface {
	FindByID(ctx context.Context, id int) (*model.Category, error)
	FindAll(ctx context.Context) ([]model.Category, error)
}

// CategoryWriter defines write operations for categories
type CategoryWriter interface {
	Create(ctx context.Context, c model.Category) (*model.Category, error)
	Update(ctx context.Context, id int, c model.Category) (*model.Category, error)
	Delete(ctx context.Context, id int) error
}

// TransactionWriter defines write operations for transactions
type TransactionWriter interface {
	CreateTransaction(ctx context.Context, items []model.CheckoutItem) (*model.Transaction, error)
}

// ReportReader defines read operations for reports
type ReportReader interface {
	GetTodayReport(ctx context.Context) (*model.ReportSummary, error)
	GetReportByDateRange(ctx context.Context, startDate, endDate string) (*model.ReportSummary, error)
}
