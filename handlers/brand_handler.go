package handlers

import (
	"device-store/models"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) GetBrands(w http.ResponseWriter, r *http.Request) {
	var brands []models.Brand
	h.DB.Find(&brands)
	json.NewEncoder(w).Encode(brands)
}

func (h *Handler) GetBrandByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var brand models.Brand
	if err := h.DB.First(&brand, id).Error; err != nil {
		http.Error(w, "Brand not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(brand)
}

func (h *Handler) CreateBrand(w http.ResponseWriter, r *http.Request) {
	var brand models.Brand
	if err := json.NewDecoder(r.Body).Decode(&brand); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if brand.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	h.DB.Create(&brand)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(brand)
}

func (h *Handler) UpdateBrand(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var brand models.Brand
	if err := h.DB.First(&brand, id).Error; err != nil {
		http.Error(w, "Brand not found", http.StatusNotFound)
		return
	}
	var updated models.Brand
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	brand.Name = updated.Name
	h.DB.Save(&brand)
	json.NewEncoder(w).Encode(brand)
}

func (h *Handler) DeleteBrand(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.DB.Delete(&models.Brand{}, id).Error; err != nil {
		http.Error(w, "Failed to delete", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
