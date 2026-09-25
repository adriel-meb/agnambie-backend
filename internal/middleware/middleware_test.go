// File: internal/middleware/middleware_test.go
// Purpose: Unit tests for CORS, SecurityHeaders, and RequestLogger middlewares.
// Author: Backend Team
// Created: 2026-09-25
// Last Modified: 2026-09-25

package middleware

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCORS_Wildcard(t *testing.T) {
	tests := []struct {
		name           string
		allowedOrigins []string
		requestOrigin  string
	}{
		{
			name:           "wildcard only with origin header",
			allowedOrigins: []string{"*"},
			requestOrigin:  "https://example.com",
		},
		{
			name:           "wildcard only without origin header",
			allowedOrigins: []string{"*"},
			requestOrigin:  "",
		},
		{
			name:           "wildcard with multiple origins and localhost origin",
			allowedOrigins: []string{"https://app.example.com", "*"},
			requestOrigin:  "http://localhost:3000",
		},
		{
			name:           "wildcard with custom domain origin",
			allowedOrigins: []string{"*"},
			requestOrigin:  "https://subdomain.test.org:8080",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			nextCalled := false
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				w.WriteHeader(http.StatusOK)
			})

			handler := CORS(tc.allowedOrigins)(next)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tc.requestOrigin != "" {
				req.Header.Set("Origin", tc.requestOrigin)
			}
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if !nextCalled {
				t.Errorf("expected next handler to be called for GET request")
			}

			// Verify wildcard origin header is set
			gotAllowOrigin := rec.Header().Get("Access-Control-Allow-Origin")
			if gotAllowOrigin != "*" {
				t.Errorf("Access-Control-Allow-Origin = %q; want %q", gotAllowOrigin, "*")
			}

			// In wildcard mode, Vary header should not be set
			if gotVary := rec.Header().Get("Vary"); gotVary != "" {
				t.Errorf("Vary = %q; want empty for wildcard CORS", gotVary)
			}

			// Verify standard CORS headers
			if gotMethods := rec.Header().Get("Access-Control-Allow-Methods"); gotMethods != "GET, OPTIONS" {
				t.Errorf("Access-Control-Allow-Methods = %q; want %q", gotMethods, "GET, OPTIONS")
			}
			if gotHeaders := rec.Header().Get("Access-Control-Allow-Headers"); gotHeaders != "Content-Type, Accept" {
				t.Errorf("Access-Control-Allow-Headers = %q; want %q", gotHeaders, "Content-Type, Accept")
			}
			if gotMaxAge := rec.Header().Get("Access-Control-Max-Age"); gotMaxAge != "86400" {
				t.Errorf("Access-Control-Max-Age = %q; want %q", gotMaxAge, "86400")
			}
		})
	}
}

func TestCORS_SpecificOrigins(t *testing.T) {
	allowedOrigins := []string{
		"https://example.com",
		"http://localhost:3000",
		"https://trailing-slash.com/",
	}

	tests := []struct {
		name                string
		allowedOrigins      []string
		requestOrigin       string
		wantAllowOrigin     string
		wantVaryHeader      bool
		wantNextHandlerCall bool
	}{
		{
			name:                "matching origin exact",
			allowedOrigins:      allowedOrigins,
			requestOrigin:       "https://example.com",
			wantAllowOrigin:     "https://example.com",
			wantVaryHeader:      true,
			wantNextHandlerCall: true,
		},
		{
			name:                "matching localhost with port",
			allowedOrigins:      allowedOrigins,
			requestOrigin:       "http://localhost:3000",
			wantAllowOrigin:     "http://localhost:3000",
			wantVaryHeader:      true,
			wantNextHandlerCall: true,
		},
		{
			name:                "matching allowed origin defined with trailing slash",
			allowedOrigins:      allowedOrigins,
			requestOrigin:       "https://trailing-slash.com",
			wantAllowOrigin:     "https://trailing-slash.com",
			wantVaryHeader:      true,
			wantNextHandlerCall: true,
		},
		{
			name:                "non-matching origin does not set allow origin header",
			allowedOrigins:      allowedOrigins,
			requestOrigin:       "https://unauthorized.org",
			wantAllowOrigin:     "",
			wantVaryHeader:      false,
			wantNextHandlerCall: true,
		},
		{
			name:                "subdomain mismatch does not set allow origin header",
			allowedOrigins:      allowedOrigins,
			requestOrigin:       "https://sub.example.com",
			wantAllowOrigin:     "",
			wantVaryHeader:      false,
			wantNextHandlerCall: true,
		},
		{
			name:                "empty origin header does not set allow origin header",
			allowedOrigins:      allowedOrigins,
			requestOrigin:       "",
			wantAllowOrigin:     "",
			wantVaryHeader:      false,
			wantNextHandlerCall: true,
		},
		{
			name:                "empty allowed origins list does not match",
			allowedOrigins:      []string{},
			requestOrigin:       "https://example.com",
			wantAllowOrigin:     "",
			wantVaryHeader:      false,
			wantNextHandlerCall: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			nextCalled := false
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				w.WriteHeader(http.StatusOK)
			})

			handler := CORS(tc.allowedOrigins)(next)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tc.requestOrigin != "" {
				req.Header.Set("Origin", tc.requestOrigin)
			}
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if nextCalled != tc.wantNextHandlerCall {
				t.Errorf("next handler called = %v; want %v", nextCalled, tc.wantNextHandlerCall)
			}

			gotAllowOrigin := rec.Header().Get("Access-Control-Allow-Origin")
			if gotAllowOrigin != tc.wantAllowOrigin {
				t.Errorf("Access-Control-Allow-Origin = %q; want %q", gotAllowOrigin, tc.wantAllowOrigin)
			}

			gotVary := rec.Header().Get("Vary")
			if tc.wantVaryHeader && gotVary != "Origin" {
				t.Errorf("Vary = %q; want %q", gotVary, "Origin")
			} else if !tc.wantVaryHeader && gotVary != "" {
				t.Errorf("Vary = %q; want empty", gotVary)
			}

			// General CORS headers should always be present
			if gotMethods := rec.Header().Get("Access-Control-Allow-Methods"); gotMethods != "GET, OPTIONS" {
				t.Errorf("Access-Control-Allow-Methods = %q; want %q", gotMethods, "GET, OPTIONS")
			}
			if gotHeaders := rec.Header().Get("Access-Control-Allow-Headers"); gotHeaders != "Content-Type, Accept" {
				t.Errorf("Access-Control-Allow-Headers = %q; want %q", gotHeaders, "Content-Type, Accept")
			}
			if gotMaxAge := rec.Header().Get("Access-Control-Max-Age"); gotMaxAge != "86400" {
				t.Errorf("Access-Control-Max-Age = %q; want %q", gotMaxAge, "86400")
			}
		})
	}
}

