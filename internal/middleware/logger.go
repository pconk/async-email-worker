package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type contextKey string

const RequestIDKey contextKey = "requestID"

// responseWriter wrapper untuk menangkap status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func Logger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// 1. Generate Request ID (Correlation ID)
			requestID := uuid.New().String()

			// 2. Masukkan Request ID ke Header Response (opsional, untuk debug client)
			w.Header().Set("X-Request-ID", requestID)

			// 3. Masukkan Request ID ke Context agar bisa dipakai di Handler/DB
			ctx := context.WithValue(r.Context(), RequestIDKey, requestID)

			// 4. Wrap ResponseWriter
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			// 5. Jalankan Handler berikutnya
			next.ServeHTTP(wrapped, r.WithContext(ctx))

			// 6. Log Terstruktur dengan slog + Request ID
			logger.Info("HTTP Request",
				"request_id", requestID,
				"method", r.Method,
				"path", r.URL.Path,
				"status", wrapped.statusCode,
				"duration", time.Since(start),
				"ip", r.RemoteAddr,
			)
		})
	}
}

// GetRequestID adalah helper untuk mengambil ID di dalam Handler
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return "unknown"
}
