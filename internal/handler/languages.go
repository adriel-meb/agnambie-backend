// File: internal/handler/languages.go
// Purpose: HTTP handler for GET /api/languages.
// Author: Backend Team
// Created: 2026-09-25
// Last Modified: 2026-09-25

package handler

import (
	"encoding/json"
	"net/http"

	"github.com/adriel-meb/agnambie-backend/internal/domain"
)

// Languages returns the hardcoded Gabon language list and their filesets.
// No network call — data is static from domain/gabon.go.
// This is the first endpoint the Flutter app calls on launch.
//
// Supports ETag / If-None-Match for 304 Not Modified — mobile clients
// skip re-downloading the same language list on subsequent launches.
func (h *Handler) Languages(w http.ResponseWriter, r *http.Request) {
	data, _ := json.Marshal(Envelope{
		Data: domain.GabonLanguages,
		Meta: &Meta{Count: len(domain.GabonLanguages)},
	})

	// Static data — cache for 24 hours on client, never changes without deploy.
	writeData(w, r, data, true, 86400)
}
