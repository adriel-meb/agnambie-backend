// File: internal/middleware/compress_test.go
// Purpose: Unit tests for the gzip compression middleware.
// Author: Backend Team
// Created: 2026-09-25
// Last Modified: 2026-09-25

package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCompress(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"test":"compression_works_fine_and_saves_bytes"}`))
	})

	handler := Compress(nextHandler)

	t.Run("compresses when Accept-Encoding gzip is present", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Header().Get("Content-Encoding") != "gzip" {
			t.Errorf("expected Content-Encoding: gzip, got %q", rr.Header().Get("Content-Encoding"))
		}
		if rr.Header().Get("Vary") != "Accept-Encoding" {
			t.Errorf("expected Vary: Accept-Encoding, got %q", rr.Header().Get("Vary"))
		}

		// Verify we can decompress it
		gr, err := gzip.NewReader(rr.Body)
		if err != nil {
			t.Fatalf("failed to create gzip reader: %v", err)
		}
		defer gr.Close()

		decompressed, err := io.ReadAll(gr)
		if err != nil {
			t.Fatalf("failed to read decompressed body: %v", err)
		}

		expected := `{"test":"compression_works_fine_and_saves_bytes"}`
		if string(decompressed) != expected {
			t.Errorf("got %q, want %q", string(decompressed), expected)
		}
	})

	t.Run("does not compress without Accept-Encoding gzip", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Header().Get("Content-Encoding") == "gzip" {
			t.Error("unexpected Content-Encoding: gzip")
		}

		expected := `{"test":"compression_works_fine_and_saves_bytes"}`
		if !bytes.Equal(rr.Body.Bytes(), []byte(expected)) {
			t.Errorf("got %q, want %q", rr.Body.String(), expected)
		}
	})
}
