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
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
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

// get makes an authenticated GET request to the Bible Brain API.
// It returns the raw response body or an UpstreamError on non-200 status.
func (c *Client) get(ctx context.Context, path string, params url.Values) ([]byte, error) {
	endpoint := c.buildURL(path, params)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("building request for %s: %w", path, err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calling bible brain %s: %w", path, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response from %s: %w", path, err)
	}

	if resp.StatusCode != http.StatusOK {
		// Truncate body for logging to avoid storing large error pages.
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

	return body, nil
}
