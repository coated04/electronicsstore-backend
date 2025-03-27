package handlers

import (
    "encoding/json"
    "net/http"
    "bookstore/models"
)

var categories = []models.Category{}
var nextCategoryID = 1

func GetCategories(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(categories)
}

func CreateCategory(w http.ResponseWriter, r *http.Request) {
    var category models.Category
    json.NewDecoder(r.Body).Decode(&category)
    if category.Name == "" {
        http.Error(w, "name is required", http.StatusBadRequest)
        return
    }
    category.ID = nextCategoryID
    nextCategoryID++
    categories = append(categories, category)
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(category)
}
