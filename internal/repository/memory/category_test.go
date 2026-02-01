package memory

import (
	"context"
	"testing"

	"kasir-api/internal/model"
)

func TestCategoryRepository_Create(t *testing.T) {
	repo := NewCategoryRepository()
	ctx := context.Background()

	category := model.Category{
		Name:        "Food",
		Description: "Food items",
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
	repo := NewCategoryRepository()
	ctx := context.Background()

	category := model.Category{Name: "Food", Description: "Food items"}
	created, _ := repo.Create(ctx, category)

	found, err := repo.FindByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if found.ID != created.ID {
		t.Errorf("ID = %v, want %v", found.ID, created.ID)
	}
}

func TestCategoryRepository_FindByID_NotFound(t *testing.T) {
	repo := NewCategoryRepository()
	ctx := context.Background()

	_, err := repo.FindByID(ctx, 999)
	if err != model.ErrNotFound {
		t.Errorf("FindByID() error = %v, want %v", err, model.ErrNotFound)
	}
}

func TestCategoryRepository_FindAll(t *testing.T) {
	repo := NewCategoryRepository()
	ctx := context.Background()

	repo.Create(ctx, model.Category{Name: "Food", Description: "Food items"})
	repo.Create(ctx, model.Category{Name: "Beverage", Description: "Drinks"})

	categories, err := repo.FindAll(ctx)
	if err != nil {
		t.Fatalf("FindAll() error = %v", err)
	}

	if len(categories) != 2 {
		t.Errorf("FindAll() returned %d categories, want 2", len(categories))
	}
}

func TestCategoryRepository_Update(t *testing.T) {
	repo := NewCategoryRepository()
	ctx := context.Background()

	category := model.Category{Name: "Food", Description: "Food items"}
	created, _ := repo.Create(ctx, category)

	updated := model.Category{Name: "Food & Beverage", Description: "Food and drinks"}
	result, err := repo.Update(ctx, created.ID, updated)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if result.Name != updated.Name {
		t.Errorf("Name = %v, want %v", result.Name, updated.Name)
	}
}

func TestCategoryRepository_Update_NotFound(t *testing.T) {
	repo := NewCategoryRepository()
	ctx := context.Background()

	_, err := repo.Update(ctx, 999, model.Category{Name: "Test", Description: "Test"})
	if err != model.ErrNotFound {
		t.Errorf("Update() error = %v, want %v", err, model.ErrNotFound)
	}
}

func TestCategoryRepository_Delete(t *testing.T) {
	repo := NewCategoryRepository()
	ctx := context.Background()

	category := model.Category{Name: "Food", Description: "Food items"}
	created, _ := repo.Create(ctx, category)

	err := repo.Delete(ctx, created.ID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err = repo.FindByID(ctx, created.ID)
	if err != model.ErrNotFound {
		t.Errorf("After delete, FindByID() error = %v, want %v", err, model.ErrNotFound)
	}
}

func TestCategoryRepository_Delete_NotFound(t *testing.T) {
	repo := NewCategoryRepository()
	ctx := context.Background()

	err := repo.Delete(ctx, 999)
	if err != model.ErrNotFound {
		t.Errorf("Delete() error = %v, want %v", err, model.ErrNotFound)
	}
}
