// File: internal/biblebrain/books.go
// Purpose: Fetch the list of books for a Bible from Bible Brain.
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

// booksResponse matches Bible Brain's JSON envelope for the books endpoint.
type booksResponse struct {
	Data []domain.Book `json:"data"`
}

// FetchBooks calls GET /bibles/{bibleID}/book on the Bible Brain API.
// Returns the books array as JSON bytes suitable for caching and proxying.
func (c *Client) FetchBooks(ctx context.Context, bibleID string) ([]byte, error) {
	params := url.Values{}
	params.Set("verify_content", "true")

	raw, err := c.get(ctx, fmt.Sprintf("/bibles/%s/book", bibleID), params)
	if err != nil {
		return nil, fmt.Errorf("fetching books for bible %s: %w", bibleID, err)
	}

	var parsed booksResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, fmt.Errorf("parsing books json for %s: %w", bibleID, err)
	}

	return json.Marshal(parsed.Data)
}
