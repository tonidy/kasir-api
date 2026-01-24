package main

import (
	"encoding/json"
	"net/http"
)

func setupRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", ContentTypeJSON)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Kasir API",
		})
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", ContentTypeJSON)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "OK",
			"pesan":  MsgAPIRunning,
		})
	})

	http.HandleFunc("/api/products/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getProductByID(w, r)
		case http.MethodPut:
			updateProduct(w, r)
		case http.MethodDelete:
			deleteProduct(w, r)
		}
	})

	http.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", ContentTypeJSON)
			json.NewEncoder(w).Encode(repo.FindAll())
		case http.MethodPost:
			var newProduct Product
			err := json.NewDecoder(r.Body).Decode(&newProduct)
			if err != nil {
				http.Error(w, MsgInvalidRequest, http.StatusBadRequest)
				return
			}

			created := repo.Create(newProduct)
			w.Header().Set("Content-Type", ContentTypeJSON)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(created)
		}
	})

	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getCategoryByID(w, r)
		case http.MethodPut:
			updateCategory(w, r)
		case http.MethodDelete:
			deleteCategory(w, r)
		}
	})

	http.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", ContentTypeJSON)
			json.NewEncoder(w).Encode(categoryRepo.FindAll())
		case http.MethodPost:
			var newCategory Category
			err := json.NewDecoder(r.Body).Decode(&newCategory)
			if err != nil {
				http.Error(w, MsgInvalidRequest, http.StatusBadRequest)
				return
			}

			created := categoryRepo.Create(newCategory)
			w.Header().Set("Content-Type", ContentTypeJSON)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(created)
		}
	})
}
