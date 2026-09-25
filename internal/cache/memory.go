// File: internal/cache/memory.go
// Purpose: In-memory cache implementation for local development without Redis.
// Author: Backend Team
// Created: 2026-09-25
// Last Modified: 2026-09-25

package cache

import (
	"context"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

// entry is a single cached value with an expiration time.
type entry struct {
	data      []byte
	expiresAt time.Time
}

// isExpired returns true when the entry has passed its TTL.
func (e entry) isExpired() bool {
	return time.Now().After(e.expiresAt)
}

// Memory is a simple thread-safe in-memory cache.
// Suitable for local development; not shared across instances.
type Memory struct {
	mu    sync.RWMutex
	store map[string]entry
	sg    singleflight.Group
}

// NewMemory creates a ready-to-use in-memory cache.
func NewMemory() *Memory {
	return &Memory{
		store: make(map[string]entry),
	}
}

// Get returns cached bytes for a key. Returns hit=false on a miss or expired entry.
func (m *Memory) Get(_ context.Context, key string) ([]byte, bool, error) {
	m.mu.RLock()
	e, ok := m.store[key]
	m.mu.RUnlock()

	if !ok || e.isExpired() {
		return nil, false, nil
	}

	// Return a copy to prevent callers from mutating cached data.
	cp := make([]byte, len(e.data))
	copy(cp, e.data)
	return cp, true, nil
}

// Set stores bytes under a key with a time-to-live.
func (m *Memory) Set(_ context.Context, key string, value []byte, ttl time.Duration) error {
	cp := make([]byte, len(value))
	copy(cp, value)

	m.mu.Lock()
	m.store[key] = entry{data: cp, expiresAt: time.Now().Add(ttl)}
	m.mu.Unlock()
	return nil
}

// Del removes a key from the cache.
func (m *Memory) Del(_ context.Context, key string) error {
	m.mu.Lock()
	delete(m.store, key)
	m.mu.Unlock()
	return nil
}

// GetOrSet returns cached data if present (hit=true). On a miss it calls
// fetch(), stores the result, and returns it (hit=false). Uses singleflight
// to prevent cache stampedes on concurrent misses.
func (m *Memory) GetOrSet(ctx context.Context, key string, ttl time.Duration, fetch func() ([]byte, error)) ([]byte, bool, error) {
	if cached, found, err := m.Get(ctx, key); err != nil {
		return nil, false, err
	} else if found {
		return cached, true, nil
	}

	v, err, _ := m.sg.Do(key, func() (interface{}, error) {
		// Double check if it was populated while waiting
		if cached, found, _ := m.Get(ctx, key); found {
			return cached, nil
		}
		fresh, err := fetch()
		if err != nil {
			return nil, err
		}
		_ = m.Set(ctx, key, fresh, ttl)
		return fresh, nil
	})

	if err != nil {
		return nil, false, err
	}
	return v.([]byte), false, nil
}

// Ping always succeeds for in-memory cache.
func (m *Memory) Ping(_ context.Context) error {
	return nil
}

// Close is a no-op for in-memory cache.
func (m *Memory) Close() error {
	return nil
}
