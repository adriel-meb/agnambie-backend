// File: internal/handler/audio.go
// Purpose: HTTP handler for GET /api/audio.
// Author: Backend Team
// Created: 2026-09-25
// Last Modified: 2026-09-25

package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/adriel-meb/agnambie-backend/internal/domain"
	"github.com/adriel-meb/agnambie-backend/internal/logger"
)

// bookIDPattern validates USFM book IDs (2-5 uppercase letters/digits).
var bookIDPattern = regexp.MustCompile(`^[A-Z1-9]{2,5}$`)

// Audio returns chapter audio file URL(s) for a given fileset.
// Query params:
//   - fileset_id (required): Fileset ID, e.g. "FANBSGN2DA"
//   - book (required): USFM book ID, e.g. "MAT"
//   - chapter (required): Chapter number, e.g. "1"
//   - quality (optional): "data_saver" for opus16 variant
//
// Also respects the standard Save-Data HTTP header — if the client sends
// Save-Data: on, the handler automatically uses the opus16 data-saver fileset
// when available. This is critical for Gabon's low-bandwidth mobile users.
func (h *Handler) Audio(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query()

	filesetID := strings.TrimSpace(q.Get("fileset_id"))
	book := strings.TrimSpace(strings.ToUpper(q.Get("book")))
	chapterStr := strings.TrimSpace(q.Get("chapter"))

	// Validate required params.
	if filesetID == "" || book == "" || chapterStr == "" {
		writeError(w, http.StatusBadRequest, "fileset_id, book, and chapter are required")
		return
	}

	// Validate book format.
	if !bookIDPattern.MatchString(book) {
		writeError(w, http.StatusBadRequest, "book must be a valid USFM book ID (e.g. MAT, GEN)")
		return
	}

	chapter, err := strconv.Atoi(chapterStr)
	if err != nil || chapter < 1 || chapter > 200 {
		writeError(w, http.StatusBadRequest, "chapter must be a positive integer (1-200)")
		return
	}

	// Data Saver: respect both ?quality=data_saver AND the standard Save-Data header.
	// The Save-Data: on header is sent automatically by Android/Chrome when the user
	// enables data-saving mode in system settings.
	effectiveID := filesetID
	wantDataSaver := q.Get("quality") == "data_saver" || r.Header.Get("Save-Data") == "on"
	if wantDataSaver {
		dataSaverID := filesetID + "-opus16"
		if domain.IsAllowedFileset(dataSaverID) {
			effectiveID = dataSaverID
		}
	}

	// Enforce the Gabon fileset allowlist server-side.
	if !domain.IsAllowedFileset(effectiveID) {
		writeError(w, http.StatusForbidden, "fileset not available")
		return
	}

	cacheKey := "audio:" + effectiveID + ":" + book + ":" + chapterStr

	// Audio CDN URLs may expire, so use a shorter cache TTL.
	raw, hit, err := h.cache.GetOrSet(ctx, cacheKey, 30*time.Minute, func() ([]byte, error) {
		return h.bb.FetchAudio(ctx, effectiveID, book, chapter)
	})
	if err != nil {
		logger.FromContext(ctx).Error("failed to fetch audio",
			slog.String("fileset_id", effectiveID),
			slog.String("book", book),
			slog.Int("chapter", chapter),
			slog.String("error", err.Error()),
		)
		writeError(w, http.StatusBadGateway, "could not fetch audio")
		return
	}

	var items []json.RawMessage
	_ = json.Unmarshal(raw, &items)

	data, _ := json.Marshal(Envelope{
		Data: json.RawMessage(raw),
		Meta: &Meta{Count: len(items)},
	})

	// Shorter client-side cache (15 min) since CDN URLs expire.
	writeData(w, r, data, hit, 900)
}
