// File: internal/middleware/ratelimit.go
// Purpose: Token-bucket rate limiting middleware based on client IP.
// Author: Backend Team
// Created: 2026-09-26
// Last Modified: 2026-09-26

package middleware

import (
	"encoding/json"
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

// rateLimiter struct holds a map of rate limiters per IP.
type rateLimiter struct {
	mu       sync.Mutex
	limiters map[string]*rate.Limiter
	r        rate.Limit
	b        int
}

func newRateLimiter(r rate.Limit, b int) *rateLimiter {
	return &rateLimiter{
		limiters: make(map[string]*rate.Limiter),
		r:        r,
		b:        b,
	}
}

// getLimiter returns the rate limiter for the provided IP address.
func (rl *rateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.r, rl.b)
		rl.limiters[ip] = limiter
	}
	return limiter
}

// RateLimit is a middleware that implements token bucket rate limiting.
func RateLimit(reqsPerSec float64, burst int) func(http.Handler) http.Handler {
	limiter := newRateLimiter(rate.Limit(reqsPerSec), burst)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract IP. Ideally this comes from X-Forwarded-For or RemoteAddr
			// if behind a reverse proxy.
			ip := r.RemoteAddr
			if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
				ip = forwarded
			}

			if !limiter.getLimiter(ip).Allow() {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(map[string]any{
					"error": map[string]any{
						"code":    http.StatusTooManyRequests,
						"message": "Too many requests. Please slow down.",
					},
				})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
