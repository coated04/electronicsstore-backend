package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		start := time.Now()
		log.Printf("[Start] [RequestID: %s] %s %s", requestID, r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		duration := time.Since(start).Milliseconds()
		log.Printf("[End] [RequestID: %s] %s %s - Duration: %dms", requestID, r.Method, r.URL.Path, duration)
	})
}
