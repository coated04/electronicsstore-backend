package main

import (
	"device-store/device-service/handlers"
	"device-store/device-service/middleware"
	"device-store/device-service/models"
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

	db.AutoMigrate(&models.Device{})

	h := &handlers.Handler{DB: db}

	r := mux.NewRouter()

	r.Use(middleware.LoggingMiddleware)

	r.HandleFunc("/devices", h.GetDevices).Methods("GET")
	r.HandleFunc("/devices/{id}", h.GetDeviceByID).Methods("GET")
	r.HandleFunc("/devices", h.CreateDevice).Methods("POST")
	r.HandleFunc("/devices/{id}", h.UpdateDevice).Methods("PUT")
	r.HandleFunc("/devices/{id}", h.DeleteDevice).Methods("DELETE")

	r.HandleFunc("/devices/{id}/user", h.GetUserFromDevice).Methods("GET")


	log.Println("Device Service is running on http://localhost:8002")
	http.ListenAndServe(":8002", r)
}
