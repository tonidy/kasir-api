package postgres

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"kasir-api/internal/model"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	// Skip if no test database configured
	// Use 127.0.0.1 instead of localhost to force IPv4
	dsn := "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=kasir_test sslmode=disable"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Skip("Test database not available:", err)
	}

	if err := db.Ping(); err != nil {
		t.Skip("Test database not available:", err)
	}

	// Clean up tables
	db.Exec("DELETE FROM products")
	db.Exec("DELETE FROM categories")

	return db
}

func TestProductRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewProductRepository(db)
	ctx := context.Background()

	product := model.Product{
		Name:  "Test Product",
		Price: 5000,
		Stock: 50,
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
	db := setupTestDB(t)
	defer db.Close()

	repo := NewProductRepository(db)
	ctx := context.Background()

	product := model.Product{Name: "Test Product", Price: 5000, Stock: 50}
	created, _ := repo.Create(ctx, product)

	found, err := repo.FindByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if found.ID != created.ID {
		t.Errorf("ID = %v, want %v", found.ID, created.ID)
	}
	if found.Name != created.Name {
		t.Errorf("Name = %v, want %v", found.Name, created.Name)
	}
}

func TestProductRepository_FindByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewProductRepository(db)
	ctx := context.Background()

	_, err := repo.FindByID(ctx, 99999)
	if err != model.ErrNotFound {
		t.Errorf("FindByID() error = %v, want %v", err, model.ErrNotFound)
	}
}

func TestProductRepository_FindAll(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewProductRepository(db)
	ctx := context.Background()

	repo.Create(ctx, model.Product{Name: "Product 1", Price: 1000, Stock: 10})
	repo.Create(ctx, model.Product{Name: "Product 2", Price: 2000, Stock: 20})

	products, err := repo.FindAll(ctx)
	if err != nil {
		t.Fatalf("FindAll() error = %v", err)
	}

	if len(products) < 2 {
		t.Errorf("FindAll() returned %d products, want at least 2", len(products))
	}
}

func TestProductRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewProductRepository(db)
	ctx := context.Background()

	product := model.Product{Name: "Original", Price: 3000, Stock: 30}
	created, _ := repo.Create(ctx, product)

	updated := model.Product{Name: "Updated", Price: 4000, Stock: 40}
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
	db := setupTestDB(t)
	defer db.Close()

	repo := NewProductRepository(db)
	ctx := context.Background()

	_, err := repo.Update(ctx, 99999, model.Product{Name: "Test", Price: 1000, Stock: 10})
	if err != model.ErrNotFound {
		t.Errorf("Update() error = %v, want %v", err, model.ErrNotFound)
	}
}

func TestProductRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewProductRepository(db)
	ctx := context.Background()

	product := model.Product{Name: "To Delete", Price: 5000, Stock: 50}
	created, _ := repo.Create(ctx, product)

	err := repo.Delete(ctx, created.ID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err = repo.FindByID(ctx, created.ID)
	if err != model.ErrNotFound {
		t.Errorf("After delete, FindByID() error = %v, want %v", err, model.ErrNotFound)
	}
}

func TestProductRepository_Delete_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewProductRepository(db)
	ctx := context.Background()

	err := repo.Delete(ctx, 99999)
	if err != model.ErrNotFound {
		t.Errorf("Delete() error = %v, want %v", err, model.ErrNotFound)
	}
}
