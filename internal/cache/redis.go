// File: internal/cache/redis.go
// Purpose: Redis-backed cache implementation for production use.
// Author: Backend Team
// Created: 2026-09-25
// Last Modified: 2026-09-25

package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

// Redis wraps a Redis client implementing the Cacher interface.
type Redis struct {
	rdb *redis.Client
	sg  singleflight.Group
}

// NewRedis parses a redis:// URL and returns a ready Redis cache.
func NewRedis(redisURL string) (*Redis, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	return &Redis{rdb: redis.NewClient(opts)}, nil
}

// Get returns the raw cached bytes for a key. hit=false on a cache miss.
func (r *Redis) Get(ctx context.Context, key string) ([]byte, bool, error) {
	val, err := r.rdb.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, false, nil // normal miss
	}
	if err != nil {
		return nil, false, err // real Redis error
	}
	return val, true, nil
}

// Set stores bytes under a key with a time-to-live.
func (r *Redis) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return r.rdb.Set(ctx, key, value, ttl).Err()
}

// Del removes a key from the cache.
func (r *Redis) Del(ctx context.Context, key string) error {
	return r.rdb.Del(ctx, key).Err()
}

// GetOrSet returns cached bytes for key if present (hit=true). On a miss it
// calls fetch(), stores the result under key with ttl, and returns it (hit=false).
// Uses singleflight to prevent cache stampedes on concurrent misses.
func (r *Redis) GetOrSet(
	ctx context.Context,
	key string,
	ttl time.Duration,
	fetch func() ([]byte, error),
) (data []byte, hit bool, err error) {
	if cached, found, gErr := r.Get(ctx, key); gErr != nil {
		return nil, false, gErr
	} else if found {
		return cached, true, nil
	}

	v, err, _ := r.sg.Do(key, func() (interface{}, error) {
		if cached, found, _ := r.Get(ctx, key); found {
			return cached, nil
		}
		fresh, err := fetch()
		if err != nil {
			return nil, err
		}
		_ = r.Set(ctx, key, fresh, ttl)
		return fresh, nil
	})

	if err != nil {
		return nil, false, err
	}
	return v.([]byte), false, nil
}

// Ping checks the connection to Redis.
func (r *Redis) Ping(ctx context.Context) error {
	return r.rdb.Ping(ctx).Err()
}

// Close shuts down the Redis client connection pool.
func (r *Redis) Close() error {
	return r.rdb.Close()
}
