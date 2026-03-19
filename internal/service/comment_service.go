package service

import (
	"context"

	"github.com/khanhnp-2797/echo-realworld-api/internal/domain"
	"github.com/khanhnp-2797/echo-realworld-api/internal/repository"
	"github.com/khanhnp-2797/echo-realworld-api/pkg/apperrors"
)

// CommentService defines business operations on Comment.
type CommentService interface {
	// Task 3: Comments
	AddComment(ctx context.Context, slug string, authorID uint, body string) (*domain.Comment, error)
	GetComments(ctx context.Context, slug string) ([]*domain.Comment, error)

	// Task 6: Delete own comment
	DeleteComment(ctx context.Context, id uint) error
}

type commentService struct {
	commentRepo repository.CommentRepository
	articleRepo repository.ArticleRepository
}

func NewCommentService(commentRepo repository.CommentRepository, articleRepo repository.ArticleRepository) CommentService {
	return &commentService{commentRepo: commentRepo, articleRepo: articleRepo}
}

// AddComment creates a new comment on the given article slug.
func (s *commentService) AddComment(ctx context.Context, slug string, authorID uint, body string) (*domain.Comment, error) {
	article, err := s.articleRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}

	comment := &domain.Comment{
		Body:      body,
		AuthorID:  authorID,
		ArticleID: article.ID,
	}

	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return nil, err
	}

	// Re-fetch with Author preloaded
	comments, err := s.commentRepo.FindByArticleID(ctx, article.ID)
	if err != nil {
		return nil, err
	}
	for _, c := range comments {
		if c.AuthorID == authorID && c.Body == body {
			return c, nil
		}
	}
	return comment, nil
}

// GetComments returns all comments for the article identified by slug.
func (s *commentService) GetComments(ctx context.Context, slug string) ([]*domain.Comment, error) {
	article, err := s.articleRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	return s.commentRepo.FindByArticleID(ctx, article.ID)
}

// DeleteComment deletes a comment by ID.
// Ownership is enforced by the CommentOwner middleware before this is called.
func (s *commentService) DeleteComment(ctx context.Context, id uint) error {
	return s.commentRepo.Delete(ctx, id)
}
