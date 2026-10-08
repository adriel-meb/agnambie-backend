// File: internal/biblebrain/playlist.go
// Purpose: Fetch key-protected HLS playlists server-side and rewrite them so
//          clients can play them without ever seeing the Bible Brain API key.
// Author: Backend Team
// Created: 2026-10-08
// Last Modified: 2026-10-08

package biblebrain

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// ErrNoPlaylist is returned when a chapter has no key-protected HLS playlist,
// e.g. because the fileset is served as plain signed MP3 files.
var ErrNoPlaylist = errors.New("no hls playlist for this chapter")

// uriAttrPattern matches URI="..." attributes in HLS tags such as
// #EXT-X-KEY and #EXT-X-MAP.
var uriAttrPattern = regexp.MustCompile(`URI="([^"]*)"`)

// FetchPlaylist returns a client-safe HLS playlist for one chapter.
//
// It looks up the chapter's upstream playlist URL, downloads it with the API
// key, and rewrites every URI inside it to an absolute URL with any "key"
// query parameter removed. Segment URLs are signed CDN links, so clients
// fetch audio straight from the CDN and the backend carries no audio traffic.
func (c *Client) FetchPlaylist(ctx context.Context, filesetID, book string, chapter int) ([]byte, error) {
	items, err := c.fetchAudioItems(ctx, filesetID, book, chapter)
	if err != nil {
		return nil, err
	}

	var upstream string
	for _, item := range items {
		if c.IsAPIHost(item.Path) {
			upstream = item.Path
			break
		}
	}
	if upstream == "" {
		return nil, ErrNoPlaylist
	}

	base, err := url.Parse(upstream)
	if err != nil {
		return nil, fmt.Errorf("parsing upstream playlist url: %w", err)
	}

	raw, err := c.getAbsolute(ctx, upstream)
	if err != nil {
		return nil, fmt.Errorf("fetching playlist (fileset=%s, book=%s, ch=%d): %w",
			filesetID, book, chapter, err)
	}

	out := RewritePlaylist(raw, base)

	// Final safety net: never return anything containing the key, even if the
	// upstream format changes in a way the rewriter doesn't anticipate.
	if c.apiKey != "" && bytes.Contains(out, []byte(c.apiKey)) {
		return nil, fmt.Errorf("playlist for fileset=%s still references the api key", filesetID)
	}
	return out, nil
}

// RewritePlaylist makes every URI in an HLS playlist absolute (resolved against
// base) and strips any "key" query parameter. URIs without a key parameter keep
// their original query string byte-for-byte, because CloudFront signatures
// can break if query parameters are re-encoded or reordered.
func RewritePlaylist(playlist []byte, base *url.URL) []byte {
	lines := strings.Split(string(playlist), "\n")
	for i, line := range lines {
		trimmed := strings.TrimRight(line, "\r")
		switch {
		case trimmed == "":
			continue
		case strings.HasPrefix(trimmed, "#"):
			// Tags may carry URIs (encryption keys, init segments).
			lines[i] = uriAttrPattern.ReplaceAllStringFunc(trimmed, func(m string) string {
				inner := uriAttrPattern.FindStringSubmatch(m)[1]
				return `URI="` + rewriteURI(inner, base) + `"`
			})
		default:
			lines[i] = rewriteURI(strings.TrimSpace(trimmed), base)
		}
	}
	return []byte(strings.Join(lines, "\n"))
}

// rewriteURI resolves ref against base and removes any "key" query parameter.
// Unparseable URIs are returned with only the key parameter removed.
func rewriteURI(ref string, base *url.URL) string {
	u, err := url.Parse(ref)
	if err != nil {
		return ref
	}
	resolved := base.ResolveReference(u)
	resolved.RawQuery = removeQueryParam(resolved.RawQuery, "key")
	return resolved.String()
}

// removeQueryParam drops every occurrence of name from rawQuery while keeping
// the remaining parameters in their original order and encoding.
func removeQueryParam(rawQuery, name string) string {
	if rawQuery == "" {
		return rawQuery
	}
	parts := strings.Split(rawQuery, "&")
	kept := parts[:0]
	for _, p := range parts {
		k, _, _ := strings.Cut(p, "=")
		if decoded, err := url.QueryUnescape(k); err == nil && strings.EqualFold(decoded, name) {
			continue
		}
		kept = append(kept, p)
	}
	return strings.Join(kept, "&")
}
