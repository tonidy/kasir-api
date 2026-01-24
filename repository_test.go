package main

import "testing"

// Unit Tests - Repository
func TestProductRepository_FindByID(t *testing.T) {
	repo := &ProductRepository{
		data: []Product{
			{ID: 1, Name: "Test Product", Price: 1000, Stock: 5},
		},
	}

	product, found := repo.FindByID(1)
	if !found || product.Name != "Test Product" {
		t.Errorf("Expected to find product, got found=%v", found)
	}

	_, found = repo.FindByID(999)
	if found {
		t.Error("Expected not to find product with ID 999")
	}
}

func TestProductRepository_FindAll(t *testing.T) {
	repo := &ProductRepository{
		data: []Product{
			{ID: 1, Name: "Product 1", Price: 1000, Stock: 5},
			{ID: 2, Name: "Product 2", Price: 2000, Stock: 10},
		},
	}

	products := repo.FindAll()
	if len(products) != 2 {
		t.Errorf("Expected 2 products, got %d", len(products))
	}
}

func TestProductRepository_Create(t *testing.T) {
	repo := &ProductRepository{data: []Product{}}

	newProduct := Product{Name: "New Product", Price: 2000, Stock: 10}
	created := repo.Create(newProduct)

	if created.ID != 1 || created.Name != "New Product" {
		t.Errorf("Expected ID=1, got %d", created.ID)
	}
}

func TestProductRepository_Update(t *testing.T) {
	repo := &ProductRepository{
		data: []Product{
			{ID: 1, Name: "Old Name", Price: 1000, Stock: 5},
		},
	}

	updated, found := repo.Update(1, Product{Name: "New Name", Price: 2000, Stock: 10})
	if !found || updated.Name != "New Name" {
		t.Error("Expected to update product")
	}

	_, found = repo.Update(999, Product{})
	if found {
		t.Error("Expected not to find product with ID 999")
	}
}

func TestProductRepository_Delete(t *testing.T) {
	repo := &ProductRepository{
		data: []Product{
			{ID: 1, Name: "Test", Price: 1000, Stock: 5},
		},
	}

	if !repo.Delete(1) {
		t.Error("Expected to delete product")
	}

	if repo.Delete(999) {
		t.Error("Expected not to delete non-existent product")
	}
}

func TestCategoryRepository_FindByID(t *testing.T) {
	repo := &InMemoryCategoryRepository{
		data: []Category{
			{ID: 1, Name: "Test Category", Description: "Test Desc"},
		},
	}

	category, found := repo.FindByID(1)
	if !found || category.Name != "Test Category" {
		t.Errorf("Expected to find category, got found=%v", found)
	}

	_, found = repo.FindByID(999)
	if found {
		t.Error("Expected not to find category with ID 999")
	}
}

func TestCategoryRepository_FindAll(t *testing.T) {
	repo := &InMemoryCategoryRepository{
		data: []Category{
			{ID: 1, Name: "Category 1", Description: "Desc 1"},
			{ID: 2, Name: "Category 2", Description: "Desc 2"},
		},
	}

	categories := repo.FindAll()
	if len(categories) != 2 {
		t.Errorf("Expected 2 categories, got %d", len(categories))
	}
}

func TestCategoryRepository_Create(t *testing.T) {
	repo := &InMemoryCategoryRepository{data: []Category{}}

	newCategory := Category{Name: "New Category", Description: "New Desc"}
	created := repo.Create(newCategory)

	if created.ID != 1 || created.Name != "New Category" {
		t.Errorf("Expected ID=1, got %d", created.ID)
	}
}

func TestCategoryRepository_Update(t *testing.T) {
	repo := &InMemoryCategoryRepository{
		data: []Category{
			{ID: 1, Name: "Old Name", Description: "Old Desc"},
		},
	}

	updated, found := repo.Update(1, Category{Name: "New Name", Description: "New Desc"})
	if !found || updated.Name != "New Name" {
		t.Error("Expected to update category")
	}

	_, found = repo.Update(999, Category{})
	if found {
		t.Error("Expected not to find category with ID 999")
	}
}

func TestCategoryRepository_Delete(t *testing.T) {
	repo := &InMemoryCategoryRepository{
		data: []Category{
			{ID: 1, Name: "Test", Description: "Test"},
		},
	}

	if !repo.Delete(1) {
		t.Error("Expected to delete category")
	}

	if repo.Delete(999) {
		t.Error("Expected not to delete non-existent category")
	}
}
