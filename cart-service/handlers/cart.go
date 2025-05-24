package handlers

import (
	"device-store/cart-service/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type Handler struct {
	DB *gorm.DB
}

func (h *Handler) AddToCart(w http.ResponseWriter, r *http.Request) {
	var item models.CartItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if err := h.DB.Create(&item).Error; err != nil {
		http.Error(w, "Failed to add item", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(item)
}

func (h *Handler) GetCart(w http.ResponseWriter, r *http.Request) {
	userIDStr := mux.Vars(r)["userID"]
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var items []models.CartItem
	if err := h.DB.Where("user_id = ?", userID).Find(&items).Error; err != nil {
		http.Error(w, "Failed to retrieve cart", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(items)
}

func (h *Handler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid cart item ID", http.StatusBadRequest)
		return
	}
	if err := h.DB.Delete(&models.CartItem{}, id).Error; err != nil {
		http.Error(w, "Failed to delete item", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
