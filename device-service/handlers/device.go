package handlers

import (
	"encoding/json"
	"device-store/device-service/models"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

type Handler struct {
	DB *gorm.DB
}



func (h *Handler) GetDevices(w http.ResponseWriter, r *http.Request) {
	var devices []models.Device
	h.DB.Find(&devices)
	json.NewEncoder(w).Encode(devices)
}

func (h *Handler) GetDeviceByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var device models.Device

	
	if err := h.DB.Preload("Brand").Preload("Category").First(&device, id).Error; err != nil {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(device)
}


func (h *Handler) CreateDevice(w http.ResponseWriter, r *http.Request) {
	var device models.Device
	if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	h.DB.Create(&device)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(device)
}

func (h *Handler) UpdateDevice(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var device models.Device
	if err := h.DB.First(&device, id).Error; err != nil {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}
	var updatedDevice models.Device
	if err := json.NewDecoder(r.Body).Decode(&updatedDevice); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	updatedDevice.ID = device.ID 
	h.DB.Save(&updatedDevice)
	json.NewEncoder(w).Encode(updatedDevice)
}

func (h *Handler) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.DB.Delete(&models.Device{}, id).Error; err != nil {
		http.Error(w, "Failed to delete device", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}





// ======================= Brand Handlers =======================

func (h *Handler) CreateBrand(w http.ResponseWriter, r *http.Request) {
	var brand models.Brand
	if err := json.NewDecoder(r.Body).Decode(&brand); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if err := h.DB.Create(&brand).Error; err != nil {
		http.Error(w, "Failed to create brand", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(brand)
}

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
		http.Error(w, "Failed to delete brand", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ======================= Category Handlers =======================

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if err := h.DB.Create(&category).Error; err != nil {
		http.Error(w, "Failed to create category", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

func (h *Handler) GetCategories(w http.ResponseWriter, r *http.Request) {
	var categories []models.Category
	if err := h.DB.Find(&categories).Error; err != nil {
		http.Error(w, "Failed to retrieve categories", http.StatusInternalServerError)
		return
	}
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
	if err := h.DB.Save(&category).Error; err != nil {
		http.Error(w, "Failed to update category", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(category)
}

func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.DB.Delete(&models.Category{}, id).Error; err != nil {
		http.Error(w, "Failed to delete category", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
