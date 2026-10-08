// File: internal/biblebrain/playlist_test.go
// Purpose: Table-driven tests for HLS playlist rewriting and key scrubbing.
// Author: Backend Team
// Created: 2026-10-08
// Last Modified: 2026-10-08

package biblebrain

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

// signedSegment mimics a real CloudFront-signed segment URL. Its query string
// must survive rewriting byte-for-byte or the signature becomes invalid.
const signedSegment = "https://cdn.example.net/audio/X/X_B01_MAT_001.mp3?x-amz-transaction=1&Expires=99&Signature=ab~c-d_&Key-Pair-Id=APK"

func TestRewritePlaylist(t *testing.T) {
	base, _ := url.Parse("https://api.example.com/api/bible/filesets/X/MAT-1--/playlist.m3u8?v=4&key=SECRET")

	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "absolute signed segment is left untouched",
			in:   "#EXTM3U\n#EXTINF:164.000,\n" + signedSegment + "\n#EXT-X-ENDLIST\n",
			want: "#EXTM3U\n#EXTINF:164.000,\n" + signedSegment + "\n#EXT-X-ENDLIST\n",
		},
		{
			name: "relative segment is resolved against the playlist url",
			in:   "#EXTINF:10,\nseg1.mp3\n",
			want: "#EXTINF:10,\nhttps://api.example.com/api/bible/filesets/X/MAT-1--/seg1.mp3\n",
		},
		{
			name: "key param is stripped and other params keep their order",
			in:   "seg.mp3?b=2&key=SECRET&a=1\n",
			want: "https://api.example.com/api/bible/filesets/X/MAT-1--/seg.mp3?b=2&a=1\n",
		},
		{
			name: "URI attributes in tags are rewritten",
			in:   `#EXT-X-MAP:URI="init.mp4?key=SECRET"` + "\n",
			want: `#EXT-X-MAP:URI="https://api.example.com/api/bible/filesets/X/MAT-1--/init.mp4"` + "\n",
		},
		{
			name: "CRLF line endings are normalised",
			in:   "#EXTM3U\r\nseg.mp3\r\n",
			want: "#EXTM3U\nhttps://api.example.com/api/bible/filesets/X/MAT-1--/seg.mp3\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(RewritePlaylist([]byte(tt.in), base))
			if got != tt.want {
				t.Errorf("RewritePlaylist()\n got: %q\nwant: %q", got, tt.want)
			}
			if strings.Contains(got, "SECRET") {
				t.Errorf("output still contains the key: %q", got)
			}
		})
	}
}

func TestRemoveQueryParam(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"", ""},
		{"key=x", ""},
		{"a=1&key=x&b=2", "a=1&b=2"},
		{"KEY=x&a=1", "a=1"},
		{"a=1&b=2", "a=1&b=2"},
		{"monkey=1", "monkey=1"},
	}
	for _, tt := range tests {
		if got := removeQueryParam(tt.in, "key"); got != tt.want {
			t.Errorf("removeQueryParam(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

// newHLSFakeServer returns a fake Bible Brain API whose audio endpoint points
// at a key-protected playlist on the same host.
func newHLSFakeServer(t *testing.T, playlistBody string) *httptest.Server {
	t.Helper()
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/playlist.m3u8"):
			if r.URL.Query().Get("key") != "test-key" {
				http.Error(w, "missing key", http.StatusUnprocessableEntity)
				return
			}
			w.Write([]byte(playlistBody))
		case strings.HasPrefix(r.URL.Path, "/bibles/filesets/"):
			w.Write([]byte(`{"data":[{"book_id":"MAT","path":"` + srv.URL + `/bible/filesets/BNGCIEP1DA/MAT-1--/playlist.m3u8"}]}`))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestFetchAudio_RewritesKeyedPlaylistToBackendRoute(t *testing.T) {
	srv := newHLSFakeServer(t, "#EXTM3U\n")
	c := NewClient(srv.URL, "test-key", 5*time.Second)

	out, err := c.FetchAudio(context.Background(), "BNGCIEP1DA", "MAT", 1)
	if err != nil {
		t.Fatalf("FetchAudio: %v", err)
	}
	body := string(out)
	if strings.Contains(body, "test-key") {
		t.Fatalf("FetchAudio leaked the api key: %s", body)
	}
	if !strings.Contains(body, PlaylistRoute+"?book=MAT\\u0026chapter=1\\u0026fileset_id=BNGCIEP1DA") {
		t.Errorf("path not rewritten to playlist route: %s", body)
	}
}

func TestFetchPlaylist(t *testing.T) {
	srv := newHLSFakeServer(t, "#EXTM3U\n#EXTINF:164.000,\n"+signedSegment+"\n#EXT-X-ENDLIST\n")
	c := NewClient(srv.URL, "test-key", 5*time.Second)

	out, err := c.FetchPlaylist(context.Background(), "BNGCIEP1DA", "MAT", 1)
	if err != nil {
		t.Fatalf("FetchPlaylist: %v", err)
	}
	if strings.Contains(string(out), "test-key") {
		t.Fatalf("playlist leaked the api key: %s", out)
	}
	if !strings.Contains(string(out), signedSegment) {
		t.Errorf("signed segment url was altered: %s", out)
	}
}

func TestFetchPlaylist_RefusesOutputContainingKey(t *testing.T) {
	// An upstream that echoes the key in an unexpected place must be rejected.
	srv := newHLSFakeServer(t, "#EXTM3U\n#EXT-X-SESSION-DATA:VALUE=\"test-key\"\n")
	c := NewClient(srv.URL, "test-key", 5*time.Second)

	if _, err := c.FetchPlaylist(context.Background(), "BNGCIEP1DA", "MAT", 1); err == nil {
		t.Fatal("expected an error when the rewritten playlist still contains the key")
	}
}
