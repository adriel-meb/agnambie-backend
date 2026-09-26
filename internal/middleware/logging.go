// File: internal/middleware/logging.go
// Purpose: Structured HTTP request logging middleware using log/slog.
// Author: Backend Team
// Created: 2026-09-24
// Last Modified: 2026-09-25

package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/adriel-meb/agnambie-backend/internal/logger"
	chimw "github.com/go-chi/chi/v5/middleware"
)

// wrappedWriter captures the HTTP status code and bytes written by downstream handlers.
type wrappedWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

// WriteHeader captures the status code before delegating to the underlying writer.
func (w *wrappedWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *wrappedWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.bytesWritten += n
	return n, err
}

// RequestLogger returns middleware that logs every completed HTTP request
// and injects a request-scoped logger (with request_id) into the request context.
func RequestLogger(baseLogger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapped := &wrappedWriter{ResponseWriter: w, statusCode: http.StatusOK}

			// Extract Chi Request ID if available
			reqID := chimw.GetReqID(r.Context())
			
			// Create a request-scoped logger
			reqLogger := baseLogger
			if reqID != "" {
				reqLogger = baseLogger.With(slog.String("request_id", reqID))
			}

			// Inject into context
			ctx := logger.WithLogger(r.Context(), reqLogger)
			r = r.WithContext(ctx)

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)

			attrs := []any{
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
			}

			if r.URL.RawQuery != "" {
				attrs = append(attrs, slog.String("query", r.URL.RawQuery))
			}
			
			attrs = append(attrs,
				slog.Int("status", wrapped.statusCode),
				slog.Int("bytes", wrapped.bytesWritten),
				slog.Duration("duration", duration),
				slog.String("remote", r.RemoteAddr),
			)

			// Log at Info level for all requests. The handlers will log specific application
			// errors with more context. This acts purely as an access log.
			reqLogger.Info("request completed", attrs...)
		})
	}
}
