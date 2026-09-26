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
	"strings"


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

	// Hotfix for HLS streams (4.dbt.io) missing key and v=4 parameters
	for i, item := range parsed.Data {
		if strings.Contains(item.Path, "4.dbt.io") && strings.Contains(item.Path, "playlist.m3u8") {
			if u, err := url.Parse(item.Path); err == nil {
				q := u.Query()
				q.Set("v", "4")
				q.Set("key", c.apiKey)
				u.RawQuery = q.Encode()
				parsed.Data[i].Path = u.String()
			}
		}
	}


	return json.Marshal(parsed.Data)
}