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
	GetByID(ctx context.Context, id int) (*model.Product, error)
	GetAll(ctx context.Context) ([]model.Product, error)
	GetByFilters(ctx context.Context, name string, active *bool) ([]model.Product, error)
	Create(ctx context.Context, p model.Product) (*model.Product, error)
	Update(ctx context.Context, id int, p model.Product) (*model.Product, error)
	Delete(ctx context.Context, id int) error
}

type ProductHandler struct {
	svc ProductService
}

func NewProductHandler(svc ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	activeStr := r.URL.Query().Get("active")

	var products []model.Product
	var err error

	if name != "" || activeStr != "" {
		var active *bool
		if activeStr != "" {
			val := activeStr == "true"
			active = &val
		}
		products, err = h.svc.GetByFilters(r.Context(), name, active)
	} else {
		products, err = h.svc.GetAll(r.Context())
	}

	if err != nil {
		httputil.WriteError(w, httputil.ErrorStatus(err), err.Error())
		return
	}

	response := make([]dto.ProductResponse, len(products))
	for i, p := range products {
		var catResp *dto.CategoryResponse
		if p.Category != nil {
			catResp = &dto.CategoryResponse{
				ID:          p.Category.ID,
				Name:        p.Category.Name,
				Description: p.Category.Description,
			}
		}
		response[i] = dto.ProductResponse{
			ID:       p.ID,
			Name:     p.Name,
			Price:    p.Price,
			Stock:    p.Stock,
			Active:   p.Active,
			Category: catResp,
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

	var categoryResp *dto.CategoryResponse
	if product.Category != nil {
		categoryResp = &dto.CategoryResponse{
			ID:          product.Category.ID,
			Name:        product.Category.Name,
			Description: product.Category.Description,
		}
	}

	response := dto.ProductResponse{
		ID:       product.ID,
		Name:     product.Name,
		Price:    product.Price,
		Stock:    product.Stock,
		Active:   product.Active,
		Category: categoryResp,
	}

	httputil.WriteJSON(w, http.StatusOK, response)
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var product model.Product
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

	var product model.Product
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
