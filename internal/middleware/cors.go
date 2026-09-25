// File: internal/middleware/cors.go
// Purpose: CORS middleware for cross-origin request support.
// Author: Backend Team
// Created: 2026-09-25
// Last Modified: 2026-09-25

package middleware

import (
	"net/http"
	"strings"
)

// CORS returns middleware that sets CORS headers based on the allowed origins.
// In development mode (origins contains "*"), all origins are allowed.
// In production, only the listed origins are permitted.
func CORS(allowedOrigins []string) func(http.Handler) http.Handler {
	// Pre-compute whether wildcard is in the list for fast lookups.
	allowAll := false
	originSet := make(map[string]bool, len(allowedOrigins))
	for _, o := range allowedOrigins {
		if o == "*" {
			allowAll = true
		}
		originSet[strings.TrimRight(o, "/")] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if allowAll {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			} else if origin != "" && originSet[strings.TrimRight(origin, "/")] {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
			w.Header().Set("Access-Control-Max-Age", "86400")

			// Handle preflight requests.
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
