// File: internal/cache/cache_test.go
// Purpose: Comprehensive unit tests for cache factory and in-memory cache implementation.
// Author: Backend Team
// Created: 2026-09-25
// Last Modified: 2026-09-25

package cache

import (
	"bytes"
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		redisURL    string
		wantType    string
		wantErr     bool
		errContains string
	}{
		{
			name:     "empty URL returns in-memory cache",
			redisURL: "",
			wantType: "*cache.Memory",
			wantErr:  false,
		},
		{
			name:        "invalid redis URL returns error",
			redisURL:    "invalid-url",
			wantType:    "",
			wantErr:     true,
			errContains: "redis: invalid URL scheme",
		},
		{
			name:        "unsupported scheme returns error",
			redisURL:    "http://localhost:6379",
			wantType:    "",
			wantErr:     true,
			errContains: "redis: invalid URL scheme",
		},
		{
			name:     "valid redis URL returns redis cache without connecting",
			redisURL: "redis://localhost:6379",
			wantType: "*cache.Redis",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := New(tt.redisURL)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("New(%q) expected error, got nil", tt.redisURL)
				}
				if tt.errContains != "" && !bytes.Contains([]byte(err.Error()), []byte(tt.errContains)) {
					t.Errorf("New(%q) error %q does not contain %q", tt.redisURL, err.Error(), tt.errContains)
				}
				if c != nil {
					t.Errorf("New(%q) expected nil Cacher on error, got %v", tt.redisURL, c)
				}
				return
			}

			if err != nil {
				t.Fatalf("New(%q) unexpected error: %v", tt.redisURL, err)
			}
			if c == nil {
				t.Fatalf("New(%q) returned nil Cacher", tt.redisURL)
			}

			switch tt.wantType {
			case "*cache.Memory":
				if _, ok := c.(*Memory); !ok {
					t.Errorf("New(%q) type = %T, want %s", tt.redisURL, c, tt.wantType)
				}
			case "*cache.Redis":
				if _, ok := c.(*Redis); !ok {
					t.Errorf("New(%q) type = %T, want %s", tt.redisURL, c, tt.wantType)
				}
			default:
				t.Fatalf("unhandled wantType: %s", tt.wantType)
			}
		})
	}
}

func TestMemory_Get_MissingKey(t *testing.T) {
	ctx := context.Background()
	mem := NewMemory()

	tests := []struct {
		name string
		key  string
	}{
		{
			name: "non-existent normal key",
			key:  "missing-key",
		},
		{
			name: "empty string key",
			key:  "",
		},
		{
			name: "key with special characters",
			key:  "users:123:session/token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, hit, err := mem.Get(ctx, tt.key)
			if err != nil {
				t.Errorf("Get(%q) unexpected error: %v", tt.key, err)
			}
			if hit {
				t.Errorf("Get(%q) hit = true, want false", tt.key)
			}
			if data != nil {
				t.Errorf("Get(%q) data = %v, want nil", tt.key, data)
			}
		})
	}
}

func TestMemory_SetAndGet(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		key       string
		value     []byte
		ttl       time.Duration
		wantValue []byte
	}{
		{
			name:      "standard string payload",
			key:       "user:profile:1",
			value:     []byte(`{"id":1,"name":"Alice"}`),
			ttl:       time.Minute,
			wantValue: []byte(`{"id":1,"name":"Alice"}`),
		},
		{
			name:      "empty byte slice",
			key:       "empty-val",
			value:     []byte{},
			ttl:       time.Minute,
			wantValue: []byte{},
		},
		{
			name:      "nil byte slice stores as empty slice",
			key:       "nil-val",
			value:     nil,
			ttl:       time.Minute,
			wantValue: []byte{},
		},
		{
			name:      "empty key",
			key:       "",
			value:     []byte("value-for-empty-key"),
			ttl:       time.Minute,
			wantValue: []byte("value-for-empty-key"),
		},
		{
			name:      "binary data with null bytes",
			key:       "binary-key",
			value:     []byte{0x00, 0xFF, 0xFE, 0x00, 0x12, 0x34},
			ttl:       time.Minute,
			wantValue: []byte{0x00, 0xFF, 0xFE, 0x00, 0x12, 0x34},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := NewMemory()

			if err := mem.Set(ctx, tt.key, tt.value, tt.ttl); err != nil {
				t.Fatalf("Set(%q) unexpected error: %v", tt.key, err)
			}

			data, hit, err := mem.Get(ctx, tt.key)
			if err != nil {
				t.Fatalf("Get(%q) unexpected error: %v", tt.key, err)
			}
			if !hit {
				t.Fatalf("Get(%q) hit = false, want true", tt.key)
			}
			if !bytes.Equal(data, tt.wantValue) {
				t.Errorf("Get(%q) data = %v, want %v", tt.key, data, tt.wantValue)
			}
		})
	}
}

