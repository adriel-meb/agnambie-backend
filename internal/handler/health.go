// File: internal/handler/health.go
// Purpose: HTTP handlers for health/readiness probes.
// Author: Backend Team
// Created: 2026-09-25
// Last Modified: 2026-09-25

package handler

import (
	"net/http"
)

// Health is a liveness probe — returns 200 OK if the process is running.
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Ready is a readiness probe — verifies the cache backend is reachable.
// Returns 503 if the cache is unreachable so load balancers stop routing traffic.
func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	if err := h.cache.Ping(r.Context()); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{
			"status": "not ready",
			"error":  "cache backend unreachable",
		})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}
