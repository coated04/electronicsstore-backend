package handlers

import (
	"encoding/json"
	"device-store/device-service/models"
	"github.com/go-resty/resty/v2"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"strconv"
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

func (h *Handler) GetUserFromDevice(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var device models.Device
	if err := h.DB.First(&device, id).Error; err != nil {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}

	userIDStr := strconv.FormatUint(uint64(device.UserID), 10)

	client := resty.New()
	userResponse, err := client.R().Get("http://localhost:8001/users/" + userIDStr)

	if err != nil || userResponse.StatusCode() != http.StatusOK {
		http.Error(w, "User service is unavailable", http.StatusServiceUnavailable)
		return
	}

	w.Write(userResponse.Body())
}
