package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestRepo() {
	repo = &ProductRepository{
		data: []Product{
			{ID: 1, Name: "Product 1", Price: 1000, Stock: 5},
			{ID: 2, Name: "Product 2", Price: 2000, Stock: 10},
		},
	}
}

func TestGetProductByID_Success(t *testing.T) {
	setupTestRepo()

	req := httptest.NewRequest(http.MethodGet, "/api/products/1", nil)
	w := httptest.NewRecorder()

	getProductByID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var product Product
	json.NewDecoder(w.Body).Decode(&product)
	if product.ID != 1 {
		t.Errorf("Expected product ID 1, got %d", product.ID)
	}
}

func TestGetProductByID_NotFound(t *testing.T) {
	setupTestRepo()

	req := httptest.NewRequest(http.MethodGet, "/api/products/999", nil)
	w := httptest.NewRecorder()

	getProductByID(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestGetProductByID_InvalidID(t *testing.T) {
	setupTestRepo()

	req := httptest.NewRequest(http.MethodGet, "/api/products/abc", nil)
	w := httptest.NewRecorder()

	getProductByID(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestUpdateProduct_Success(t *testing.T) {
	setupTestRepo()

	updatedProduct := Product{Name: "Updated Name", Price: 2000, Stock: 10}
	body, _ := json.Marshal(updatedProduct)

	req := httptest.NewRequest(http.MethodPut, "/api/products/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	updateProduct(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestUpdateProduct_NotFound(t *testing.T) {
	setupTestRepo()

	updatedProduct := Product{Name: "Updated", Price: 2000, Stock: 10}
	body, _ := json.Marshal(updatedProduct)

	req := httptest.NewRequest(http.MethodPut, "/api/products/999", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	updateProduct(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestUpdateProduct_InvalidID(t *testing.T) {
	setupTestRepo()

	body := bytes.NewBufferString(`{"name":"Test"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/products/abc", body)
	w := httptest.NewRecorder()

	updateProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestUpdateProduct_InvalidJSON(t *testing.T) {
	setupTestRepo()

	req := httptest.NewRequest(http.MethodPut, "/api/products/1", bytes.NewBufferString("invalid"))
	w := httptest.NewRecorder()

	updateProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestDeleteProduct_Success(t *testing.T) {
	setupTestRepo()

	req := httptest.NewRequest(http.MethodDelete, "/api/products/1", nil)
	w := httptest.NewRecorder()

	deleteProduct(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestDeleteProduct_NotFound(t *testing.T) {
	setupTestRepo()

	req := httptest.NewRequest(http.MethodDelete, "/api/products/999", nil)
	w := httptest.NewRecorder()

	deleteProduct(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestDeleteProduct_InvalidID(t *testing.T) {
	setupTestRepo()

	req := httptest.NewRequest(http.MethodDelete, "/api/products/abc", nil)
	w := httptest.NewRecorder()

	deleteProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestGetCategoryByID_Success(t *testing.T) {
	categoryRepo = &InMemoryCategoryRepository{
		data: []Category{
			{ID: 1, Name: "Category 1", Description: "Desc 1"},
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/api/categories/1", nil)
	w := httptest.NewRecorder()

	getCategoryByID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestGetCategoryByID_NotFound(t *testing.T) {
	categoryRepo = &InMemoryCategoryRepository{data: []Category{}}

	req := httptest.NewRequest(http.MethodGet, "/api/categories/999", nil)
	w := httptest.NewRecorder()

	getCategoryByID(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestGetCategoryByID_InvalidID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/categories/abc", nil)
	w := httptest.NewRecorder()

	getCategoryByID(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestUpdateCategory_Success(t *testing.T) {
	categoryRepo = &InMemoryCategoryRepository{
		data: []Category{
			{ID: 1, Name: "Old", Description: "Old Desc"},
		},
	}

	updated := Category{Name: "Updated", Description: "Updated Desc"}
	body, _ := json.Marshal(updated)

	req := httptest.NewRequest(http.MethodPut, "/api/categories/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	updateCategory(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestUpdateCategory_NotFound(t *testing.T) {
	categoryRepo = &InMemoryCategoryRepository{data: []Category{}}

	updated := Category{Name: "Updated", Description: "Desc"}
	body, _ := json.Marshal(updated)

	req := httptest.NewRequest(http.MethodPut, "/api/categories/999", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	updateCategory(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestUpdateCategory_InvalidID(t *testing.T) {
	body := bytes.NewBufferString(`{"name":"Test"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/categories/abc", body)
	w := httptest.NewRecorder()

	updateCategory(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestUpdateCategory_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/api/categories/1", bytes.NewBufferString("invalid"))
	w := httptest.NewRecorder()

	updateCategory(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestDeleteCategory_Success(t *testing.T) {
	categoryRepo = &InMemoryCategoryRepository{
		data: []Category{
			{ID: 1, Name: "Category", Description: "Desc"},
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/categories/1", nil)
	w := httptest.NewRecorder()

	deleteCategory(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestDeleteCategory_NotFound(t *testing.T) {
	categoryRepo = &InMemoryCategoryRepository{data: []Category{}}

	req := httptest.NewRequest(http.MethodDelete, "/api/categories/999", nil)
	w := httptest.NewRecorder()

	deleteCategory(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestDeleteCategory_InvalidID(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/categories/abc", nil)
	w := httptest.NewRecorder()

	deleteCategory(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}