func TestMemory_Set_OverwriteExistingKey(t *testing.T) {
	ctx := context.Background()
	mem := NewMemory()
	key := "item:10"

	if err := mem.Set(ctx, key, []byte("first"), time.Minute); err != nil {
		t.Fatalf("initial Set failed: %v", err)
	}

	if err := mem.Set(ctx, key, []byte("second"), time.Minute); err != nil {
		t.Fatalf("overwrite Set failed: %v", err)
	}

	data, hit, err := mem.Get(ctx, key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if !hit {
		t.Fatal("expected cache hit after overwrite")
	}
	if !bytes.Equal(data, []byte("second")) {
		t.Errorf("data = %s, want second", string(data))
	}
}

func TestMemory_TTL_Expiration(t *testing.T) {
	ctx := context.Background()
	mem := NewMemory()

	key := "temp-key"
	val := []byte("temp-value")
	ttl := 1 * time.Millisecond

	if err := mem.Set(ctx, key, val, ttl); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Immediate get before expiry should hit
	data, hit, err := mem.Get(ctx, key)
	if err != nil {
		t.Fatalf("Get before expiry failed: %v", err)
	}
	if !hit {
		t.Fatal("expected hit immediately after Set")
	}
	if !bytes.Equal(data, val) {
		t.Fatalf("got %s, want %s", string(data), string(val))
	}

	// Sleep 5ms to allow 1ms TTL to expire
	time.Sleep(5 * time.Millisecond)

	data, hit, err = mem.Get(ctx, key)
	if err != nil {
		t.Fatalf("Get after expiration failed: %v", err)
	}
	if hit {
		t.Errorf("Get after TTL expiry returned hit = true, want false")
	}
	if data != nil {
		t.Errorf("Get after TTL expiry returned data = %v, want nil", data)
	}
}

func TestMemory_TTL_ImmediateExpiration(t *testing.T) {
	ctx := context.Background()
	mem := NewMemory()

	// Zero or negative TTL should expire immediately
	tests := []struct {
		name string
		ttl  time.Duration
	}{
		{name: "zero TTL", ttl: 0},
		{name: "negative TTL", ttl: -1 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := "expired-immediately"
			if err := mem.Set(ctx, key, []byte("val"), tt.ttl); err != nil {
				t.Fatalf("Set failed: %v", err)
			}

			data, hit, err := mem.Get(ctx, key)
			if err != nil {
				t.Fatalf("Get failed: %v", err)
			}
			if hit {
				t.Errorf("Get with TTL %v returned hit = true, want false", tt.ttl)
			}
			if data != nil {
				t.Errorf("Get with TTL %v returned data = %v, want nil", tt.ttl, data)
			}
		})
	}
}

func TestMemory_Del(t *testing.T) {
	ctx := context.Background()

	t.Run("delete existing key", func(t *testing.T) {
		mem := NewMemory()
		key := "to-delete"
		if err := mem.Set(ctx, key, []byte("data"), time.Minute); err != nil {
			t.Fatalf("Set failed: %v", err)
		}

		if err := mem.Del(ctx, key); err != nil {
			t.Fatalf("Del failed: %v", err)
		}

		data, hit, err := mem.Get(ctx, key)
		if err != nil {
			t.Fatalf("Get after Del failed: %v", err)
		}
		if hit {
			t.Errorf("Get after Del returned hit = true, want false")
		}
		if data != nil {
			t.Errorf("Get after Del returned data = %v, want nil", data)
		}
	})

	t.Run("delete non-existent key returns nil error", func(t *testing.T) {
		mem := NewMemory()
		if err := mem.Del(ctx, "does-not-exist"); err != nil {
			t.Errorf("Del on missing key returned error: %v", err)
		}
	})

	t.Run("delete empty string key returns nil error", func(t *testing.T) {
		mem := NewMemory()
		if err := mem.Del(ctx, ""); err != nil {
			t.Errorf("Del on empty key returned error: %v", err)
		}
	})
}

