package service

import (
	"context"

	"github.com/khanhnp-2797/echo-realworld-api/internal/domain"
	"github.com/khanhnp-2797/echo-realworld-api/internal/repository"
	"github.com/khanhnp-2797/echo-realworld-api/pkg/apperrors"
)

// ArticleService defines business operations on Article.
type ArticleService interface {
	// Task 2: Read-only CRUD
	GetBySlug(ctx context.Context, slug string) (*domain.Article, error)
	List(ctx context.Context, filter repository.ArticleFilter) ([]*domain.Article, int64, error)

	// Task 5: Feed
	Feed(ctx context.Context, userID uint, filter repository.ArticleFilter) ([]*domain.Article, int64, error)

	// Task 6: Favorites
	Favorite(ctx context.Context, userID uint, slug string) (*domain.Article, error)
	Unfavorite(ctx context.Context, userID uint, slug string) (*domain.Article, error)
}

type articleService struct {
	articleRepo repository.ArticleRepository
}

func NewArticleService(articleRepo repository.ArticleRepository) ArticleService {
	return &articleService{articleRepo: articleRepo}
}

// GetBySlug fetches a single article by its slug (with Author, Tags, Comments preloaded).
func (s *articleService) GetBySlug(ctx context.Context, slug string) (*domain.Article, error) {
	article, err := s.articleRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	return article, nil
}

// List returns a paginated list of articles with optional tag/author/favorited filters.
func (s *articleService) List(ctx context.Context, filter repository.ArticleFilter) ([]*domain.Article, int64, error) {
	return s.articleRepo.List(ctx, filter)
}

// Feed returns articles from authors that userID follows.
func (s *articleService) Feed(ctx context.Context, userID uint, filter repository.ArticleFilter) ([]*domain.Article, int64, error) {
	return s.articleRepo.Feed(ctx, userID, filter)
}

// Favorite adds the article to the user's favorites.
func (s *articleService) Favorite(ctx context.Context, userID uint, slug string) (*domain.Article, error) {
	article, err := s.articleRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	if err := s.articleRepo.Favorite(ctx, userID, article.ID); err != nil {
		return nil, err
	}
	// Re-fetch to get updated FavoritedBy list
	return s.articleRepo.FindBySlug(ctx, slug)
}

// Unfavorite removes the article from the user's favorites.
func (s *articleService) Unfavorite(ctx context.Context, userID uint, slug string) (*domain.Article, error) {
	article, err := s.articleRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	if err := s.articleRepo.Unfavorite(ctx, userID, article.ID); err != nil {
		return nil, err
	}
	// Re-fetch to get updated FavoritedBy list
	return s.articleRepo.FindBySlug(ctx, slug)
}
