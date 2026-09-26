// File: cmd/server/main.go
// Purpose: Application entrypoint — wires dependencies and starts the server.
// Author: Backend Team
// Created: 2026-09-25
// Last Modified: 2026-09-25

package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/adriel-meb/agnambie-backend/internal/biblebrain"
	"github.com/adriel-meb/agnambie-backend/internal/cache"
	"github.com/adriel-meb/agnambie-backend/internal/config"
	"github.com/adriel-meb/agnambie-backend/internal/handler"
	"github.com/adriel-meb/agnambie-backend/internal/router"
)

func main() {
	// ── 1. Load environment ─────────────────────────────────────────────
	if err := godotenv.Load(); err != nil {
		slog.Warn("no .env file found, reading from environment")
	}

	cfg := config.Load()

	// ── 2. Structured logging ───────────────────────────────────────────
	logger := setupLogger(cfg)
	slog.SetDefault(logger)

	// ── 3. Validate required config ─────────────────────────────────────
	if cfg.BibleBrainAPIKey == "" {
		logger.Error("BIBLE_BRAIN_API_KEY is required but not set")
		os.Exit(1)
	}

	// ── 4. Cache ────────────────────────────────────────────────────────
	c, err := cache.New(cfg.RedisURL)
	if err != nil {
		logger.Error("failed to create cache", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer c.Close()

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()
	if err := c.Ping(pingCtx); err != nil {
		logger.Error("cannot reach cache backend", slog.String("error", err.Error()))
		os.Exit(1)
	}

	cacheType := "memory"
	if cfg.RedisURL != "" {
		cacheType = "redis"
	}
	logger.Info("cache connected", slog.String("type", cacheType))

	// ── 5. Bible Brain client ───────────────────────────────────────────
	bbClient := biblebrain.NewClient(
		cfg.BibleBrainBaseURL,
		cfg.BibleBrainAPIKey,
		cfg.BibleBrainTimeout,
	)

	// ── 6. Handler + Router ─────────────────────────────────────────────
	h := handler.NewHandler(bbClient, c)
	r := router.New(h, cfg.CORSAllowedOrigins, cfg.RateLimitReqsPerSec, cfg.RateLimitBurst, logger)

	// ── 7. Start server with graceful shutdown ──────────────────────────
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("server starting",
			slog.String("port", cfg.Port),
			slog.String("env", cfg.Env),
		)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	sig := <-shutdownCh
	logger.Info("shutting down", slog.String("signal", sig.String()))

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("forced shutdown", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("server stopped")
}

// setupLogger creates a structured logger depending on the configuration.
func setupLogger(cfg config.Config) *slog.Logger {
	var logLevel slog.Level
	switch strings.ToLower(cfg.LogLevel) {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	var handler slog.Handler
	if strings.ToLower(cfg.LogFormat) == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}
