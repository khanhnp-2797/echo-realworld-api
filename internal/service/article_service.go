package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/khanhnp-2797/echo-realworld-api/internal/domain"
	"github.com/khanhnp-2797/echo-realworld-api/internal/repository"
	"github.com/khanhnp-2797/echo-realworld-api/pkg/apperrors"
)

// ArticleService defines business operations on Article.
type ArticleService interface {
	// CRUD
	Create(ctx context.Context, authorID uint, title, description, body string, tagList []string) (*domain.Article, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Article, error)
	List(ctx context.Context, filter repository.ArticleFilter) ([]*domain.Article, int64, error)
	Update(ctx context.Context, slug string, authorID uint, title, description, body *string) (*domain.Article, error)
	DeleteBySlug(ctx context.Context, slug string, authorID uint) error

	// Task 5: Feed
	Feed(ctx context.Context, userID uint, filter repository.ArticleFilter) ([]*domain.Article, int64, error)

	// Task 6: Favorites
	Favorite(ctx context.Context, userID uint, slug string) (*domain.Article, error)
	Unfavorite(ctx context.Context, userID uint, slug string) (*domain.Article, error)
}

type articleService struct {
	articleRepo repository.ArticleRepository
	tagRepo     repository.TagRepository
}

func NewArticleService(articleRepo repository.ArticleRepository, tagRepo repository.TagRepository) ArticleService {
	return &articleService{articleRepo: articleRepo, tagRepo: tagRepo}
}

// slugify converts a title to a URL-friendly lowercase slug.
func slugify(title string) string {
	s := strings.ToLower(title)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, s)
	return s
}

// Create builds and persists a new article, resolving tag names into Tag records.
func (s *articleService) Create(ctx context.Context, authorID uint, title, description, body string, tagList []string) (*domain.Article, error) {
	slug := fmt.Sprintf("%s-%d", slugify(title), authorID)

	var tags []*domain.Tag
	if len(tagList) > 0 {
		var err error
		tags, err = s.tagRepo.FindOrCreateByNames(ctx, tagList)
		if err != nil {
			return nil, err
		}
	}

	article := &domain.Article{
		Slug:        slug,
		Title:       title,
		Description: description,
		Body:        body,
		AuthorID:    authorID,
		Tags:        tags,
	}

	if err := s.articleRepo.Create(ctx, article); err != nil {
		if strings.Contains(err.Error(), "idx_articles_slug") || strings.Contains(err.Error(), "unique") {
			return nil, apperrors.ErrConflict
		}
		return nil, err
	}

	return s.articleRepo.FindBySlug(ctx, article.Slug)
}

// Update applies non-nil field patches to an existing article owned by authorID.
func (s *articleService) Update(ctx context.Context, slug string, authorID uint, title, description, body *string) (*domain.Article, error) {
	article, err := s.articleRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	if article.AuthorID != authorID {
		return nil, apperrors.ErrForbidden
	}

	if title != nil {
		article.Title = *title
		article.Slug = fmt.Sprintf("%s-%d", slugify(*title), authorID)
	}
	if description != nil {
		article.Description = *description
	}
	if body != nil {
		article.Body = *body
	}

	if err := s.articleRepo.Update(ctx, article); err != nil {
		return nil, err
	}
	return s.articleRepo.FindBySlug(ctx, article.Slug)
}

// DeleteBySlug soft-deletes an article owned by authorID.
func (s *articleService) DeleteBySlug(ctx context.Context, slug string, authorID uint) error {
	return s.articleRepo.Delete(ctx, slug, authorID)
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
