package handler

import (
	"net/http"
)

func SetupRoutes(mux *http.ServeMux, productHandler *ProductHandler, categoryHandler *CategoryHandler, healthHandler *HealthHandler) {
	// Health endpoints
	mux.HandleFunc("/", healthHandler.Root)
	mux.HandleFunc("/health", healthHandler.Check)

	// Documentation
	mux.Handle("/docs/", http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs"))))

	// Product endpoints
	mux.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			productHandler.GetAll(w, r)
		case http.MethodPost:
			productHandler.Create(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/products/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			productHandler.GetByID(w, r)
		case http.MethodPut:
			productHandler.Update(w, r)
		case http.MethodDelete:
			productHandler.Delete(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Category endpoints
	mux.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			categoryHandler.GetAll(w, r)
		case http.MethodPost:
			categoryHandler.Create(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			categoryHandler.GetByID(w, r)
		case http.MethodPut:
			categoryHandler.Update(w, r)
		case http.MethodDelete:
			categoryHandler.Delete(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
