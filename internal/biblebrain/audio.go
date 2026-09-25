// File: internal/biblebrain/audio.go
// Purpose: Fetch chapter audio content from Bible Brain filesets.
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

// audioResponse matches Bible Brain's chapter-audio JSON envelope.
type audioResponse struct {
	Data []domain.AudioChapter `json:"data"`
}

// FetchAudio returns the audio file URL(s) for one chapter of a fileset.
// Bible Brain endpoint: GET /bibles/filesets/{fileset_id}/{book}/{chapter}
func (c *Client) FetchAudio(ctx context.Context, filesetID, book string, chapter int) ([]byte, error) {
	path := fmt.Sprintf("/bibles/filesets/%s/%s/%d", filesetID, book, chapter)

	raw, err := c.get(ctx, path, url.Values{})
	if err != nil {
		return nil, fmt.Errorf("fetching audio (fileset=%s, book=%s, ch=%d): %w",
			filesetID, book, chapter, err)
	}

	var parsed audioResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, fmt.Errorf("parsing audio json (fileset=%s): %w", filesetID, err)
	}

	return json.Marshal(parsed.Data)
}