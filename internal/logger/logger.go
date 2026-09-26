// Package logger provides context-aware structured logging capabilities.
//
// Purpose: Enables passing a request-scoped logger (containing request_id) through the call stack.
// Author: Backend Team
// Created: 2026-09-25
// Last Modified: 2026-09-25

package logger

import (
	"context"
	"log/slog"
)

type contextKey string

const loggerKey = contextKey("logger")

// WithLogger returns a new context with the provided logger embedded.
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext retrieves the logger from the context.
// If no logger is present, it returns the default slog.Logger.
func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}
