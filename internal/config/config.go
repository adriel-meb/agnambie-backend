// Package config provides centralized, typed configuration loaded from environment variables.
//
// Purpose: Single source of truth for all application configuration.
// Author: Backend Team
// Created: 2026-09-25
// Last Modified: 2026-09-25

package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all application configuration values, loaded from the environment.
type Config struct {
	// Port is the HTTP server listen port (default: "8080").
	Port string

	// BibleBrainBaseURL is the base URL for the Bible Brain API (default: "https://4.dbt.io/api").
	BibleBrainBaseURL string

	// BibleBrainAPIKey is the API key for authenticating with Bible Brain.
	BibleBrainAPIKey string

	// BibleBrainTimeout is the HTTP client timeout for Bible Brain requests.
	BibleBrainTimeout time.Duration

	// RedisURL is the Redis connection URL. If empty, an in-memory cache is used.
	RedisURL string

	// CORSAllowedOrigins is the list of allowed CORS origins.
	CORSAllowedOrigins []string

	// LogLevel controls the minimum log level: "debug", "info", "warn", "error".
	LogLevel string

	// LogFormat controls the structured log format: "text" (for dev) or "json" (for prod).
	LogFormat string

	// Env is the deployment environment: "development" or "production".
	Env string

	// RateLimitReqsPerSec is the allowed requests per second per IP (default: 10).
	RateLimitReqsPerSec float64

	// RateLimitBurst is the maximum burst size for rate limiting per IP (default: 30).
	RateLimitBurst int
}

// Load reads configuration from environment variables with sensible defaults.
func Load() Config {
	return Config{
		Port:                envOrDefault("PORT", "8080"),
		BibleBrainBaseURL:   envOrDefault("BIBLE_BRAIN_BASE_URL", "https://4.dbt.io/api"),
		BibleBrainAPIKey:    os.Getenv("BIBLE_BRAIN_API_KEY"),
		BibleBrainTimeout:   15 * time.Second,
		RedisURL:            os.Getenv("REDIS_URL"),
		CORSAllowedOrigins:  parseCORSOrigins(os.Getenv("CORS_ALLOWED_ORIGINS")),
		LogLevel:            envOrDefault("LOG_LEVEL", "info"),
		LogFormat:           envOrDefault("LOG_FORMAT", "text"),
		Env:                 envOrDefault("ENV", "development"),
		RateLimitReqsPerSec: envOrDefaultFloat64("RATE_LIMIT_REQS_PER_SEC", 10.0),
		RateLimitBurst:      envOrDefaultInt("RATE_LIMIT_BURST", 30),
	}
}

// IsDevelopment returns true when running in development mode.
func (c Config) IsDevelopment() bool {
	return c.Env == "development"
}

// envOrDefault returns the environment variable value or the fallback if unset/empty.
func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// parseCORSOrigins splits a comma-separated origins string into a slice.
func parseCORSOrigins(raw string) []string {
	if raw == "" {
		return []string{"*"}
	}
	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			origins = append(origins, trimmed)
		}
	}
	if len(origins) == 0 {
		return []string{"*"}
	}
	return origins
}

func envOrDefaultFloat64(key string, fallback float64) float64 {
	if v := os.Getenv(key); v != "" {
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			return parsed
		}
	}
	return fallback
}

func envOrDefaultInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			return parsed
		}
	}
	return fallback
}
