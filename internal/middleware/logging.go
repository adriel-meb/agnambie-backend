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
)

// wrappedWriter captures the HTTP status code written by downstream handlers.
type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code before delegating to the underlying writer.
func (w *wrappedWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// RequestLogger returns middleware that logs every completed HTTP request
// with structured fields: method, path, status, duration, and remote address.
func RequestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapped := &wrappedWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)
			logger.Info("request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", wrapped.statusCode),
				slog.Duration("duration", duration),
				slog.String("remote", r.RemoteAddr),
			)
		})
	}
}
