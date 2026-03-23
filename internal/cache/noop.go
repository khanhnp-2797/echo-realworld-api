package cache

import (
	"context"
	"time"
)

// NoopCache is a cache that always misses. Useful for tests and when Redis
// is disabled.
type NoopCache struct{}

func (n *NoopCache) Get(_ context.Context, _ string) ([]byte, error)                  { return nil, nil }
func (n *NoopCache) Set(_ context.Context, _ string, _ []byte, _ time.Duration) error { return nil }
func (n *NoopCache) Del(_ context.Context, _ ...string) error                         { return nil }
