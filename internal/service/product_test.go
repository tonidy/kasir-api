package service

import (
	"context"
	"testing"

	"kasir-api/internal/model"
	"kasir-api/internal/repository/memory"
)

func TestProductService_Create(t *testing.T) {
	repo := memory.NewProductRepository()
	svc := NewProductService(repo, repo)
	ctx := context.Background()

	product := domain.Product{
		Name:  "Indomie",
		Price: 3500,
		Stock: 100,
	}

	created, err := svc.Create(ctx, product)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if created.ID == 0 {
		t.Error("Created product should have ID")
	}
}

func TestProductService_Create_ValidationError(t *testing.T) {
	repo := memory.NewProductRepository()
	svc := NewProductService(repo, repo)
	ctx := context.Background()

	product := domain.Product{
		Name:  "", // Invalid: empty name
		Price: 3500,
		Stock: 100,
	}

	_, err := svc.Create(ctx, product)
	if err == nil {
		t.Error("Create() should return validation error for empty name")
	}
}

func TestProductService_GetByID(t *testing.T) {
	repo := memory.NewProductRepository()
	svc := NewProductService(repo, repo)
	ctx := context.Background()

	product := domain.Product{Name: "Indomie", Price: 3500, Stock: 100}
	created, _ := svc.Create(ctx, product)

	found, err := svc.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	if found.ID != created.ID {
		t.Errorf("ID = %v, want %v", found.ID, created.ID)
	}
}

func TestProductService_GetByID_NotFound(t *testing.T) {
	repo := memory.NewProductRepository()
	svc := NewProductService(repo, repo)
	ctx := context.Background()

	_, err := svc.GetByID(ctx, 999)
	if err != domain.ErrNotFound {
		t.Errorf("GetByID() error = %v, want %v", err, domain.ErrNotFound)
	}
}

func TestProductService_GetAll(t *testing.T) {
	repo := memory.NewProductRepository()
	svc := NewProductService(repo, repo)
	ctx := context.Background()

	svc.Create(ctx, domain.Product{Name: "Product 1", Price: 1000, Stock: 10})
	svc.Create(ctx, domain.Product{Name: "Product 2", Price: 2000, Stock: 20})

	products, err := svc.GetAll(ctx)
	if err != nil {
		t.Fatalf("GetAll() error = %v", err)
	}

	if len(products) != 2 {
		t.Errorf("GetAll() returned %d products, want 2", len(products))
	}
}

func TestProductService_Update(t *testing.T) {
	repo := memory.NewProductRepository()
	svc := NewProductService(repo, repo)
	ctx := context.Background()

	product := domain.Product{Name: "Indomie", Price: 3500, Stock: 100}
	created, _ := svc.Create(ctx, product)

	updated := domain.Product{Name: "Indomie Goreng", Price: 4000, Stock: 150}
	result, err := svc.Update(ctx, created.ID, updated)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if result.Name != updated.Name {
		t.Errorf("Name = %v, want %v", result.Name, updated.Name)
	}
}

func TestProductService_Update_ValidationError(t *testing.T) {
	repo := memory.NewProductRepository()
	svc := NewProductService(repo, repo)
	ctx := context.Background()

	product := domain.Product{Name: "Indomie", Price: 3500, Stock: 100}
	created, _ := svc.Create(ctx, product)

	updated := domain.Product{Name: "", Price: 4000, Stock: 150} // Invalid
	_, err := svc.Update(ctx, created.ID, updated)
	if err == nil {
		t.Error("Update() should return validation error for empty name")
	}
}

func TestProductService_Delete(t *testing.T) {
	repo := memory.NewProductRepository()
	svc := NewProductService(repo, repo)
	ctx := context.Background()

	product := domain.Product{Name: "Indomie", Price: 3500, Stock: 100}
	created, _ := svc.Create(ctx, product)

	err := svc.Delete(ctx, created.ID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err = svc.GetByID(ctx, created.ID)
	if err != domain.ErrNotFound {
		t.Errorf("After delete, GetByID() error = %v, want %v", err, domain.ErrNotFound)
	}
}
