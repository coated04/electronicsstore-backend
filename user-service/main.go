package main

import (
	"device-store/user-service/handlers"
	"device-store/user-service/middleware"
	"device-store/user-service/models"
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

	db.AutoMigrate(&models.User{})

	h := &handlers.Handler{DB: db}

	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)

	r.HandleFunc("/register", h.Register).Methods("POST")
	r.HandleFunc("/login", h.Login).Methods("POST")


	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	log.Println("User Service is running on http://localhost:8001")
	http.ListenAndServe(":8001", handler)
}
