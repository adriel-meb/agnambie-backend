// File: internal/handler/copyright.go
// Purpose: HTTP handler for GET /api/copyright.
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

// Copyright returns the copyright information for a Bible.
// Query params:
//   - bible_id (required): Bible ID/abbreviation, e.g. "FANBSG"
//
// Required by Bible Brain's terms of service — copyright must be displayed
// before presenting Bible content to end users.
func (h *Handler) Copyright(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	bibleID := strings.TrimSpace(r.URL.Query().Get("bible_id"))
	if bibleID == "" {
		writeError(w, http.StatusBadRequest, "bible_id is required")
		return
	}

	if !domain.IsAllowedBibleID(bibleID) {
		writeError(w, http.StatusForbidden, "bible not available")
		return
	}

	cacheKey := "copyright:" + bibleID

	raw, hit, err := h.cache.GetOrSet(ctx, cacheKey, 24*time.Hour, func() ([]byte, error) {
		return h.bb.FetchCopyright(ctx, bibleID)
	})
	if err != nil {
		h.logger.Error("failed to fetch copyright",
			slog.String("bible_id", bibleID),
			slog.String("error", err.Error()),
		)
		writeError(w, http.StatusBadGateway, "could not fetch copyright")
		return
	}

	data, _ := json.Marshal(Envelope{
		Data: json.RawMessage(raw),
	})

	// Copyright rarely changes — cache for 24 hours on client.
	writeData(w, r, data, hit, 86400)
}
