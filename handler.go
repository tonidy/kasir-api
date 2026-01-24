package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

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
