// File: internal/router/router.go
// Purpose: HTTP route registration, separated from main.go for testability.
// Author: Backend Team
// Created: 2026-09-25
// Last Modified: 2026-09-25

package router

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/adriel-meb/agnambie-backend/internal/handler"
	"github.com/adriel-meb/agnambie-backend/internal/middleware"
)

// New creates a fully configured chi.Router with all middleware and routes.
// Separating route registration from main.go makes the routing testable
// and scannable in one file.
func New(h *handler.Handler, corsOrigins []string, rlReqsPerSec float64, rlBurst int, logger *slog.Logger) chi.Router {
	r := chi.NewRouter()

	// ── Global middleware stack (order matters) ──────────────────────────
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(middleware.RequestLogger(logger))
	r.Use(middleware.CORS(corsOrigins))
	r.Use(middleware.Compress)     // Gzip — critical for low-bandwidth clients
	r.Use(middleware.SecurityHeaders)
	r.Use(chimw.Recoverer)

	// ── Probes — outside /api for load balancer access ───────────────────
	r.Get("/health", h.Health)
	r.Get("/ready", h.Ready)

	// ── API routes ──────────────────────────────────────────────────────
	apiRoutes := func(r chi.Router) {
		r.Use(middleware.RateLimit(rlReqsPerSec, rlBurst))
		r.Get("/languages", h.Languages)
		r.Get("/bibles", h.Bibles)
		r.Get("/books", h.Books)
		r.Get("/audio", h.Audio)
		r.Get("/copyright", h.Copyright)
		r.Post("/metrics", h.Metrics)
	}

	r.Route("/api", apiRoutes)
	r.Route("/api/v1", apiRoutes)

	return r
}
