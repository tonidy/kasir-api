package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"kasir-api/internal/model"
	"kasir-api/pkg/httputil"
)

type CategoryService interface {
	GetByID(ctx context.Context, id int) (*model.Category, error)
	GetAll(ctx context.Context) ([]model.Category, error)
	Create(ctx context.Context, c model.Category) (*model.Category, error)
	Update(ctx context.Context, id int, c model.Category) (*model.Category, error)
	Delete(ctx context.Context, id int) error
}

type CategoryHandler struct {
	svc CategoryService
}

func NewCategoryHandler(svc CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := h.svc.GetAll(r.Context())
	if err != nil {
		httputil.HandleError(w, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, categories)
}

func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseID(r)
	if err != nil {
		httputil.HandleError(w, err)
		return
	}

	category, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		httputil.HandleError(w, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, category)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var category model.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	created, err := h.svc.Create(r.Context(), category)
	if err != nil {
		httputil.HandleError(w, err)
		return
	}
	httputil.WriteJSON(w, http.StatusCreated, created)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseID(r)
	if err != nil {
		httputil.HandleError(w, err)
		return
	}

	var category model.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		httputil.HandleError(w, err)
		return
	}

	updated, err := h.svc.Update(r.Context(), id, category)
	if err != nil {
		httputil.HandleError(w, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, updated)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseID(r)
	if err != nil {
		httputil.HandleError(w, err)
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		httputil.HandleError(w, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, map[string]string{"message": "Data successfully deleted"})
}
