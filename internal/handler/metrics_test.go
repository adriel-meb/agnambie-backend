package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMetricsHandler(t *testing.T) {
	h := NewHandler(nil, nil)

	t.Run("valid batch", func(t *testing.T) {
		events := []MetricEvent{
			{
				Name:      "chapter_played",
				Timestamp: time.Now().UTC(),
				Properties: map[string]any{
					"language_iso": "fan",
					"bible_id":     "FANBSG",
					"book_id":      "MAT",
					"chapter":      1,
					"is_offline":   true,
				},
			},
		}

		body, err := json.Marshal(events)
		if err != nil {
			t.Fatalf("failed to marshal events: %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/api/v1/metrics", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		h.Metrics(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d, body: %s", http.StatusOK, rr.Code, rr.Body.String())
		}
	})

	t.Run("empty array returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/metrics", bytes.NewReader([]byte("[]")))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		h.Metrics(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("invalid json returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/metrics", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		h.Metrics(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("method not allowed on GET", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/metrics", nil)
		rr := httptest.NewRecorder()

		h.Metrics(rr, req)

		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
		}
	})
}
