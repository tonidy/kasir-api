package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetupRoutes(t *testing.T) {
	repo = &ProductRepository{
		data: []Product{
			{ID: 1, Name: "Product 1", Price: 1000, Stock: 5},
			{ID: 2, Name: "Product 2", Price: 2000, Stock: 10},
		},
	}
	http.DefaultServeMux = new(http.ServeMux)
	setupRoutes()

	tests := []struct {
		method string
		path   string
		status int
	}{
		{http.MethodGet, "/api/products", http.StatusOK},
		{http.MethodGet, "/api/products/1", http.StatusOK},
		{http.MethodGet, "/health", http.StatusOK},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.path, nil)
		w := httptest.NewRecorder()
		http.DefaultServeMux.ServeHTTP(w, req)
		if w.Code != tt.status {
			t.Errorf("%s %s: expected %d, got %d", tt.method, tt.path, tt.status, w.Code)
		}
	}
}

func TestProductsRoute_POST_Success(t *testing.T) {
	repo = &ProductRepository{data: []Product{}}
	http.DefaultServeMux = new(http.ServeMux)
	setupRoutes()

	newProduct := Product{Name: "New Product", Price: 3000, Stock: 15}
	body, _ := json.Marshal(newProduct)

	req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	http.DefaultServeMux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var created Product
	json.NewDecoder(w.Body).Decode(&created)
	if created.Name != "New Product" {
		t.Errorf("Expected name 'New Product', got %s", created.Name)
	}
}

func TestProductsRoute_POST_InvalidJSON(t *testing.T) {
	repo = &ProductRepository{data: []Product{}}
	http.DefaultServeMux = new(http.ServeMux)
	setupRoutes()

	req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBufferString("invalid"))
	w := httptest.NewRecorder()

	http.DefaultServeMux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestProductsWithIDRoute_PUT(t *testing.T) {
	repo = &ProductRepository{
		data: []Product{
			{ID: 1, Name: "Product 1", Price: 1000, Stock: 5},
		},
	}
	http.DefaultServeMux = new(http.ServeMux)
	setupRoutes()

	updatedProduct := Product{Name: "Updated", Price: 5000, Stock: 20}
	body, _ := json.Marshal(updatedProduct)

	req := httptest.NewRequest(http.MethodPut, "/api/products/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	http.DefaultServeMux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestProductsWithIDRoute_DELETE(t *testing.T) {
	repo = &ProductRepository{
		data: []Product{
			{ID: 1, Name: "Product 1", Price: 1000, Stock: 5},
		},
	}
	http.DefaultServeMux = new(http.ServeMux)
	setupRoutes()

	req := httptest.NewRequest(http.MethodDelete, "/api/products/1", nil)
	w := httptest.NewRecorder()

	http.DefaultServeMux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestCategoriesRoute_GET(t *testing.T) {
	categoryRepo = &InMemoryCategoryRepository{
		data: []Category{
			{ID: 1, Name: "Category 1", Description: "Desc 1"},
		},
	}
	http.DefaultServeMux = new(http.ServeMux)
	setupRoutes()

	req := httptest.NewRequest(http.MethodGet, "/api/categories", nil)
	w := httptest.NewRecorder()

	http.DefaultServeMux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestCategoriesRoute_POST_Success(t *testing.T) {
	categoryRepo = &InMemoryCategoryRepository{data: []Category{}}
	http.DefaultServeMux = new(http.ServeMux)
	setupRoutes()

	newCategory := Category{Name: "New Category", Description: "New Desc"}
	body, _ := json.Marshal(newCategory)

	req := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	http.DefaultServeMux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
}

func TestCategoriesWithIDRoute_PUT(t *testing.T) {
	categoryRepo = &InMemoryCategoryRepository{
		data: []Category{
			{ID: 1, Name: "Category 1", Description: "Desc 1"},
		},
	}
	http.DefaultServeMux = new(http.ServeMux)
	setupRoutes()

	updatedCategory := Category{Name: "Updated", Description: "Updated Desc"}
	body, _ := json.Marshal(updatedCategory)

	req := httptest.NewRequest(http.MethodPut, "/api/categories/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	http.DefaultServeMux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestCategoriesWithIDRoute_DELETE(t *testing.T) {
	categoryRepo = &InMemoryCategoryRepository{
		data: []Category{
			{ID: 1, Name: "Category 1", Description: "Desc 1"},
		},
	}
	http.DefaultServeMux = new(http.ServeMux)
	setupRoutes()

	req := httptest.NewRequest(http.MethodDelete, "/api/categories/1", nil)
	w := httptest.NewRecorder()

	http.DefaultServeMux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
