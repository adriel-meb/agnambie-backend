// File: internal/biblebrain/bibles.go
// Purpose: Fetch Bible listings and copyright from the Bible Brain API.
// Author: Backend Team
// Created: 2026-09-25
// Last Modified: 2026-09-25

package biblebrain

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/adriel-meb/agnambie-backend/internal/domain"
)

// biblesListResponse matches Bible Brain's paginated Bible listing envelope.
type biblesListResponse struct {
	Data []domain.BibleSummary `json:"data"`
}

// FetchBibles calls GET /bibles?language_code={code} on the Bible Brain API.
// Returns the bibles array as JSON bytes suitable for caching and proxying.
func (c *Client) FetchBibles(ctx context.Context, languageCode string) ([]byte, error) {
	params := url.Values{}
	params.Set("language_code", languageCode)
	params.Set("page", "1")
	params.Set("limit", "50")

	raw, err := c.get(ctx, "/bibles", params)
	if err != nil {
		return nil, fmt.Errorf("fetching bibles for language %s: %w", languageCode, err)
	}

	var parsed biblesListResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, fmt.Errorf("parsing bibles json for %s: %w", languageCode, err)
	}

	return json.Marshal(parsed.Data)
}

// FetchCopyright calls GET /bibles/{bibleID}/copyright on the Bible Brain API.
// Returns the raw copyright JSON for proxying to the Flutter client.
func (c *Client) FetchCopyright(ctx context.Context, bibleID string) ([]byte, error) {
	raw, err := c.get(ctx, fmt.Sprintf("/bibles/%s/copyright", bibleID), url.Values{})
	if err != nil {
		return nil, fmt.Errorf("fetching copyright for bible %s: %w", bibleID, err)
	}

	return raw, nil
}
