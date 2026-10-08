// File: internal/handler/playlist_test.go
// Purpose: Handler tests ensuring /audio and /audio/playlist never expose the
//          Bible Brain API key and enforce validation and the allowlist.
// Author: Backend Team
// Created: 2026-10-08
// Last Modified: 2026-10-08

package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/adriel-meb/agnambie-backend/internal/biblebrain"
	"github.com/adriel-meb/agnambie-backend/internal/cache"
)

const testSegment = "https://cdn.example.net/a.mp3?Expires=1&Signature=s~g&Key-Pair-Id=K"

// setupHLSEnv creates a Handler wired to a fake Bible Brain server that serves
// a key-protected HLS playlist, mirroring the real 4.dbt.io behaviour.
func setupHLSEnv(t *testing.T) *Handler {
	t.Helper()
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/playlist.m3u8"):
			if r.URL.Query().Get("key") != "test-key" {
				http.Error(w, "missing key", http.StatusUnprocessableEntity)
				return
			}
			w.Write([]byte("#EXTM3U\n#EXTINF:10.0,\n" + testSegment + "\n#EXT-X-ENDLIST\n"))
		case strings.HasPrefix(r.URL.Path, "/bibles/filesets/"):
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"data":[{"book_id":"MAT","path":"` + srv.URL + `/bible/filesets/BNGCIEP1DA/MAT-1--/playlist.m3u8"}]}`))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)

	return NewHandler(biblebrain.NewClient(srv.URL, "test-key", 5*time.Second), cache.NewMemory())
}

func TestAudio_NeverReturnsAPIKey(t *testing.T) {
	h := setupHLSEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/api/audio?fileset_id=BNGCIEP1DA&book=MAT&chapter=1", nil)
	rr := httptest.NewRecorder()
	h.Audio(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rr.Code, rr.Body)
	}
	body := rr.Body.String()
	if strings.Contains(body, "test-key") {
		t.Fatalf("response leaked the api key: %s", body)
	}
	if !strings.Contains(body, biblebrain.PlaylistRoute) {
		t.Errorf("expected path to point at %s, got %s", biblebrain.PlaylistRoute, body)
	}
}

func TestPlaylist(t *testing.T) {
	h := setupHLSEnv(t)

	tests := []struct {
		name       string
		query      string
		wantStatus int
	}{
		{"valid request", "fileset_id=BNGCIEP1DA&book=MAT&chapter=1", http.StatusOK},
		{"missing params", "fileset_id=BNGCIEP1DA", http.StatusBadRequest},
		{"bad book id", "fileset_id=BNGCIEP1DA&book=matthew!&chapter=1", http.StatusBadRequest},
		{"bad chapter", "fileset_id=BNGCIEP1DA&book=MAT&chapter=0", http.StatusBadRequest},
		{"fileset outside allowlist", "fileset_id=ENGESVN1SA&book=MAT&chapter=1", http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/audio/playlist?"+tt.query, nil)
			rr := httptest.NewRecorder()
			h.Playlist(rr, req)

			if rr.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", rr.Code, tt.wantStatus, rr.Body)
			}
			if strings.Contains(rr.Body.String(), "test-key") {
				t.Fatalf("response leaked the api key: %s", rr.Body)
			}
			if tt.wantStatus != http.StatusOK {
				return
			}
			if ct := rr.Header().Get("Content-Type"); ct != "application/vnd.apple.mpegurl" {
				t.Errorf("Content-Type = %q", ct)
			}
			if !strings.Contains(rr.Body.String(), testSegment) {
				t.Errorf("signed segment missing or altered: %s", rr.Body)
			}
		})
	}
}

func TestPlaylist_NoHLSForMP3Fileset(t *testing.T) {
	// The default fake server returns a plain CDN MP3 path, so there is no
	// key-protected playlist to serve.
	h, cleanup := setupTestEnv(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/api/audio/playlist?fileset_id=FANBSGN2DA&book=MAT&chapter=1", nil)
	rr := httptest.NewRecorder()
	h.Playlist(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", rr.Code)
	}
}
