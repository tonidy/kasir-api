package main

import (
	_ "embed"
	"net/http"
)

//go:embed docs/index.html
var docsHTML []byte

func setupRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		rootHandler(w, r)
	})

	http.HandleFunc("/health", healthHandler)

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
			getProducts(w, r)
		case http.MethodPost:
			createProduct(w, r)
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
			getCategories(w, r)
		case http.MethodPost:
			createCategory(w, r)
		}
	})

	http.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(docsHTML)
	})

	http.Handle("/docs/", http.StripPrefix("/docs/", http.FileServer(http.Dir("docs"))))
}
