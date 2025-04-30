package main

import (
	"device-store/user-service/handlers"
	"device-store/user-service/middleware"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"device-store/user-service/models"
)

func main() {
	dsn := "host=localhost user=postgres password=2001 dbname=devicestore port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	db.AutoMigrate(&models.User{})

	h := &handlers.Handler{DB: db}

	r := mux.NewRouter()

	
	r.Use(middleware.LoggingMiddleware)

	r.HandleFunc("/users", h.GetUsers).Methods("GET")
	r.HandleFunc("/users/{id}", h.GetUserByID).Methods("GET")
	r.HandleFunc("/users", h.CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", h.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", h.DeleteUser).Methods("DELETE")

	log.Println("User Service is running on http://localhost:8001")
	http.ListenAndServe(":8001", r)
}
