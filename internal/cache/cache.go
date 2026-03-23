package cache

import (
	"context"
	"time"
)

// Cache is a generic key-value cache interface.
// Implementations can be Redis, in-memory, or a no-op for tests.
type Cache interface {
	// Get retrieves the cached bytes for key. Returns (nil, nil) on cache miss.
	Get(ctx context.Context, key string) ([]byte, error)

	// Set stores value for key with the given TTL.
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error

	// Del removes one or more keys.
	Del(ctx context.Context, keys ...string) error
}