func TestCORS_PreflightOptions(t *testing.T) {
	tests := []struct {
		name            string
		method          string
		allowedOrigins  []string
		requestOrigin   string
		wantStatusCode  int
		wantAllowOrigin string
		wantNextCall    bool
	}{
		{
			name:            "OPTIONS preflight with wildcard origin returns 204 without calling next",
			method:          http.MethodOptions,
			allowedOrigins:  []string{"*"},
			requestOrigin:   "https://example.com",
			wantStatusCode:  http.StatusNoContent,
			wantAllowOrigin: "*",
			wantNextCall:    false,
		},
		{
			name:            "OPTIONS preflight with matching specific origin returns 204 without calling next",
			method:          http.MethodOptions,
			allowedOrigins:  []string{"https://app.example.com"},
			requestOrigin:   "https://app.example.com",
			wantStatusCode:  http.StatusNoContent,
			wantAllowOrigin: "https://app.example.com",
			wantNextCall:    false,
		},
		{
			name:            "OPTIONS preflight with non-matching origin returns 204 without calling next",
			method:          http.MethodOptions,
			allowedOrigins:  []string{"https://app.example.com"},
			requestOrigin:   "https://evil.com",
			wantStatusCode:  http.StatusNoContent,
			wantAllowOrigin: "",
			wantNextCall:    false,
		},
		{
			name:            "GET request calls next handler and returns downstream status",
			method:          http.MethodGet,
			allowedOrigins:  []string{"https://app.example.com"},
			requestOrigin:   "https://app.example.com",
			wantStatusCode:  http.StatusOK,
			wantAllowOrigin: "https://app.example.com",
			wantNextCall:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			nextCalled := false
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				w.WriteHeader(http.StatusOK)
			})

			handler := CORS(tc.allowedOrigins)(next)

			req := httptest.NewRequest(tc.method, "/api/data", nil)
			if tc.requestOrigin != "" {
				req.Header.Set("Origin", tc.requestOrigin)
			}
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatusCode {
				t.Errorf("status code = %d; want %d", rec.Code, tc.wantStatusCode)
			}

			if nextCalled != tc.wantNextCall {
				t.Errorf("next handler called = %v; want %v", nextCalled, tc.wantNextCall)
			}

			gotAllowOrigin := rec.Header().Get("Access-Control-Allow-Origin")
			if gotAllowOrigin != tc.wantAllowOrigin {
				t.Errorf("Access-Control-Allow-Origin = %q; want %q", gotAllowOrigin, tc.wantAllowOrigin)
			}

			// Preflight responses must have no body content
			if tc.method == http.MethodOptions && rec.Body.Len() > 0 {
				t.Errorf("expected empty body for OPTIONS preflight, got %q", rec.Body.String())
			}
		})
	}
}

