package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/khanhnp-2797/echo-realworld-api/internal/cache"
	"github.com/khanhnp-2797/echo-realworld-api/internal/domain"
	"github.com/khanhnp-2797/echo-realworld-api/internal/repository"
)

const articleListCacheTTL = 2 * time.Minute

// articleListCacheKey builds a deterministic key from the filter parameters.
func articleListCacheKey(filter repository.ArticleFilter) string {
	return fmt.Sprintf("articles:list:tag=%s:author=%s:fav=%s:limit=%d:offset=%d",
		filter.Tag, filter.Author, filter.Favorited, filter.Limit, filter.Offset)
}

// articleListResult is used for JSON marshalling the List response.
type articleListResult struct {
	Articles []*domain.Article `json:"articles"`
	Count    int64             `json:"count"`
}

// cachedArticleService wraps an ArticleService and caches List results.
type cachedArticleService struct {
	inner ArticleService
	cache cache.Cache
}

// NewCachedArticleService decorates svc with a cache layer for List only.
// All other operations (GetBySlug, Feed, Favorite, Unfavorite) bypass cache.
func NewCachedArticleService(svc ArticleService, c cache.Cache) ArticleService {
	return &cachedArticleService{inner: svc, cache: c}
}

func (s *cachedArticleService) List(ctx context.Context, filter repository.ArticleFilter) ([]*domain.Article, int64, error) {
	key := articleListCacheKey(filter)

	if data, err := s.cache.Get(ctx, key); err == nil && data != nil {
		var result articleListResult
		if err := json.Unmarshal(data, &result); err == nil {
			return result.Articles, result.Count, nil
		}
	}

	articles, count, err := s.inner.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	if data, err := json.Marshal(articleListResult{Articles: articles, Count: count}); err == nil {
		_ = s.cache.Set(ctx, key, data, articleListCacheTTL)
	}

	return articles, count, nil
}

// All remaining methods delegate directly to the inner service.

func (s *cachedArticleService) Create(ctx context.Context, authorID uint, title, description, body string, tagList []string) (*domain.Article, error) {
	return s.inner.Create(ctx, authorID, title, description, body, tagList)
}

func (s *cachedArticleService) GetBySlug(ctx context.Context, slug string) (*domain.Article, error) {
	return s.inner.GetBySlug(ctx, slug)
}

func (s *cachedArticleService) Update(ctx context.Context, slug string, authorID uint, title, description, body *string) (*domain.Article, error) {
	return s.inner.Update(ctx, slug, authorID, title, description, body)
}

func (s *cachedArticleService) DeleteBySlug(ctx context.Context, slug string, authorID uint) error {
	return s.inner.DeleteBySlug(ctx, slug, authorID)
}

func (s *cachedArticleService) Feed(ctx context.Context, userID uint, filter repository.ArticleFilter) ([]*domain.Article, int64, error) {
	return s.inner.Feed(ctx, userID, filter)
}

func (s *cachedArticleService) Favorite(ctx context.Context, userID uint, slug string) (*domain.Article, error) {
	return s.inner.Favorite(ctx, userID, slug)
}

func (s *cachedArticleService) Unfavorite(ctx context.Context, userID uint, slug string) (*domain.Article, error) {
	return s.inner.Unfavorite(ctx, userID, slug)
}
