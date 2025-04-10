package handlers

import (
	"device-store/models"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) GetCategories(w http.ResponseWriter, r *http.Request) {
	var categories []models.Category
	h.DB.Find(&categories)
	json.NewEncoder(w).Encode(categories)
}

func (h *Handler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var category models.Category
	if err := h.DB.First(&category, id).Error; err != nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(category)
}

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if category.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	h.DB.Create(&category)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var category models.Category
	if err := h.DB.First(&category, id).Error; err != nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}
	var updated models.Category
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	category.Name = updated.Name
	h.DB.Save(&category)
	json.NewEncoder(w).Encode(category)
}

func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.DB.Delete(&models.Category{}, id).Error; err != nil {
		http.Error(w, "Failed to delete", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
