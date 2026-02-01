package handler

import (
	"context"
	"database/sql"
	"net/http"

	"kasir-api/pkg/httputil"
)

type HealthHandler struct {
	db *sql.DB
}

func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status":  "OK",
		"message": "API is running",
	}

	// Check database connection if available
	if h.db != nil {
		if err := h.db.PingContext(context.Background()); err != nil {
			response["database"] = "disconnected"
			response["status"] = "degraded"
			httputil.WriteJSON(w, http.StatusServiceUnavailable, response)
			return
		}
		response["database"] = "connected"
	} else {
		response["storage"] = "in-memory"
	}

	httputil.WriteJSON(w, http.StatusOK, response)
}

func (h *HealthHandler) Root(w http.ResponseWriter, r *http.Request) {
	httputil.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Kasir API",
	})
}
