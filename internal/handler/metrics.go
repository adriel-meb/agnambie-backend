// File: internal/handler/metrics.go
// Purpose: Ingestion endpoint for offline-batched client telemetry events.
// Author: Backend Team
// Created: 2026-09-26
// Last Modified: 2026-09-26

package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

// MetricEvent represents a single client-side telemetry event.
type MetricEvent struct {
	Name       string         `json:"name"`
	Timestamp  time.Time      `json:"timestamp"`
	Properties map[string]any `json:"properties,omitempty"`
}

// Metrics handles batch uploads of telemetry events from mobile clients.
// POST /api/v1/metrics
// POST /api/metrics
func (h *Handler) Metrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, r)
		return
	}

	// Limit request body to 512 KB to prevent memory exhaustion
	r.Body = http.MaxBytesReader(w, r.Body, 512*1024)

	var events []MetricEvent
	if err := json.NewDecoder(r.Body).Decode(&events); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json payload or body too large")
		return
	}

	if len(events) == 0 {
		writeError(w, http.StatusBadRequest, "empty events array")
		return
	}

	// Log batch summary and individual events with structured context
	slog.Info("metrics_batch_received", "count", len(events))
	for _, ev := range events {
		slog.Info("client_event",
			"event", ev.Name,
			"client_time", ev.Timestamp,
			"properties", ev.Properties,
		)
	}

	writeJSON(w, http.StatusOK, Envelope{
		Data: map[string]any{"accepted": len(events)},
		Meta: &Meta{Count: len(events)},
	})
}
