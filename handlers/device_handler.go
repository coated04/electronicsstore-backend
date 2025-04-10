package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"device-store/models"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func (h *Handler) GetDevices(w http.ResponseWriter, r *http.Request) {
	var devices []models.Device

	brandID := r.URL.Query().Get("brand_id")
	categoryID := r.URL.Query().Get("category_id")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	query := h.DB.Limit(limit).Offset(offset)
	if brandID != "" {
		query = query.Where("brand_id = ?", brandID)
	}
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	query.Find(&devices)
	json.NewEncoder(w).Encode(devices)
}

func (h *Handler) GetDeviceByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var device models.Device
	if err := h.DB.First(&device, id).Error; err != nil {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(device)
}

func (h *Handler) CreateDevice(w http.ResponseWriter, r *http.Request) {
	var device models.Device
	if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if device.Name == "" || device.Price <= 0 || device.BrandID == 0 || device.CategoryID == 0 {
		http.Error(w, "Missing or invalid fields", http.StatusBadRequest)
		return
	}

	h.DB.Create(&device)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(device)
}

func (h *Handler) UpdateDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var existing models.Device
	if err := h.DB.First(&existing, id).Error; err != nil {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}

	var updated models.Device
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if updated.Name == "" || updated.Price <= 0 || updated.BrandID == 0 || updated.CategoryID == 0 {
		http.Error(w, "Missing or invalid fields", http.StatusBadRequest)
		return
	}

	updated.ID = existing.ID
	h.DB.Save(&updated)
	json.NewEncoder(w).Encode(updated)
}

func (h *Handler) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.DB.Delete(&models.Device{}, id).Error; err != nil {
		http.Error(w, "Failed to delete", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
