package repository

import (
	"context"

	"github.com/khanhnp-2797/echo-realworld-api/internal/domain"
	"gorm.io/gorm"
)

type tagRepository struct {
	db *gorm.DB
}

// NewTagRepository returns a GORM-backed TagRepository.
func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) FindAll(ctx context.Context) ([]*domain.Tag, error) {
	var tags []*domain.Tag
	err := r.db.WithContext(ctx).Order("name ASC").Find(&tags).Error
	return tags, err
}

func (r *tagRepository) FindOrCreate(ctx context.Context, name string) (*domain.Tag, error) {
	tag := &domain.Tag{Name: name}
	result := r.db.WithContext(ctx).Where(domain.Tag{Name: name}).FirstOrCreate(tag)
	return tag, result.Error
}
