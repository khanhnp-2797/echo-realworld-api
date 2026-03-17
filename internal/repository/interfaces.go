package repository

import (
	"context"

	"github.com/khanhnp-2797/echo-realworld-api/internal/domain"
)

// ArticleFilter holds optional query parameters when listing articles.
type ArticleFilter struct {
	Tag    string
	Author string
	Limit  int
	Offset int
}

// UserRepository — persistence operations for User.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id uint) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
}

// ArticleRepository — persistence operations for Article.
type ArticleRepository interface {
	FindBySlug(ctx context.Context, slug string) (*domain.Article, error)
	List(ctx context.Context, filter ArticleFilter) ([]*domain.Article, int64, error)
}

// CommentRepository — persistence operations for Comment.
type CommentRepository interface {
	Create(ctx context.Context, comment *domain.Comment) error
	FindByArticleID(ctx context.Context, articleID uint) ([]*domain.Comment, error)
}

// TagRepository — persistence operations for Tag.
type TagRepository interface {
	FindAll(ctx context.Context) ([]*domain.Tag, error)
}
