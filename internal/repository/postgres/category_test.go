package postgres

import (
	"context"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"kasir-api/internal/model"
)

func TestCategoryRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCategoryRepository(db)
	ctx := context.Background()

	category := domain.Category{
		Name:        "Test Category",
		Description: "Test Description",
	}

	created, err := repo.Create(ctx, category)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if created.ID == 0 {
		t.Error("Created category should have ID")
	}
	if created.Name != category.Name {
		t.Errorf("Name = %v, want %v", created.Name, category.Name)
	}
}

func TestCategoryRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCategoryRepository(db)
	ctx := context.Background()

	category := domain.Category{Name: "Test Category", Description: "Test Description"}
	created, _ := repo.Create(ctx, category)

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

func TestCategoryRepository_FindByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCategoryRepository(db)
	ctx := context.Background()

	_, err := repo.FindByID(ctx, 99999)
	if err != domain.ErrNotFound {
		t.Errorf("FindByID() error = %v, want %v", err, domain.ErrNotFound)
	}
}

func TestCategoryRepository_FindAll(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCategoryRepository(db)
	ctx := context.Background()

	repo.Create(ctx, domain.Category{Name: "Category 1", Description: "Desc 1"})
	repo.Create(ctx, domain.Category{Name: "Category 2", Description: "Desc 2"})

	categories, err := repo.FindAll(ctx)
	if err != nil {
		t.Fatalf("FindAll() error = %v", err)
	}

	if len(categories) < 2 {
		t.Errorf("FindAll() returned %d categories, want at least 2", len(categories))
	}
}

func TestCategoryRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCategoryRepository(db)
	ctx := context.Background()

	category := domain.Category{Name: "Original", Description: "Original Desc"}
	created, _ := repo.Create(ctx, category)

	updated := domain.Category{Name: "Updated", Description: "Updated Desc"}
	result, err := repo.Update(ctx, created.ID, updated)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if result.Name != updated.Name {
		t.Errorf("Name = %v, want %v", result.Name, updated.Name)
	}
	if result.Description != updated.Description {
		t.Errorf("Description = %v, want %v", result.Description, updated.Description)
	}
}

func TestCategoryRepository_Update_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCategoryRepository(db)
	ctx := context.Background()

	_, err := repo.Update(ctx, 99999, domain.Category{Name: "Test", Description: "Test"})
	if err != domain.ErrNotFound {
		t.Errorf("Update() error = %v, want %v", err, domain.ErrNotFound)
	}
}

func TestCategoryRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCategoryRepository(db)
	ctx := context.Background()

	category := domain.Category{Name: "To Delete", Description: "Will be deleted"}
	created, _ := repo.Create(ctx, category)

	err := repo.Delete(ctx, created.ID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err = repo.FindByID(ctx, created.ID)
	if err != domain.ErrNotFound {
		t.Errorf("After delete, FindByID() error = %v, want %v", err, domain.ErrNotFound)
	}
}

func TestCategoryRepository_Delete_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCategoryRepository(db)
	ctx := context.Background()

	err := repo.Delete(ctx, 99999)
	if err != domain.ErrNotFound {
		t.Errorf("Delete() error = %v, want %v", err, domain.ErrNotFound)
	}
}
