// File: internal/handler/playlist.go
// Purpose: HTTP handler for GET /api/audio/playlist — serves key-free HLS
//          playlists so the Bible Brain API key never reaches clients.
// Author: Backend Team
// Created: 2026-10-08
// Last Modified: 2026-10-08

package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/adriel-meb/agnambie-backend/internal/biblebrain"
	"github.com/adriel-meb/agnambie-backend/internal/domain"
	"github.com/adriel-meb/agnambie-backend/internal/logger"
)

// playlistCacheTTL is deliberately shorter than the audio JSON cache: the
// playlist embeds signed CDN segment URLs, which eventually expire.
const playlistCacheTTL = 10 * time.Minute

// Playlist returns a client-safe HLS playlist for one chapter.
// Query params (all required): fileset_id, book, chapter.
//
// The /audio endpoint links here instead of returning the upstream playlist
// URL, which would require the secret API key.
func (h *Handler) Playlist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	p, msg := parseChapterParams(r.URL.Query())
	if msg != "" {
		writeError(w, http.StatusBadRequest, msg)
		return
	}

	// Enforce the Gabon fileset allowlist server-side, same as /audio.
	if !domain.IsAllowedFileset(p.filesetID) {
		writeError(w, http.StatusForbidden, "fileset not available")
		return
	}

	cacheKey := "playlist:" + p.filesetID + ":" + p.book + ":" + strconv.Itoa(p.chapter)
	raw, hit, err := h.cache.GetOrSet(ctx, cacheKey, playlistCacheTTL, func() ([]byte, error) {
		return h.bb.FetchPlaylist(ctx, p.filesetID, p.book, p.chapter)
	})
	if errors.Is(err, biblebrain.ErrNoPlaylist) {
		writeError(w, http.StatusNotFound, "no playlist for this chapter")
		return
	}
	if err != nil {
		logger.FromContext(ctx).Error("failed to fetch playlist",
			slog.String("fileset_id", p.filesetID),
			slog.String("book", p.book),
			slog.Int("chapter", p.chapter),
			slog.String("error", err.Error()),
		)
		writeError(w, http.StatusBadGateway, "could not fetch playlist")
		return
	}

	setCacheHeader(w, hit)
	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	// Short client cache: segment signatures expire upstream.
	w.Header().Set("Cache-Control", "private, max-age=300")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(raw); err != nil {
		logger.FromContext(ctx).Warn("writing playlist response failed",
			slog.String("error", err.Error()))
	}
}
