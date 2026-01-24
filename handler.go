package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// @Summary Root endpoint
// @Tags Root
// @Success 200 {object} map[string]string
// @Router / [get]
func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", ContentTypeJSON)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Kasir API",
	})
}

// @Summary Health check
// @Tags Health
// @Success 200 {object} map[string]string
// @Router /health [get]
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", ContentTypeJSON)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "OK",
		"pesan":  MsgAPIRunning,
	})
}

// @Summary Get all products
// @Tags Products
// @Success 200 {array} Product
// @Router /api/products [get]
func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", ContentTypeJSON)
	json.NewEncoder(w).Encode(repo.FindAll())
}

// @Summary Create product
// @Tags Products
// @Param product body Product true "Product data"
// @Success 201 {object} Product
// @Failure 400 {string} string
// @Router /api/products [post]
func createProduct(w http.ResponseWriter, r *http.Request) {
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

// @Summary Get product by ID
// @Tags Products
// @Param id path int true "Product ID"
// @Success 200 {object} Product
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Router /api/products/{id} [get]
func getProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, MsgInvalidID, http.StatusBadRequest)
		return
	}

	product, found := repo.FindByID(id)
	if !found {
		http.Error(w, MsgNotFound, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", ContentTypeJSON)
	json.NewEncoder(w).Encode(product)
}

// @Summary Update product
// @Tags Products
// @Param id path int true "Product ID"
// @Param product body Product true "Product data"
// @Success 200 {object} Product
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Router /api/products/{id} [put]
func updateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, MsgInvalidID, http.StatusBadRequest)
		return
	}

	var updatedProduct Product
	err = json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		http.Error(w, MsgInvalidRequest, http.StatusBadRequest)
		return
	}

	product, found := repo.Update(id, updatedProduct)
	if !found {
		http.Error(w, MsgNotFound, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", ContentTypeJSON)
	json.NewEncoder(w).Encode(product)
}

// @Summary Delete product
// @Tags Products
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Router /api/products/{id} [delete]
func deleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, MsgInvalidID, http.StatusBadRequest)
		return
	}

	if !repo.Delete(id) {
		http.Error(w, MsgNotFound, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", ContentTypeJSON)
	json.NewEncoder(w).Encode(map[string]string{
		"pesan": MsgDeleteSuccess,
	})
}

// @Summary Get all categories
// @Tags Categories
// @Success 200 {array} Category
// @Router /api/categories [get]
func getCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", ContentTypeJSON)
	json.NewEncoder(w).Encode(categoryRepo.FindAll())
}

// @Summary Create category
// @Tags Categories
// @Param category body Category true "Category data"
// @Success 201 {object} Category
// @Failure 400 {string} string
// @Router /api/categories [post]
func createCategory(w http.ResponseWriter, r *http.Request) {
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

// @Summary Get category by ID
// @Tags Categories
// @Param id path int true "Category ID"
// @Success 200 {object} Category
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Router /api/categories/{id} [get]
func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, MsgInvalidID, http.StatusBadRequest)
		return
	}

	category, found := categoryRepo.FindByID(id)
	if !found {
		http.Error(w, MsgNotFound, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", ContentTypeJSON)
	json.NewEncoder(w).Encode(category)
}

// @Summary Update category
// @Tags Categories
// @Param id path int true "Category ID"
// @Param category body Category true "Category data"
// @Success 200 {object} Category
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Router /api/categories/{id} [put]
func updateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, MsgInvalidID, http.StatusBadRequest)
		return
	}

	var updatedCategory Category
	err = json.NewDecoder(r.Body).Decode(&updatedCategory)
	if err != nil {
		http.Error(w, MsgInvalidRequest, http.StatusBadRequest)
		return
	}

	category, found := categoryRepo.Update(id, updatedCategory)
	if !found {
		http.Error(w, MsgNotFound, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", ContentTypeJSON)
	json.NewEncoder(w).Encode(category)
}

// @Summary Delete category
// @Tags Categories
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Router /api/categories/{id} [delete]
func deleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, MsgInvalidID, http.StatusBadRequest)
		return
	}

	if !categoryRepo.Delete(id) {
		http.Error(w, MsgNotFound, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", ContentTypeJSON)
	json.NewEncoder(w).Encode(map[string]string{
		"pesan": MsgDeleteSuccess,
	})
}
