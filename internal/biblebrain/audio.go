// File: internal/biblebrain/audio.go
// Purpose: Fetch chapter audio content from Bible Brain filesets.
// Author: Backend Team
// Created: 2026-09-25
// Last Modified: 2026-10-08

package biblebrain

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/adriel-meb/agnambie-backend/internal/domain"
)

// PlaylistRoute is the backend route that serves key-free HLS playlists.
// FetchAudio rewrites key-requiring playlist links to this relative path so the
// Bible Brain API key never leaves the server.
const PlaylistRoute = "/api/v1/audio/playlist"

// audioResponse matches Bible Brain's chapter-audio JSON envelope.
type audioResponse struct {
	Data []domain.AudioChapter `json:"data"`
}

// FetchAudio returns the audio file URL(s) for one chapter of a fileset.
// Bible Brain endpoint: GET /bibles/filesets/{fileset_id}/{book}/{chapter}
//
// Paths that point at the Bible Brain API host (HLS playlists that need the
// API key) are replaced with a relative link to PlaylistRoute. Signed CDN
// paths are returned untouched.
func (c *Client) FetchAudio(ctx context.Context, filesetID, book string, chapter int) ([]byte, error) {
	items, err := c.fetchAudioItems(ctx, filesetID, book, chapter)
	if err != nil {
		return nil, err
	}

	for i, item := range items {
		if c.IsAPIHost(item.Path) {
			items[i].Path = playlistPath(filesetID, book, chapter)
		}
	}

	return json.Marshal(items)
}

// fetchAudioItems returns the raw, unmodified audio items from Bible Brain.
// The result may contain key-requiring URLs and must not be sent to clients.
func (c *Client) fetchAudioItems(ctx context.Context, filesetID, book string, chapter int) ([]domain.AudioChapter, error) {
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
	return parsed.Data, nil
}

// playlistPath builds the relative backend URL that serves a chapter's playlist.
func playlistPath(filesetID, book string, chapter int) string {
	q := url.Values{}
	q.Set("fileset_id", filesetID)
	q.Set("book", book)
	q.Set("chapter", strconv.Itoa(chapter))
	return PlaylistRoute + "?" + q.Encode()
}
