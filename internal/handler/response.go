// File: internal/handler/response.go
// Purpose: Standardized JSON response envelope, helpers, and ETag support.
// Author: Backend Team
// Created: 2026-09-25
// Last Modified: 2026-09-25

package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
)

// Envelope wraps all successful API responses in a predictable structure.
// The Flutter client can check for "data" (success) vs "error" (failure)
// without inspecting HTTP status codes — critical for offline caching.
type Envelope struct {
	Data any   `json:"data"`
	Meta *Meta `json:"meta,omitempty"`
}

// Meta provides response metadata useful for the client.
type Meta struct {
	Count int `json:"count,omitempty"` // Number of items in data array
}

// ErrorResponse is the structured error envelope.
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains machine-readable code and human-readable message.
type ErrorDetail struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// --- Response writers ---

// writeJSON serializes data in the standard envelope and sends it.
// Uses compact JSON (no trailing newline) to save bytes for mobile clients.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	// json.Marshal produces compact JSON without trailing newline.
	b, err := json.Marshal(data)
	if err != nil {
		// Fallback: if marshalling fails, send a minimal error.
		w.Write([]byte(`{"error":{"code":500,"message":"internal error"}}`))
		return
	}
	w.Write(b)
}

// writeData sends data wrapped in the standard envelope with cache metadata.
// Adds ETag and Cache-Control headers for client-side caching.
func writeData(w http.ResponseWriter, r *http.Request, data []byte, hit bool, cacheMaxAge int) {
	// Set cache-related headers.
	setCacheHeader(w, hit)

	// Generate ETag from content hash — allows 304 Not Modified responses.
	etag := generateETag(data)
	w.Header().Set("ETag", etag)

	// If the client already has this version, send 304 to save bandwidth.
	if r.Header.Get("If-None-Match") == etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	// Set Cache-Control for HTTP client caching.
	if cacheMaxAge > 0 {
		w.Header().Set("Cache-Control", "public, max-age="+intToStr(cacheMaxAge))
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// writeError sends a structured error response with code and message.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, ErrorResponse{
		Error: ErrorDetail{
			Code:    status,
			Message: msg,
		},
	})
}

// setCacheHeader sets the X-Cache header based on whether the response was cached.
func setCacheHeader(w http.ResponseWriter, hit bool) {
	if hit {
		w.Header().Set("X-Cache", "HIT")
	} else {
		w.Header().Set("X-Cache", "MISS")
	}
}

// generateETag creates a weak ETag from a SHA-256 hash of the content.
// Weak ETags are appropriate since our content may be serialized differently.
func generateETag(data []byte) string {
	h := sha256.Sum256(data)
	return `W/"` + hex.EncodeToString(h[:8]) + `"`
}

// intToStr converts an int to its string representation without importing strconv
// in this small file. Only handles positive integers.
func intToStr(n int) string {
	if n == 0 {
		return "0"
	}
	digits := make([]byte, 0, 10)
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}
