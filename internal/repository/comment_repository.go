package repository

import (
	"context"

	"github.com/khanhnp-2797/echo-realworld-api/internal/domain"
	"gorm.io/gorm"
)

type commentRepository struct {
	db *gorm.DB
}

// NewCommentRepository returns a GORM-backed CommentRepository.
func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(ctx context.Context, comment *domain.Comment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

// FindByArticleID fetches all comments for an article, with Author preloaded.
func (r *commentRepository) FindByArticleID(ctx context.Context, articleID uint) ([]*domain.Comment, error) {
	var comments []*domain.Comment
	err := r.db.WithContext(ctx).
		Preload("Author").
		Where("article_id = ?", articleID).
		Order("created_at DESC").
		Find(&comments).Error
	return comments, err
}
