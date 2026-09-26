// File: internal/handler/books.go
// Purpose: HTTP handler for GET /api/books.
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
	"github.com/adriel-meb/agnambie-backend/internal/logger"
)

// Books returns the list of books available in a given Bible.
// Query params:
//   - bible_id (required): Bible ID/abbreviation, e.g. "FANBSG"
//
// The bible_id is validated against the Gabon allowlist server-side.
func (h *Handler) Books(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	bibleID := strings.TrimSpace(r.URL.Query().Get("bible_id"))
	if bibleID == "" {
		writeError(w, http.StatusBadRequest, "bible_id is required")
		return
	}

	// Enforce the Gabon allowlist.
	if !domain.IsAllowedBibleID(bibleID) {
		writeError(w, http.StatusForbidden, "bible not available")
		return
	}

	cacheKey := "books:" + bibleID

	raw, hit, err := h.cache.GetOrSet(ctx, cacheKey, 1*time.Hour, func() ([]byte, error) {
		return h.bb.FetchBooks(ctx, bibleID)
	})
	if err != nil {
		logger.FromContext(ctx).Error("failed to fetch books",
			slog.String("bible_id", bibleID),
			slog.String("error", err.Error()),
		)
		writeError(w, http.StatusBadGateway, "could not fetch books")
		return
	}

	var items []json.RawMessage
	_ = json.Unmarshal(raw, &items)

	data, _ := json.Marshal(Envelope{
		Data: json.RawMessage(raw),
		Meta: &Meta{Count: len(items)},
	})

	writeData(w, r, data, hit, 3600)
}
