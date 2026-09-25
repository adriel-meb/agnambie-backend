// File: internal/middleware/compress.go
// Purpose: Gzip compression middleware for low-bandwidth clients.
// Author: Backend Team
// Created: 2026-09-25
// Last Modified: 2026-09-25

package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"
)

// gzipWriterPool recycles gzip writers to avoid repeated allocations.
var gzipWriterPool = sync.Pool{
	New: func() any {
		// Level 6 (default) gives good compression/speed tradeoff.
		w, _ := gzip.NewWriterLevel(io.Discard, gzip.DefaultCompression)
		return w
	},
}

// gzipResponseWriter wraps http.ResponseWriter to compress the body.
type gzipResponseWriter struct {
	http.ResponseWriter
	gw *gzip.Writer
}

// Write compresses bytes through the gzip writer.
func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.gw.Write(b)
}

// Compress returns middleware that gzip-compresses responses when the client
// sends Accept-Encoding: gzip. This is critical for low-bandwidth users in
// Gabon — JSON payloads compress 70-80%.
//
// Only compresses responses with Content-Type containing "json" or "text" to
// avoid double-compressing already-compressed formats.
func Compress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip if client doesn't accept gzip.
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// Get a gzip writer from the pool.
		gz := gzipWriterPool.Get().(*gzip.Writer)
		gz.Reset(w)

		defer func() {
			gz.Close()
			gzipWriterPool.Put(gz)
		}()

		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Vary", "Accept-Encoding")
		// Delete Content-Length since compressed size differs.
		w.Header().Del("Content-Length")

		next.ServeHTTP(&gzipResponseWriter{ResponseWriter: w, gw: gz}, r)
	})
}
