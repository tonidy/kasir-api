package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"kasir-api/internal/model"
	"kasir-api/pkg/httputil"
)

type TransactionService interface {
	Checkout(ctx context.Context, req model.CheckoutRequest) (*model.Transaction, error)
}

type TransactionHandler struct {
	svc TransactionService
}

func NewTransactionHandler(svc TransactionService) *TransactionHandler {
	return &TransactionHandler{svc: svc}
}

func (h *TransactionHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	var req model.CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	transaction, err := h.svc.Checkout(r.Context(), req)
	if err != nil {
		httputil.WriteError(w, httputil.ErrorStatus(err), err.Error())
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, transaction)
}
