// File: internal/biblebrain/client_test.go
// Purpose: Tests for Bible Brain HTTP client — URL building, error handling.
// Author: Backend Team
// Created: 2026-09-25
// Last Modified: 2026-09-25

package biblebrain

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestBuildURL(t *testing.T) {
	c := NewClient("https://4.dbt.io/api", "test-key", 5*time.Second)

	tests := []struct {
		name       string
		path       string
		params     url.Values
		wantHost   string
		wantPath   string
		wantV      string
		wantKey    string
		wantExtra  map[string]string
	}{
		{
			name:     "simple path no extra params",
			path:     "/bibles",
			params:   url.Values{},
			wantHost: "4.dbt.io",
			wantPath: "/api/bibles",
			wantV:    "4",
			wantKey:  "test-key",
		},
		{
			name:   "path with extra params",
			path:   "/bibles",
			params: url.Values{"language_code": []string{"FAN"}, "page": []string{"1"}},
			wantHost: "4.dbt.io",
			wantPath: "/api/bibles",
			wantV:    "4",
			wantKey:  "test-key",
			wantExtra: map[string]string{
				"language_code": "FAN",
				"page":          "1",
			},
		},
		{
			name:     "nested path",
			path:     "/bibles/FANBSG/book",
			params:   url.Values{"verify_content": []string{"true"}},
			wantHost: "4.dbt.io",
			wantPath: "/api/bibles/FANBSG/book",
			wantV:    "4",
			wantKey:  "test-key",
			wantExtra: map[string]string{
				"verify_content": "true",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := c.buildURL(tt.path, tt.params)
			parsed, err := url.Parse(result)
			if err != nil {
				t.Fatalf("buildURL returned invalid URL: %v", err)
			}
			if parsed.Host != tt.wantHost {
				t.Errorf("host = %q, want %q", parsed.Host, tt.wantHost)
			}
			if parsed.Path != tt.wantPath {
				t.Errorf("path = %q, want %q", parsed.Path, tt.wantPath)
			}
			q := parsed.Query()
			if q.Get("v") != tt.wantV {
				t.Errorf("v = %q, want %q", q.Get("v"), tt.wantV)
			}
			if q.Get("key") != tt.wantKey {
				t.Errorf("key = %q, want %q", q.Get("key"), tt.wantKey)
			}
			for k, v := range tt.wantExtra {
				if q.Get(k) != v {
					t.Errorf("param %s = %q, want %q", k, q.Get(k), v)
				}
			}
		})
	}
}

func TestGet_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":"ok"}`))
	}))
	defer server.Close()

	c := NewClient(server.URL, "test-key", 5*time.Second)
	body, err := c.get(context.Background(), "/test", url.Values{})
	if err != nil {
		t.Fatalf("get returned error: %v", err)
	}
	if string(body) != `{"data":"ok"}` {
		t.Errorf("body = %q, want %q", string(body), `{"data":"ok"}`)
	}
}

func TestGet_Non200_ReturnsUpstreamError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"not found"}`))
	}))
	defer server.Close()

	c := NewClient(server.URL, "test-key", 5*time.Second)
	_, err := c.get(context.Background(), "/missing", url.Values{})
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}

	upErr, ok := err.(*UpstreamError)
	if !ok {
		t.Fatalf("expected *UpstreamError, got %T", err)
	}
	if upErr.StatusCode != 404 {
		t.Errorf("StatusCode = %d, want 404", upErr.StatusCode)
	}
}

func TestGet_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.Write([]byte(`{"data":"late"}`))
	}))
	defer server.Close()

	// Use a very short timeout.
	c := NewClient(server.URL, "test-key", 50*time.Millisecond)
	_, err := c.get(context.Background(), "/slow", url.Values{})
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

func TestGet_ContextCancelled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	c := NewClient(server.URL, "test-key", 5*time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately.

	_, err := c.get(ctx, "/cancelled", url.Values{})
	if err == nil {
		t.Fatal("expected context cancelled error, got nil")
	}
}

func TestNewClient_DefaultBaseURL(t *testing.T) {
	c := NewClient("", "key", 5*time.Second)
	if c.baseURL != "https://4.dbt.io/api" {
		t.Errorf("baseURL = %q, want default", c.baseURL)
	}
}
