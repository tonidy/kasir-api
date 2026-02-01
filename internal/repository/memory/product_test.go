package memory

import (
	"context"
	"testing"

	"kasir-api/internal/model"
)

func TestProductRepository_Create(t *testing.T) {
	repo := NewProductRepository()
	ctx := context.Background()

	product := domain.Product{
		Name:  "Indomie",
		Price: 3500,
		Stock: 100,
	}

	created, err := repo.Create(ctx, product)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if created.ID == 0 {
		t.Error("Created product should have ID")
	}
	if created.Name != product.Name {
		t.Errorf("Name = %v, want %v", created.Name, product.Name)
	}
}

func TestProductRepository_FindByID(t *testing.T) {
	repo := NewProductRepository()
	ctx := context.Background()

	product := domain.Product{Name: "Indomie", Price: 3500, Stock: 100}
	created, _ := repo.Create(ctx, product)

	found, err := repo.FindByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if found.ID != created.ID {
		t.Errorf("ID = %v, want %v", found.ID, created.ID)
	}
}

func TestProductRepository_FindByID_NotFound(t *testing.T) {
	repo := NewProductRepository()
	ctx := context.Background()

	_, err := repo.FindByID(ctx, 999)
	if err != domain.ErrNotFound {
		t.Errorf("FindByID() error = %v, want %v", err, domain.ErrNotFound)
	}
}

func TestProductRepository_FindAll(t *testing.T) {
	repo := NewProductRepository()
	ctx := context.Background()

	repo.Create(ctx, domain.Product{Name: "Product 1", Price: 1000, Stock: 10})
	repo.Create(ctx, domain.Product{Name: "Product 2", Price: 2000, Stock: 20})

	products, err := repo.FindAll(ctx)
	if err != nil {
		t.Fatalf("FindAll() error = %v", err)
	}

	if len(products) != 2 {
		t.Errorf("FindAll() returned %d products, want 2", len(products))
	}
}

func TestProductRepository_Update(t *testing.T) {
	repo := NewProductRepository()
	ctx := context.Background()

	product := domain.Product{Name: "Indomie", Price: 3500, Stock: 100}
	created, _ := repo.Create(ctx, product)

	updated := domain.Product{Name: "Indomie Goreng", Price: 4000, Stock: 150}
	result, err := repo.Update(ctx, created.ID, updated)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if result.Name != updated.Name {
		t.Errorf("Name = %v, want %v", result.Name, updated.Name)
	}
	if result.Price != updated.Price {
		t.Errorf("Price = %v, want %v", result.Price, updated.Price)
	}
}

func TestProductRepository_Update_NotFound(t *testing.T) {
	repo := NewProductRepository()
	ctx := context.Background()

	_, err := repo.Update(ctx, 999, domain.Product{Name: "Test", Price: 1000, Stock: 10})
	if err != domain.ErrNotFound {
		t.Errorf("Update() error = %v, want %v", err, domain.ErrNotFound)
	}
}

func TestProductRepository_Delete(t *testing.T) {
	repo := NewProductRepository()
	ctx := context.Background()

	product := domain.Product{Name: "Indomie", Price: 3500, Stock: 100}
	created, _ := repo.Create(ctx, product)

	err := repo.Delete(ctx, created.ID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err = repo.FindByID(ctx, created.ID)
	if err != domain.ErrNotFound {
		t.Errorf("After delete, FindByID() error = %v, want %v", err, domain.ErrNotFound)
	}
}

func TestProductRepository_Delete_NotFound(t *testing.T) {
	repo := NewProductRepository()
	ctx := context.Background()

	err := repo.Delete(ctx, 999)
	if err != domain.ErrNotFound {
		t.Errorf("Delete() error = %v, want %v", err, domain.ErrNotFound)
	}
}
