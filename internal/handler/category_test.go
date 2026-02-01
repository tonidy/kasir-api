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
type mockCategoryService struct {
	getByIDFunc func(ctx context.Context, id int) (*domain.Category, error)
	getAllFunc  func(ctx context.Context) ([]domain.Category, error)
	createFunc  func(ctx context.Context, c domain.Category) (*domain.Category, error)
	updateFunc  func(ctx context.Context, id int, c domain.Category) (*domain.Category, error)
	deleteFunc  func(ctx context.Context, id int) error
}

func (m *mockCategoryService) GetByID(ctx context.Context, id int) (*domain.Category, error) {
	return m.getByIDFunc(ctx, id)
}

func (m *mockCategoryService) GetAll(ctx context.Context) ([]domain.Category, error) {
	return m.getAllFunc(ctx)
}

func (m *mockCategoryService) Create(ctx context.Context, c domain.Category) (*domain.Category, error) {
	return m.createFunc(ctx, c)
}

func (m *mockCategoryService) Update(ctx context.Context, id int, c domain.Category) (*domain.Category, error) {
	return m.updateFunc(ctx, id, c)
}

func (m *mockCategoryService) Delete(ctx context.Context, id int) error {
	return m.deleteFunc(ctx, id)
}

func TestCategoryHandler_GetAll(t *testing.T) {
	mockSvc := &mockCategoryService{
		getAllFunc: func(ctx context.Context) ([]domain.Category, error) {
			return []domain.Category{
				{ID: 1, Name: "Food", Description: "Food items"},
				{ID: 2, Name: "Beverage", Description: "Drinks"},
			}, nil
		},
	}

	handler := NewCategoryHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/categories", nil)
	w := httptest.NewRecorder()

	handler.GetAll(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var categories []domain.Category
	json.NewDecoder(w.Body).Decode(&categories)
	if len(categories) != 2 {
		t.Errorf("Expected 2 categories, got %d", len(categories))
	}
}

func TestCategoryHandler_GetByID(t *testing.T) {
	mockSvc := &mockCategoryService{
		getByIDFunc: func(ctx context.Context, id int) (*domain.Category, error) {
			return &domain.Category{ID: 1, Name: "Food", Description: "Food items"}, nil
		},
	}

	handler := NewCategoryHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/categories/1", nil)
	w := httptest.NewRecorder()

	handler.GetByID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestCategoryHandler_GetByID_NotFound(t *testing.T) {
	mockSvc := &mockCategoryService{
		getByIDFunc: func(ctx context.Context, id int) (*domain.Category, error) {
			return nil, domain.ErrNotFound
		},
	}

	handler := NewCategoryHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/categories/999", nil)
	w := httptest.NewRecorder()

	handler.GetByID(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestCategoryHandler_Create(t *testing.T) {
	mockSvc := &mockCategoryService{
		createFunc: func(ctx context.Context, c domain.Category) (*domain.Category, error) {
			c.ID = 1
			return &c, nil
		},
	}

	handler := NewCategoryHandler(mockSvc)
	body := bytes.NewBufferString(`{"name":"New Category","description":"New Description"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/categories", body)
	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
}

func TestCategoryHandler_Create_InvalidJSON(t *testing.T) {
	mockSvc := &mockCategoryService{}
	handler := NewCategoryHandler(mockSvc)
	body := bytes.NewBufferString(`{invalid json}`)
	req := httptest.NewRequest(http.MethodPost, "/api/categories", body)
	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestCategoryHandler_Update(t *testing.T) {
	mockSvc := &mockCategoryService{
		updateFunc: func(ctx context.Context, id int, c domain.Category) (*domain.Category, error) {
			c.ID = id
			return &c, nil
		},
	}

	handler := NewCategoryHandler(mockSvc)
	body := bytes.NewBufferString(`{"name":"Updated Category","description":"Updated Description"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/categories/1", body)
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestCategoryHandler_Delete(t *testing.T) {
	mockSvc := &mockCategoryService{
		deleteFunc: func(ctx context.Context, id int) error {
			return nil
		},
	}

	handler := NewCategoryHandler(mockSvc)
	req := httptest.NewRequest(http.MethodDelete, "/api/categories/1", nil)
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
