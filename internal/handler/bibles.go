// File: internal/handler/bibles.go
// Purpose: HTTP handler for GET /api/bibles.
// Author: Backend Team
// Created: 2026-09-25
// Last Modified: 2026-09-25

package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/adriel-meb/agnambie-backend/internal/domain"
)

// Bibles returns Bible translations available for a given Gabon language.
// Query params:
//   - language_code (required): ISO 639-3 code, e.g. "FAN", "MYE", "FRA"
//
// The language_code is validated against the Gabon allowlist server-side.
func (h *Handler) Bibles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	langCode := strings.TrimSpace(strings.ToUpper(r.URL.Query().Get("language_code")))
	if langCode == "" {
		writeError(w, http.StatusBadRequest, "language_code is required")
		return
	}

	// Enforce the Gabon language allowlist server-side.
	if !domain.IsAllowedLanguage(langCode) {
		writeError(w, http.StatusForbidden, "language not available")
		return
	}

	cacheKey := "bibles:" + langCode

	raw, hit, err := h.cache.GetOrSet(ctx, cacheKey, 1*time.Hour, func() ([]byte, error) {
		return h.bb.FetchBibles(ctx, langCode)
	})
	if err != nil {
		h.logger.Error("failed to fetch bibles",
			slog.String("language_code", langCode),
			slog.String("error", err.Error()),
		)
		writeError(w, http.StatusBadGateway, "could not fetch bibles")
		return
	}

	// Wrap raw Bible Brain data in our envelope.
	// Parse count from raw JSON array for the meta field.
	var items []json.RawMessage
	_ = json.Unmarshal(raw, &items)

	data, _ := json.Marshal(Envelope{
		Data: json.RawMessage(raw),
		Meta: &Meta{Count: len(items)},
	})

	// 1 hour client-side cache, matching server-side cache TTL.
	writeData(w, r, data, hit, 3600)
}


