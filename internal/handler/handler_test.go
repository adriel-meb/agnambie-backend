// File: internal/handler/handler_test.go
// Purpose: Integration tests for all HTTP handlers, mocking the Bible Brain client.
// Author: Backend Team
// Created: 2026-09-25
// Last Modified: 2026-09-25

package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/adriel-meb/agnambie-backend/internal/biblebrain"
	"github.com/adriel-meb/agnambie-backend/internal/cache"
)

// setupTestEnv creates a Handler wired to a fake Bible Brain server.
func setupTestEnv(t *testing.T) (*Handler, func()) {
	t.Helper()

	fakeBB := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.HasPrefix(r.URL.Path, "/bibles/") && strings.HasSuffix(r.URL.Path, "/book"):
			w.Write([]byte(`{"data":[{"book_id":"MAT","name":"Matthew","testament":"NT","chapters":[1,2,3]}]}`))
		case strings.HasPrefix(r.URL.Path, "/bibles/filesets/"):
			w.Write([]byte(`{"data":[{"book_id":"MAT","chapter_start":1,"path":"https://cdn.example.com/audio.mp3","duration":120.5}]}`))
		case r.URL.Path == "/bibles" && r.URL.Query().Get("language_code") != "":
			w.Write([]byte(`{"data":[{"id":"FANBSG","name":"Fang Bible","iso":"fan"}]}`))
		case strings.HasSuffix(r.URL.Path, "/copyright"):
			w.Write([]byte(`[{"id":"FANBSGN2DA","type":"audio_drama","copyright":{"copyright":"© Test"}}]`))
		default:
			http.NotFound(w, r)
		}
	}))

	bbClient := biblebrain.NewClient(fakeBB.URL, "test-key", 5*time.Second)
	h := NewHandler(bbClient, cache.NewMemory())
	return h, fakeBB.Close
}

func TestLanguages(t *testing.T) {
	h, cleanup := setupTestEnv(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/api/languages", nil)
	rr := httptest.NewRecorder()
	h.Languages(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}

	var env Envelope
	if err := json.Unmarshal(rr.Body.Bytes(), &env); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if env.Meta == nil || env.Meta.Count == 0 {
		t.Errorf("missing or empty count in meta")
	}
}

func TestBibles(t *testing.T) {
	h, cleanup := setupTestEnv(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/api/bibles?language_code=FAN", nil)
	rr := httptest.NewRecorder()
	h.Bibles(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}

	var env Envelope
	if err := json.Unmarshal(rr.Body.Bytes(), &env); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if env.Meta == nil || env.Meta.Count == 0 {
		t.Errorf("expected count in meta")
	}

	// Test caching and 304 Not Modified
	etag := rr.Header().Get("ETag")
	if etag == "" {
		t.Fatal("expected ETag header")
	}

	req2 := httptest.NewRequest(http.MethodGet, "/api/bibles?language_code=FAN", nil)
	req2.Header.Set("If-None-Match", etag)
	rr2 := httptest.NewRecorder()
	h.Bibles(rr2, req2)

	if rr2.Code != http.StatusNotModified {
		t.Errorf("status = %d, want 304", rr2.Code)
	}
}

func TestBooks(t *testing.T) {
	h, cleanup := setupTestEnv(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/api/books?bible_id=FANBSG", nil)
	rr := httptest.NewRecorder()
	h.Books(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}

	var env Envelope
	if err := json.Unmarshal(rr.Body.Bytes(), &env); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if env.Meta == nil || env.Meta.Count != 1 {
		t.Errorf("expected count 1, got %v", env.Meta)
	}
}

func TestAudio(t *testing.T) {
	h, cleanup := setupTestEnv(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/api/audio?fileset_id=FANBSGN2DA&book=MAT&chapter=1", nil)
	rr := httptest.NewRecorder()
	h.Audio(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}

	var env Envelope
	if err := json.Unmarshal(rr.Body.Bytes(), &env); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if env.Meta == nil || env.Meta.Count != 1 {
		t.Errorf("expected count 1, got %v", env.Meta)
	}
}

func TestAudio_SaveDataHeader(t *testing.T) {
	h, cleanup := setupTestEnv(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/api/audio?fileset_id=FANBSGN2DA&book=MAT&chapter=1", nil)
	req.Header.Set("Save-Data", "on")
	rr := httptest.NewRecorder()
	h.Audio(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}
}

func TestCopyright(t *testing.T) {
	h, cleanup := setupTestEnv(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/api/copyright?bible_id=FANBSG", nil)
	rr := httptest.NewRecorder()
	h.Copyright(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}

	var env Envelope
	if err := json.Unmarshal(rr.Body.Bytes(), &env); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
}

func TestErrorResponse(t *testing.T) {
	h, cleanup := setupTestEnv(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/api/bibles", nil) // Missing lang code
	rr := httptest.NewRecorder()
	h.Bibles(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rr.Code)
	}

	var errResp ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if errResp.Error.Code != 400 || errResp.Error.Message == "" {
		t.Errorf("invalid error detail: %v", errResp)
	}
}
