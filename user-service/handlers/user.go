package handlers

import (
	"encoding/json"
	"device-store/user-service/models"
	"device-store/user-service/utils"
	"net/http"

	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	DB *gorm.DB
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error encrypting password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	if err := h.DB.Create(&user).Error; err != nil {
		http.Error(w, "Error saving user", http.StatusInternalServerError)
		return
	}

	user.Password = "" 
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    var user models.User
    if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
        http.Error(w, "Invalid email or password", http.StatusUnauthorized)
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        http.Error(w, "Invalid email or password", http.StatusUnauthorized)
        return
    }

    token, err := utils.GenerateJWT(user.ID)
    if err != nil {
        http.Error(w, "Error generating token", http.StatusInternalServerError)
        return
    }

    // Prepare response without password
    responseUser := struct {
        ID       uint   `json:"id"`
        Username string `json:"username"`
        Email    string `json:"email"`
    }{
        ID:       user.ID,
        Username: user.Username,
        Email:    user.Email,
    }

    // Send token + user info in JSON
    json.NewEncoder(w).Encode(map[string]interface{}{
        "token": token,
        "user":  responseUser,
    })
}
