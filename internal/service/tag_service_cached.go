package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/khanhnp-2797/echo-realworld-api/internal/cache"
	"github.com/khanhnp-2797/echo-realworld-api/internal/domain"
)

const tagCacheKey = "tags:all"
const tagCacheTTL = 5 * time.Minute

// cachedTagService wraps a TagService and caches GetAll results in Redis.
type cachedTagService struct {
	inner TagService
	cache cache.Cache
}

// NewCachedTagService decorates svc with a Redis-backed cache layer.
func NewCachedTagService(svc TagService, c cache.Cache) TagService {
	return &cachedTagService{inner: svc, cache: c}
}

func (s *cachedTagService) GetAll(ctx context.Context) ([]*domain.Tag, error) {
	// Try cache first
	if data, err := s.cache.Get(ctx, tagCacheKey); err == nil && data != nil {
		var tags []*domain.Tag
		if err := json.Unmarshal(data, &tags); err == nil {
			return tags, nil
		}
	}

	// Cache miss — fetch from DB
	tags, err := s.inner.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Populate cache (best-effort, ignore errors)
	if data, err := json.Marshal(tags); err == nil {
		_ = s.cache.Set(ctx, tagCacheKey, data, tagCacheTTL)
	}

	return tags, nil
}
