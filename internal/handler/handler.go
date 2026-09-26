// File: internal/handler/handler.go
// Purpose: Shared Handler struct and JSON response helpers for all API endpoints.
// Author: Backend Team
// Created: 2026-09-25
// Last Modified: 2026-09-25

package handler

import (
	"net/http"

	"github.com/adriel-meb/agnambie-backend/internal/biblebrain"
	"github.com/adriel-meb/agnambie-backend/internal/cache"
)

// Handler holds the dependencies all HTTP handlers share.
type Handler struct {
	bb    *biblebrain.Client
	cache cache.Cacher
}

// NewHandler creates a Handler with its required dependencies.
func NewHandler(bb *biblebrain.Client, c cache.Cacher) *Handler {
	return &Handler{
		bb:    bb,
		cache: c,
	}
}

// methodNotAllowed responds with 405 for unsupported HTTP methods.
// Used by the router to enforce GET-only endpoints.
func methodNotAllowed(w http.ResponseWriter, _ *http.Request) {
	writeError(w, http.StatusMethodNotAllowed, "method not allowed")
}
