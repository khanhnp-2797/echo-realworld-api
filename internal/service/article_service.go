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

// List returns a paginated list of articles with optional tag/author filters.
func (s *articleService) List(ctx context.Context, filter repository.ArticleFilter) ([]*domain.Article, int64, error) {
	return s.articleRepo.List(ctx, filter)
}
