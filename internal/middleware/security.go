// File: internal/middleware/security.go
// Purpose: Security response headers middleware.
// Author: Backend Team
// Created: 2026-09-25
// Last Modified: 2026-09-25

package middleware

import "net/http"

// SecurityHeaders returns middleware that adds standard security headers
// to every response, preventing common web vulnerabilities.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Content-Security-Policy", "default-src 'none'")

		next.ServeHTTP(w, r)
	})
}
