// Package biblebrain provides HTTP client interactions with the Bible Brain API.
//
// File: internal/biblebrain/client.go
// Purpose: Bible Brain API client — builds URLs, makes authenticated requests.
// Author: Backend Team
// Created: 2026-09-25
// Last Modified: 2026-09-25

package biblebrain

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/adriel-meb/agnambie-backend/internal/logger"
)

// Client is a configured HTTP client for the Bible Brain API.
// All requests automatically include the API key and version parameter.
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a Client with the given Bible Brain base URL, API key, and timeout.
func NewClient(baseURL, apiKey string, timeout time.Duration) *Client {
	if baseURL == "" {
		baseURL = "https://4.dbt.io/api"
	}
	return &Client{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: timeout},
	}
}

// buildURL constructs a full Bible Brain API URL, appending v=4 and the API key.
func (c *Client) buildURL(path string, params url.Values) string {
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		// Fallback — should not happen with well-formed paths.
		return c.baseURL + path
	}

	q := u.Query()
	q.Set("v", "4")
	q.Set("key", c.apiKey)
	for k, vv := range params {
		for _, v := range vv {
			q.Set(k, v)
		}
	}
	u.RawQuery = q.Encode()
	return u.String()
}

// scrubError removes the API key from a URL error string.
func (c *Client) scrubError(err error) error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	if c.apiKey != "" {
		msg = strings.ReplaceAll(msg, c.apiKey, "REDACTED")
	}
	return errors.New(msg)
}

// get makes an authenticated GET request to the Bible Brain API.
// It returns the raw response body or an UpstreamError on non-200 status.
func (c *Client) get(ctx context.Context, path string, params url.Values) ([]byte, error) {
	endpoint := c.buildURL(path, params)

	log := logger.FromContext(ctx).With(
		slog.String("component", "biblebrain"),
		slog.String("upstream_path", path),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("building request for %s: %w", path, c.scrubError(err))
	}

	start := time.Now()
	resp, err := c.httpClient.Do(req)
	duration := time.Since(start)

	if err != nil {
		scrubbed := c.scrubError(err)
		log.Error("bible brain request failed",
			slog.Duration("duration", duration),
			slog.String("error", scrubbed.Error()),
		)
		return nil, fmt.Errorf("calling bible brain %s: %w", path, scrubbed)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response from %s: %w", path, err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Warn("bible brain returned error status",
			slog.Int("upstream_status", resp.StatusCode),
			slog.Duration("duration", duration),
		)

		truncated := string(body)
		if len(truncated) > 500 {
			truncated = truncated[:500] + "..."
		}
		return nil, &UpstreamError{
			StatusCode: resp.StatusCode,
			Endpoint:   path,
			Body:       truncated,
		}
	}

	log.Debug("bible brain request completed",
		slog.Int("upstream_status", resp.StatusCode),
		slog.Duration("duration", duration),
	)

	return body, nil
}
