package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"kasir-api/internal/model"
)

// Mock service for testing
type mockProductService struct {
	getByIDFunc func(ctx context.Context, id int) (*domain.ProductWithCategory, error)
	getAllFunc  func(ctx context.Context) ([]domain.ProductWithCategory, error)
	createFunc  func(ctx context.Context, p domain.Product) (*domain.Product, error)
	updateFunc  func(ctx context.Context, id int, p domain.Product) (*domain.Product, error)
	deleteFunc  func(ctx context.Context, id int) error
}

func (m *mockProductService) GetByID(ctx context.Context, id int) (*domain.ProductWithCategory, error) {
	return m.getByIDFunc(ctx, id)
}

func (m *mockProductService) GetAll(ctx context.Context) ([]domain.ProductWithCategory, error) {
	return m.getAllFunc(ctx)
}

func (m *mockProductService) Create(ctx context.Context, p domain.Product) (*domain.Product, error) {
	return m.createFunc(ctx, p)
}

func (m *mockProductService) Update(ctx context.Context, id int, p domain.Product) (*domain.Product, error) {
	return m.updateFunc(ctx, id, p)
}

func (m *mockProductService) Delete(ctx context.Context, id int) error {
	return m.deleteFunc(ctx, id)
}

func TestProductHandler_GetAll(t *testing.T) {
	mockSvc := &mockProductService{
		getAllFunc: func(ctx context.Context) ([]domain.ProductWithCategory, error) {
			return []domain.ProductWithCategory{
				{Product: domain.Product{ID: 1, Name: "Product 1", Price: 1000, Stock: 10}},
				{Product: domain.Product{ID: 2, Name: "Product 2", Price: 2000, Stock: 20}},
			}, nil
		},
	}

	handler := NewProductHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/products", nil)
	w := httptest.NewRecorder()

	handler.GetAll(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var products []domain.ProductWithCategory
	json.NewDecoder(w.Body).Decode(&products)
	if len(products) != 2 {
		t.Errorf("Expected 2 products, got %d", len(products))
	}
}

func TestProductHandler_GetByID(t *testing.T) {
	mockSvc := &mockProductService{
		getByIDFunc: func(ctx context.Context, id int) (*domain.ProductWithCategory, error) {
			return &domain.ProductWithCategory{
				Product: domain.Product{ID: 1, Name: "Product 1", Price: 1000, Stock: 10},
			}, nil
		},
	}

	handler := NewProductHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/products/1", nil)
	w := httptest.NewRecorder()

	handler.GetByID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestProductHandler_GetByID_NotFound(t *testing.T) {
	mockSvc := &mockProductService{
		getByIDFunc: func(ctx context.Context, id int) (*domain.ProductWithCategory, error) {
			return nil, domain.ErrNotFound
		},
	}

	handler := NewProductHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/products/999", nil)
	w := httptest.NewRecorder()

	handler.GetByID(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestProductHandler_Create(t *testing.T) {
	mockSvc := &mockProductService{
		createFunc: func(ctx context.Context, p domain.Product) (*domain.Product, error) {
			p.ID = 1
			return &p, nil
		},
	}

	handler := NewProductHandler(mockSvc)
	body := bytes.NewBufferString(`{"name":"New Product","price":5000,"stock":50}`)
	req := httptest.NewRequest(http.MethodPost, "/api/products", body)
	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
}

func TestProductHandler_Create_InvalidJSON(t *testing.T) {
	mockSvc := &mockProductService{}
	handler := NewProductHandler(mockSvc)
	body := bytes.NewBufferString(`{invalid json}`)
	req := httptest.NewRequest(http.MethodPost, "/api/products", body)
	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestProductHandler_Update(t *testing.T) {
	mockSvc := &mockProductService{
		updateFunc: func(ctx context.Context, id int, p domain.Product) (*domain.Product, error) {
			p.ID = id
			return &p, nil
		},
	}

	handler := NewProductHandler(mockSvc)
	body := bytes.NewBufferString(`{"name":"Updated Product","price":6000,"stock":60}`)
	req := httptest.NewRequest(http.MethodPut, "/api/products/1", body)
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestProductHandler_Delete(t *testing.T) {
	mockSvc := &mockProductService{
		deleteFunc: func(ctx context.Context, id int) error {
			return nil
		},
	}

	handler := NewProductHandler(mockSvc)
	req := httptest.NewRequest(http.MethodDelete, "/api/products/1", nil)
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
