// File: internal/logger/pretty.go
// Purpose: Human-friendly ANSI colored terminal slog handler for development.
// Author: Backend Team
// Created: 2026-09-26
// Last Modified: 2026-09-26

package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

// ANSI color escape sequences.
const (
	reset      = "\033[0m"
	bold       = "\033[1m"
	dim        = "\033[2m"
	red        = "\033[31m"
	green      = "\033[32m"
	yellow     = "\033[33m"
	blue       = "\033[34m"
	magenta    = "\033[35m"
	cyan       = "\033[36m"
	gray       = "\033[90m"
	boldRed    = "\033[1;31m"
	boldGreen  = "\033[1;32m"
	boldYellow = "\033[1;33m"
	boldCyan   = "\033[1;36m"
	boldWhite  = "\033[1;37m"
)

// PrettyHandler formats log records into clean, readable, colored terminal output.
type PrettyHandler struct {
	opts  slog.HandlerOptions
	out   io.Writer
	mu    *sync.Mutex
	attrs []slog.Attr
}

// NewPrettyHandler creates a new colored terminal log handler.
func NewPrettyHandler(out io.Writer, opts *slog.HandlerOptions) *PrettyHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &PrettyHandler{
		opts: *opts,
		out:  out,
		mu:   &sync.Mutex{},
	}
}

// Enabled reports whether the handler emits logs at the given level.
func (h *PrettyHandler) Enabled(_ context.Context, level slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return level >= minLevel
}

// Handle formats and prints the log record.
func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	timeStr := r.Time.Format("15:04:05")

	// 1. Level badge
	var levelBadge string
	switch r.Level {
	case slog.LevelDebug:
		levelBadge = fmt.Sprintf("%sDEBUG%s", magenta, reset)
	case slog.LevelInfo:
		levelBadge = fmt.Sprintf("%sINFO %s", cyan, reset)
	case slog.LevelWarn:
		levelBadge = fmt.Sprintf("%sWARN %s", yellow, reset)
	case slog.LevelError:
		levelBadge = fmt.Sprintf("%sERROR%s", boldRed, reset)
	default:
		levelBadge = r.Level.String()
	}

	// 2. Extract common HTTP attributes if present
	var (
		method   string
		path     string
		status   int
		duration time.Duration
		other    []slog.Attr
	)

	// Combine handler-level and record-level attrs
	allAttrs := append([]slog.Attr(nil), h.attrs...)
	r.Attrs(func(a slog.Attr) bool {
		allAttrs = append(allAttrs, a)
		return true
	})

	for _, a := range allAttrs {
		switch a.Key {
		case "method":
			method = a.Value.String()
		case "path":
			path = a.Value.String()
		case "status":
			status = int(a.Value.Int64())
		case "duration":
			duration = a.Value.Duration()
		default:
			other = append(other, a)
		}
	}

	// 3. Render HTTP access logs specially
	if method != "" && path != "" {
		methodColor := blue
		switch method {
		case http.MethodGet:
			methodColor = cyan
		case http.MethodPost:
			methodColor = green
		case http.MethodDelete:
			methodColor = red
		case http.MethodPut, http.MethodPatch:
			methodColor = yellow
		}

		statusColor := green
		if status >= 300 && status < 400 {
			statusColor = cyan
		} else if status >= 400 && status < 500 {
			statusColor = yellow
		} else if status >= 500 {
			statusColor = boldRed
		}

		statusText := fmt.Sprintf("%s%3d %-2s%s", statusColor, status, http.StatusText(status), reset)
		durText := fmt.Sprintf("%s%6.1fms%s", magenta, float64(duration.Microseconds())/1000.0, reset)

		fmt.Fprintf(h.out, "%s%s%s │ %s │ %s%-6s%s %s%-24s%s │ %s │ %s",
			gray, timeStr, reset,
			levelBadge,
			methodColor, method, reset,
			boldWhite, path, reset,
			statusText,
			durText,
		)
	} else {
		// General application log line
		fmt.Fprintf(h.out, "%s%s%s │ %s │ %s%s%s",
			gray, timeStr, reset,
			levelBadge,
			boldWhite, r.Message, reset,
		)
	}

	// 4. Render trailing key-value attributes
	if len(other) > 0 {
		fmt.Fprintf(h.out, " %s", dim)
		for _, a := range other {
			fmt.Fprintf(h.out, "%s=%v ", a.Key, a.Value.Any())
		}
		fmt.Fprintf(h.out, "%s", reset)
	}

	fmt.Fprintln(h.out)
	return nil
}

// WithAttrs returns a new handler with the given attributes added.
func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, len(h.attrs)+len(attrs))
	copy(newAttrs, h.attrs)
	copy(newAttrs[len(h.attrs):], attrs)
	return &PrettyHandler{
		opts:  h.opts,
		out:   h.out,
		mu:    h.mu,
		attrs: newAttrs,
	}
}

// WithGroup returns a new handler with the given group name.
func (h *PrettyHandler) WithGroup(_ string) slog.Handler {
	return h
}
