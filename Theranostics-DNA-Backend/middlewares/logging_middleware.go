// middlewares/logging_middleware.go
package middlewares

import (
    "log"
    "net/http"
    "time"
)

// LoggingMiddleware logs each incoming HTTP request with details such as method, URI, status code, and duration.
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Capture the start time
        start := time.Now()

        // Use a ResponseWriter wrapper to capture the status code
        rw := &responseWriter{w, http.StatusOK}

        // Process the request
        next.ServeHTTP(rw, r)

        // Calculate the duration
        duration := time.Since(start)

        // Log the details
        log.Printf(
            "%s %s %d %s",
            r.Method,
            r.RequestURI,
            rw.statusCode,
            duration,
        )
    })
}

// responseWriter is a wrapper around http.ResponseWriter to capture the status code
type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

// WriteHeader captures the status code for logging
func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}
