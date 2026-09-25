// Package cache provides a caching abstraction with Redis and in-memory implementations.
//
// Purpose: Cache interface definition and factory constructor enabling Redis/memory swap.
// Author: Backend Team
// Created: 2026-09-25
// Last Modified: 2026-09-25

package cache

import (
	"context"
	"time"
)

// Cacher defines the caching contract used by API handlers.
// Implementations include Redis (for production) and in-memory (for local dev).
type Cacher interface {
	// Get returns cached bytes for a key. hit=false on a cache miss.
	Get(ctx context.Context, key string) (data []byte, hit bool, err error)

	// Set stores bytes under a key with a time-to-live.
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error

	// Del removes a key from the cache.
	Del(ctx context.Context, key string) error

	// GetOrSet returns cached data if present (hit=true). On a miss it calls
	// fetch(), stores the result with ttl, and returns it (hit=false).
	GetOrSet(ctx context.Context, key string, ttl time.Duration, fetch func() ([]byte, error)) (data []byte, hit bool, err error)

	// Ping checks that the cache backend is reachable.
	Ping(ctx context.Context) error

	// Close releases any resources held by the cache.
	Close() error
}

// New creates the appropriate Cacher implementation based on the provided URL.
// If redisURL is empty, an in-memory cache is returned (suitable for local dev).
// Otherwise a Redis-backed cache is created using the given connection URL.
func New(redisURL string) (Cacher, error) {
	if redisURL == "" {
		return NewMemory(), nil
	}
	r, err := NewRedis(redisURL)
	if err != nil {
		return nil, err
	}
	return r, nil
}
