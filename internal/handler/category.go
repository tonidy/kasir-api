package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"kasir-api/internal/model"
	"kasir-api/pkg/httputil"
)

type CategoryService interface {
	GetByID(ctx context.Context, id int) (*domain.Category, error)
	GetAll(ctx context.Context) ([]domain.Category, error)
	Create(ctx context.Context, c domain.Category) (*domain.Category, error)
	Update(ctx context.Context, id int, c domain.Category) (*domain.Category, error)
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
		httputil.WriteError(w, httputil.ErrorStatus(err), err.Error())
		return
	}
	httputil.WriteJSON(w, http.StatusOK, categories)
}

func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseID(r)
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	category, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, httputil.ErrorStatus(err), err.Error())
		return
	}
	httputil.WriteJSON(w, http.StatusOK, category)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var category domain.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	created, err := h.svc.Create(r.Context(), category)
	if err != nil {
		httputil.WriteError(w, httputil.ErrorStatus(err), err.Error())
		return
	}
	httputil.WriteJSON(w, http.StatusCreated, created)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseID(r)
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var category domain.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	updated, err := h.svc.Update(r.Context(), id, category)
	if err != nil {
		httputil.WriteError(w, httputil.ErrorStatus(err), err.Error())
		return
	}
	httputil.WriteJSON(w, http.StatusOK, updated)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
