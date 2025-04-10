package main

import (
    "device-store/handlers" 
    "device-store/models"    
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=2001 dbname=devicestore port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	db.AutoMigrate(&models.Device{}, &models.Brand{}, &models.Category{})

	h := &handlers.Handler{DB: db}

	r := mux.NewRouter()

	r.HandleFunc("/devices", h.GetDevices).Methods("GET")
	r.HandleFunc("/devices/{id}", h.GetDeviceByID).Methods("GET")
	r.HandleFunc("/devices", h.CreateDevice).Methods("POST")
	r.HandleFunc("/devices/{id}", h.UpdateDevice).Methods("PUT")
	r.HandleFunc("/devices/{id}", h.DeleteDevice).Methods("DELETE")


	r.HandleFunc("/brands", h.GetBrands).Methods("GET")
	r.HandleFunc("/brands/{id}", h.GetBrandByID).Methods("GET")
	r.HandleFunc("/brands", h.CreateBrand).Methods("POST")
	r.HandleFunc("/brands/{id}", h.UpdateBrand).Methods("PUT")
	r.HandleFunc("/brands/{id}", h.DeleteBrand).Methods("DELETE")

	r.HandleFunc("/categories", h.GetCategories).Methods("GET")
	r.HandleFunc("/categories/{id}", h.GetCategoryByID).Methods("GET")
	r.HandleFunc("/categories", h.CreateCategory).Methods("POST")
	r.HandleFunc("/categories/{id}", h.UpdateCategory).Methods("PUT")
	r.HandleFunc("/categories/{id}", h.DeleteCategory).Methods("DELETE")

	log.Println("Server is running on http://localhost:8000")
	http.ListenAndServe(":8000", r)
}