func TestSecurityHeaders(t *testing.T) {
	expectedHeaders := map[string]string{
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":        "DENY",
		"X-XSS-Protection":       "1; mode=block",
		"Content-Security-Policy": "default-src 'none'",
	}

	tests := []struct {
		name           string
		method         string
		path           string
		responseStatus int
		responseBody   string
	}{
		{
			name:           "GET request with 200 OK",
			method:         http.MethodGet,
			path:           "/test",
			responseStatus: http.StatusOK,
			responseBody:   "ok",
		},
		{
			name:           "POST request with 201 Created",
			method:         http.MethodPost,
			path:           "/submit",
			responseStatus: http.StatusCreated,
			responseBody:   `{"created": true}`,
		},
		{
			name:           "GET request with 404 Not Found",
			method:         http.MethodGet,
			path:           "/missing",
			responseStatus: http.StatusNotFound,
			responseBody:   "not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			nextCalled := false
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				w.WriteHeader(tc.responseStatus)
				_, _ = w.Write([]byte(tc.responseBody))
			})

			handler := SecurityHeaders(next)

			req := httptest.NewRequest(tc.method, tc.path, nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if !nextCalled {
				t.Fatalf("expected next handler to be called")
			}

			if rec.Code != tc.responseStatus {
				t.Errorf("status code = %d; want %d", rec.Code, tc.responseStatus)
			}

			if rec.Body.String() != tc.responseBody {
				t.Errorf("response body = %q; want %q", rec.Body.String(), tc.responseBody)
			}

			// Verify all 4 required security headers
			for header, expectedValue := range expectedHeaders {
				actualValue := rec.Header().Get(header)
				if actualValue != expectedValue {
					t.Errorf("header %q = %q; want %q", header, actualValue, expectedValue)
				}
			}
		})
	}
}

func TestRequestLogger(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		handlerStatus  int
		handlerBody    string
		explicitStatus bool
	}{
		{
			name:           "default 200 OK without explicit WriteHeader",
			method:         http.MethodGet,
			path:           "/api/items",
			handlerStatus:  http.StatusOK,
			handlerBody:    "item list",
			explicitStatus: false,
		},
		{
			name:           "explicit 201 Created",
			method:         http.MethodPost,
			path:           "/api/items",
			handlerStatus:  http.StatusCreated,
			handlerBody:    "created",
			explicitStatus: true,
		},
		{
			name:           "explicit 404 Not Found",
			method:         http.MethodGet,
			path:           "/api/unknown",
			handlerStatus:  http.StatusNotFound,
			handlerBody:    "not found",
			explicitStatus: true,
		},
		{
			name:           "explicit 500 Internal Server Error",
			method:         http.MethodGet,
			path:           "/api/error",
			handlerStatus:  http.StatusInternalServerError,
			handlerBody:    "internal error",
			explicitStatus: true,
		},
		{
			name:           "empty path normalized to root",
			method:         http.MethodGet,
			path:           "",
			handlerStatus:  http.StatusOK,
			handlerBody:    "root",
			explicitStatus: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var logBuf bytes.Buffer
			logger := slog.New(slog.NewTextHandler(&logBuf, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}))

			nextCalled := false
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				if tc.explicitStatus {
					w.WriteHeader(tc.handlerStatus)
				}
				_, _ = w.Write([]byte(tc.handlerBody))
			})

			handler := RequestLogger(logger)(next)

			reqPath := tc.path
			if reqPath == "" {
				reqPath = "/"
			}
			req := httptest.NewRequest(tc.method, reqPath, nil)
			req.RemoteAddr = "192.0.2.1:12345"
			rec := httptest.NewRecorder()

			// Ensure it doesn't panic
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("RequestLogger panicked: %v", r)
				}
			}()

			handler.ServeHTTP(rec, req)

			if !nextCalled {
				t.Fatalf("expected next handler to be called")
			}

			if rec.Code != tc.handlerStatus {
				t.Errorf("status code = %d; want %d", rec.Code, tc.handlerStatus)
			}

			if rec.Body.String() != tc.handlerBody {
				t.Errorf("response body = %q; want %q", rec.Body.String(), tc.handlerBody)
			}

			// Verify log output contains expected fields
			logOutput := logBuf.String()
			if !strings.Contains(logOutput, "msg=request") && !strings.Contains(logOutput, "request") {
				t.Errorf("log output missing msg field: %q", logOutput)
			}
			if !strings.Contains(logOutput, tc.method) {
				t.Errorf("log output missing method %q: %q", tc.method, logOutput)
			}
			if !strings.Contains(logOutput, reqPath) {
				t.Errorf("log output missing path %q: %q", reqPath, logOutput)
			}
			if !strings.Contains(logOutput, "192.0.2.1:12345") {
				t.Errorf("log output missing remote addr: %q", logOutput)
			}
		})
	}

	t.Run("does not panic with discarded logger output", func(t *testing.T) {
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		nextCalled := false
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nextCalled = true
			w.WriteHeader(http.StatusTeapot)
		})

		handler := RequestLogger(logger)(next)
		req := httptest.NewRequest(http.MethodDelete, "/resource/123", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if !nextCalled {
			t.Errorf("expected next handler to be called")
		}
		if rec.Code != http.StatusTeapot {
			t.Errorf("status code = %d; want %d", rec.Code, http.StatusTeapot)
		}
	})
}
