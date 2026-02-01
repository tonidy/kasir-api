package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"kasir-api/internal/dto"
	"kasir-api/internal/model"
	"kasir-api/pkg/httputil"
)

type ProductService interface {
	GetByID(ctx context.Context, id int) (*domain.ProductWithCategory, error)
	GetAll(ctx context.Context) ([]domain.ProductWithCategory, error)
	Create(ctx context.Context, p domain.Product) (*domain.Product, error)
	Update(ctx context.Context, id int, p domain.Product) (*domain.Product, error)
	Delete(ctx context.Context, id int) error
}

type ProductHandler struct {
	svc ProductService
}

func NewProductHandler(svc ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.svc.GetAll(r.Context())
	if err != nil {
		httputil.WriteError(w, httputil.ErrorStatus(err), err.Error())
		return
	}

	response := make([]dto.ProductResponse, len(products))
	for i, p := range products {
		response[i] = dto.ProductResponse{
			ID:           p.ID,
			Name:         p.Name,
			Price:        p.Price,
			Stock:        p.Stock,
			CategoryID:   p.CategoryID,
			CategoryName: p.CategoryName,
		}
	}

	httputil.WriteJSON(w, http.StatusOK, response)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseID(r)
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	product, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, httputil.ErrorStatus(err), err.Error())
		return
	}

	response := dto.ProductResponse{
		ID:           product.ID,
		Name:         product.Name,
		Price:        product.Price,
		Stock:        product.Stock,
		CategoryID:   product.CategoryID,
		CategoryName: product.CategoryName,
	}

	httputil.WriteJSON(w, http.StatusOK, response)
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var product domain.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	created, err := h.svc.Create(r.Context(), product)
	if err != nil {
		httputil.WriteError(w, httputil.ErrorStatus(err), err.Error())
		return
	}
	httputil.WriteJSON(w, http.StatusCreated, created)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseID(r)
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var product domain.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	updated, err := h.svc.Update(r.Context(), id, product)
	if err != nil {
		httputil.WriteError(w, httputil.ErrorStatus(err), err.Error())
		return
	}
	httputil.WriteJSON(w, http.StatusOK, updated)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseID(r)
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		httputil.WriteError(w, httputil.ErrorStatus(err), err.Error())
		return
	}
	httputil.WriteJSON(w, http.StatusOK, map[string]string{"message": "Data successfully deleted"})
}
