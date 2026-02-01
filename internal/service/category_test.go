package service

import (
	"context"
	"testing"

	"kasir-api/internal/model"
	"kasir-api/internal/repository/memory"
)

func TestCategoryService_Create(t *testing.T) {
	repo := memory.NewCategoryRepository()
	svc := NewCategoryService(repo, repo)
	ctx := context.Background()

	category := domain.Category{
		Name:        "Food",
		Description: "Food items",
	}

	created, err := svc.Create(ctx, category)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if created.ID == 0 {
		t.Error("Created category should have ID")
	}
}

func TestCategoryService_Create_ValidationError(t *testing.T) {
	repo := memory.NewCategoryRepository()
	svc := NewCategoryService(repo, repo)
	ctx := context.Background()

	category := domain.Category{
		Name:        "", // Invalid: empty name
		Description: "Food items",
	}

	_, err := svc.Create(ctx, category)
	if err == nil {
		t.Error("Create() should return validation error for empty name")
	}
}

func TestCategoryService_GetByID(t *testing.T) {
	repo := memory.NewCategoryRepository()
	svc := NewCategoryService(repo, repo)
	ctx := context.Background()

	category := domain.Category{Name: "Food", Description: "Food items"}
	created, _ := svc.Create(ctx, category)

	found, err := svc.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	if found.ID != created.ID {
		t.Errorf("ID = %v, want %v", found.ID, created.ID)
	}
}

func TestCategoryService_GetAll(t *testing.T) {
	repo := memory.NewCategoryRepository()
	svc := NewCategoryService(repo, repo)
	ctx := context.Background()

	svc.Create(ctx, domain.Category{Name: "Food", Description: "Food items"})
	svc.Create(ctx, domain.Category{Name: "Beverage", Description: "Drinks"})

	categories, err := svc.GetAll(ctx)
	if err != nil {
		t.Fatalf("GetAll() error = %v", err)
	}

	if len(categories) != 2 {
		t.Errorf("GetAll() returned %d categories, want 2", len(categories))
	}
}

func TestCategoryService_Update(t *testing.T) {
	repo := memory.NewCategoryRepository()
	svc := NewCategoryService(repo, repo)
	ctx := context.Background()

	category := domain.Category{Name: "Food", Description: "Food items"}
	created, _ := svc.Create(ctx, category)

	updated := domain.Category{Name: "Food & Beverage", Description: "Food and drinks"}
	result, err := svc.Update(ctx, created.ID, updated)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if result.Name != updated.Name {
		t.Errorf("Name = %v, want %v", result.Name, updated.Name)
	}
}

func TestCategoryService_Delete(t *testing.T) {
	repo := memory.NewCategoryRepository()
	svc := NewCategoryService(repo, repo)
	ctx := context.Background()

	category := domain.Category{Name: "Food", Description: "Food items"}
	created, _ := svc.Create(ctx, category)

	err := svc.Delete(ctx, created.ID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err = svc.GetByID(ctx, created.ID)
	if err != domain.ErrNotFound {
		t.Errorf("After delete, GetByID() error = %v, want %v", err, domain.ErrNotFound)
	}
}