func TestMemory_GetOrSet(t *testing.T) {
	ctx := context.Background()

	t.Run("miss calls fetch and caches value, hit returns cached on second call", func(t *testing.T) {
		mem := NewMemory()
		key := "item:profile"
		expectedData := []byte("fetched-data")
		fetchCount := 0

		fetchFn := func() ([]byte, error) {
			fetchCount++
			return expectedData, nil
		}

		// First call: cache miss
		data, hit, err := mem.GetOrSet(ctx, key, time.Minute, fetchFn)
		if err != nil {
			t.Fatalf("GetOrSet miss unexpected error: %v", err)
		}
		if hit {
			t.Errorf("first GetOrSet returned hit = true, want false")
		}
		if !bytes.Equal(data, expectedData) {
			t.Errorf("first GetOrSet returned %s, want %s", string(data), string(expectedData))
		}
		if fetchCount != 1 {
			t.Fatalf("fetchCount = %d, want 1", fetchCount)
		}

		// Second call: cache hit
		data2, hit2, err2 := mem.GetOrSet(ctx, key, time.Minute, fetchFn)
		if err2 != nil {
			t.Fatalf("GetOrSet hit unexpected error: %v", err2)
		}
		if !hit2 {
			t.Errorf("second GetOrSet returned hit = false, want true")
		}
		if !bytes.Equal(data2, expectedData) {
			t.Errorf("second GetOrSet returned %s, want %s", string(data2), string(expectedData))
		}
		if fetchCount != 1 {
			t.Fatalf("fetchCount = %d, want 1 (fetch should not be called again)", fetchCount)
		}

		// Direct Get verifies value is stored
		directData, directHit, directErr := mem.Get(ctx, key)
		if directErr != nil {
			t.Fatalf("direct Get unexpected error: %v", directErr)
		}
		if !directHit {
			t.Fatal("direct Get hit = false, want true")
		}
		if !bytes.Equal(directData, expectedData) {
			t.Errorf("direct Get returned %s, want %s", string(directData), string(expectedData))
		}
	})

	t.Run("fetch error propagates and does not cache", func(t *testing.T) {
		mem := NewMemory()
		key := "failing-fetch"
		expectedErr := errors.New("upstream fetch failure")

		data, hit, err := mem.GetOrSet(ctx, key, time.Minute, func() ([]byte, error) {
			return nil, expectedErr
		})

		if !errors.Is(err, expectedErr) {
			t.Fatalf("GetOrSet error = %v, want %v", err, expectedErr)
		}
		if hit {
			t.Errorf("GetOrSet on fetch error returned hit = true, want false")
		}
		if data != nil {
			t.Errorf("GetOrSet on fetch error returned data = %v, want nil", data)
		}

		// Ensure key was not cached
		_, hitAfter, errAfter := mem.Get(ctx, key)
		if errAfter != nil {
			t.Fatalf("Get after failed fetch error: %v", errAfter)
		}
		if hitAfter {
			t.Errorf("key was unexpectedly cached after fetch failure")
		}
	})

	t.Run("expired entry calls fetch again", func(t *testing.T) {
		mem := NewMemory()
		key := "expiring-fetch"
		calls := 0

		fetchFn := func() ([]byte, error) {
			calls++
			if calls == 1 {
				return []byte("first-run"), nil
			}
			return []byte("second-run"), nil
		}

		data1, hit1, err1 := mem.GetOrSet(ctx, key, 1*time.Millisecond, fetchFn)
		if err1 != nil || hit1 || !bytes.Equal(data1, []byte("first-run")) {
			t.Fatalf("first GetOrSet unexpected result: data=%s, hit=%v, err=%v", data1, hit1, err1)
		}

		time.Sleep(5 * time.Millisecond)

		data2, hit2, err2 := mem.GetOrSet(ctx, key, time.Minute, fetchFn)
		if err2 != nil {
			t.Fatalf("second GetOrSet unexpected error: %v", err2)
		}
		if hit2 {
			t.Errorf("second GetOrSet returned hit = true on expired entry, want false")
		}
		if !bytes.Equal(data2, []byte("second-run")) {
			t.Errorf("second GetOrSet returned %s, want second-run", string(data2))
		}
		if calls != 2 {
			t.Errorf("calls = %d, want 2", calls)
		}
	})
}

