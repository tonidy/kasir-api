package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"kasir-api/pkg/logger"
	"kasir-api/pkg/uuid"
)

// RequestIDKey is the key used to store request ID in context
type RequestIDKey string

const RequestIDCtxKey RequestIDKey = "request_id"

// LoggingMiddleware logs incoming requests and responses
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Generate request ID
		requestID := uuid.New()
		ctx := context.WithValue(r.Context(), RequestIDCtxKey, requestID)
		r = r.WithContext(ctx)

		// Log request
		logger.Info("Incoming request",
			"request_id", requestID,
			"method", r.Method,
			"url", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)

		// Wrap ResponseWriter to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Process request
		next.ServeHTTP(wrapped, r)

		// Log response
		duration := time.Since(start)
		logger.Info("Request completed",
			"request_id", requestID,
			"method", r.Method,
			"url", r.URL.Path,
			"status_code", wrapped.statusCode,
			"duration", duration.Milliseconds(),
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// CORSMiddleware adds CORS headers
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware provides basic authentication
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for health and docs endpoints
		if strings.HasPrefix(r.URL.Path, "/health") ||
			strings.HasPrefix(r.URL.Path, "/docs/") ||
			r.URL.Path == "/" {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Basic validation - in real app, validate token properly
		if !isValidToken(authHeader) {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// isValidToken validates the authentication token
func isValidToken(token string) bool {
	// In a real application, this would validate against a database or JWT
	// For now, we'll accept any token that starts with "Bearer "
	return strings.HasPrefix(token, "Bearer ")
}

// RecoveryMiddleware recovers from panics
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered",
					"error", fmt.Sprintf("%v", err),
					"url", r.URL.Path,
					"method", r.Method,
				)

				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
