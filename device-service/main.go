package main

import (
	"device-store/device-service/handlers"
	"device-store/device-service/middleware"
	"device-store/device-service/models"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

func connectWithRetry(dsn string) *gorm.DB {
	var db *gorm.DB
	var err error

	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			sqlDB, err := db.DB()
			if err == nil && sqlDB.Ping() == nil {
				log.Println("Connected to DB")
				return db
			}
		}
		log.Printf("DB not ready (attempt %d): %v", i+1, err)
		time.Sleep(3 * time.Second)
	}

	log.Fatalf("Failed to connect to DB after retries: %v", err)
	return nil
}

func main() {
	dsn := "host=postgres user=postgres password=2001 dbname=devicestore port=5432 sslmode=disable"
	db := connectWithRetry(dsn)

	db.AutoMigrate(&models.Brand{}, &models.Category{}, &models.Device{})

	h := &handlers.Handler{DB: db}

	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)

	r.HandleFunc("/devices", h.GetDevices).Methods("GET")
	r.HandleFunc("/devices/{id}", h.GetDeviceByID).Methods("GET")
	authMiddleware := middleware.AuthMiddleware
	r.Handle("/devices", authMiddleware(http.HandlerFunc(h.CreateDevice))).Methods("POST")
	r.Handle("/devices/{id}", authMiddleware(http.HandlerFunc(h.UpdateDevice))).Methods("PUT")
	r.Handle("/devices/{id}", authMiddleware(http.HandlerFunc(h.DeleteDevice))).Methods("DELETE")


	r.HandleFunc("/brands", h.GetBrands).Methods("GET")
	r.HandleFunc("/brands/{id}", h.GetBrandByID).Methods("GET")
	r.Handle("/brands", authMiddleware(http.HandlerFunc(h.CreateBrand))).Methods("POST")
	r.Handle("/brands/{id}", authMiddleware(http.HandlerFunc(h.UpdateBrand))).Methods("PUT")
	r.Handle("/brands/{id}", authMiddleware(http.HandlerFunc(h.DeleteBrand))).Methods("DELETE")

	r.HandleFunc("/categories", h.GetCategories).Methods("GET")
	r.HandleFunc("/categories/{id}", h.GetCategoryByID).Methods("GET")
	r.Handle("/categories", authMiddleware(http.HandlerFunc(h.CreateCategory))).Methods("POST")
	r.Handle("/categories/{id}", authMiddleware(http.HandlerFunc(h.UpdateCategory))).Methods("PUT")
	r.Handle("/categories/{id}", authMiddleware(http.HandlerFunc(h.DeleteCategory))).Methods("DELETE")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // frontend origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})


	handler := c.Handler(r)

	log.Println("Device Service running on http://localhost:8002")
	http.ListenAndServe(":8002", handler)
}
