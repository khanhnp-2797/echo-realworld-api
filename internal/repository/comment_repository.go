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

// FindByID fetches a single comment by ID with Author preloaded.
func (r *commentRepository) FindByID(ctx context.Context, id uint) (*domain.Comment, error) {
	var comment domain.Comment
	if err := r.db.WithContext(ctx).Preload("Author").First(&comment, id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
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

// Delete soft-deletes a comment by ID.
func (r *commentRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Comment{}, id).Error
}
