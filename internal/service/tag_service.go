package service

import (
	"context"

	"github.com/khanhnp-2797/echo-realworld-api/internal/domain"
	"github.com/khanhnp-2797/echo-realworld-api/internal/repository"
)

// TagService defines business operations on Tag.
type TagService interface {
	GetAll(ctx context.Context) ([]*domain.Tag, error)
}

type tagService struct {
	repo repository.TagRepository
}

func NewTagService(repo repository.TagRepository) TagService {
	return &tagService{repo: repo}
}

func (s *tagService) GetAll(ctx context.Context) ([]*domain.Tag, error) {
	return s.repo.FindAll(ctx)
}
