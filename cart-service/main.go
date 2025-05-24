package main

import (
    "device-store/cart-service/handlers"
    "device-store/cart-service/middleware"
    "device-store/cart-service/models"
    "fmt"
    "github.com/gorilla/mux"
    "github.com/rs/cors"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "log"
    "net/http"
    "os"
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
    // DB configuration
    host := os.Getenv("DB_HOST")
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME")
    port := 5432

    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
        host, user, password, dbname, port,
    )

    db := connectWithRetry(dsn)

    err := db.AutoMigrate(&models.CartItem{})
    if err != nil {
        log.Fatalf("Failed to migrate DB: %v", err)
    }

    h := &handlers.Handler{DB: db}

    r := mux.NewRouter()
    r.Use(middleware.LoggingMiddleware)

    r.HandleFunc("/cart", h.AddToCart).Methods("POST")
    r.HandleFunc("/cart/{userID}", h.GetCart).Methods("GET")
    r.HandleFunc("/cart/{id}", h.RemoveFromCart).Methods("DELETE")

    // Add CORS wrapper
    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:3000"}, // Allow your frontend
        AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
        AllowedHeaders:   []string{"Authorization", "Content-Type"},
        AllowCredentials: true,
    })

    handler := c.Handler(r) // Wrap the router with CORS

    log.Println("Cart service running on :8003")
    err = http.ListenAndServe(":8003", handler) // Use CORS-wrapped handler
    if err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