func TestMemory_Ping(t *testing.T) {
	mem := NewMemory()

	tests := []struct {
		name string
		ctx  func() (context.Context, context.CancelFunc)
	}{
		{
			name: "background context",
			ctx: func() (context.Context, context.CancelFunc) {
				return context.Background(), func() {}
			},
		},
		{
			name: "canceled context",
			ctx: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx, cancel
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := tt.ctx()
			defer cancel()

			if err := mem.Ping(ctx); err != nil {
				t.Errorf("Ping() returned error: %v, want nil", err)
			}
		})
	}
}

func TestMemory_Close(t *testing.T) {
	mem := NewMemory()
	if err := mem.Close(); err != nil {
		t.Errorf("Close() returned error: %v, want nil", err)
	}
}

func TestMemory_DataImmutability(t *testing.T) {
	ctx := context.Background()

	t.Run("mutating input slice after Set does not corrupt cached data", func(t *testing.T) {
		mem := NewMemory()
		key := "input-immutability"
		input := []byte("original-payload")

		if err := mem.Set(ctx, key, input, time.Minute); err != nil {
			t.Fatalf("Set failed: %v", err)
		}

		// Mutate input slice
		input[0] = 'X'

		data, hit, err := mem.Get(ctx, key)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if !hit {
			t.Fatal("expected hit = true")
		}
		if !bytes.Equal(data, []byte("original-payload")) {
			t.Errorf("cached data was corrupted by mutating input: got %s, want %s", string(data), "original-payload")
		}
	})

	t.Run("mutating returned slice from Get does not corrupt cached data", func(t *testing.T) {
		mem := NewMemory()
		key := "output-immutability"

		if err := mem.Set(ctx, key, []byte("immutable-value"), time.Minute); err != nil {
			t.Fatalf("Set failed: %v", err)
		}

		data1, hit1, err1 := mem.Get(ctx, key)
		if err1 != nil || !hit1 {
			t.Fatalf("first Get failed: hit=%v, err=%v", hit1, err1)
		}

		// Mutate the returned slice
		data1[0] = 'Z'

		data2, hit2, err2 := mem.Get(ctx, key)
		if err2 != nil || !hit2 {
			t.Fatalf("second Get failed: hit=%v, err=%v", hit2, err2)
		}

		if !bytes.Equal(data2, []byte("immutable-value")) {
			t.Errorf("cached data was corrupted by mutating Get output: got %s, want immutable-value", string(data2))
		}
	})

	t.Run("mutating returned slice from GetOrSet does not corrupt cached data", func(t *testing.T) {
		mem := NewMemory()
		key := "getorset-immutability"

		fetchResult := []byte("fetched-immutable")
		data1, _, err := mem.GetOrSet(ctx, key, time.Minute, func() ([]byte, error) {
			return fetchResult, nil
		})
		if err != nil {
			t.Fatalf("GetOrSet failed: %v", err)
		}

		// Mutate fetchResult and data1
		fetchResult[0] = 'Q'
		data1[0] = 'W'

		data2, hit, err := mem.Get(ctx, key)
		if err != nil || !hit {
			t.Fatalf("Get failed: hit=%v, err=%v", hit, err)
		}

		if !bytes.Equal(data2, []byte("fetched-immutable")) {
			t.Errorf("cached data was corrupted by mutating GetOrSet output: got %s, want fetched-immutable", string(data2))
		}
	})
}

func TestMemory_Concurrency(t *testing.T) {
	ctx := context.Background()
	mem := NewMemory()
	key := "concurrent-key"

	var wg sync.WaitGroup
	workers := 20
	iterations := 100

	var fetchCalls atomic.Int64

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				// Perform mixed operations concurrently
				switch (workerID + j) % 4 {
				case 0:
					_ = mem.Set(ctx, key, []byte("data"), time.Minute)
				case 1:
					_, _, _ = mem.Get(ctx, key)
				case 2:
					_ = mem.Del(ctx, key)
				case 3:
					_, _, _ = mem.GetOrSet(ctx, key, time.Minute, func() ([]byte, error) {
						fetchCalls.Add(1)
						return []byte("concurrent-fetch"), nil
					})
				}
			}
		}(i)
	}

	wg.Wait()
}
