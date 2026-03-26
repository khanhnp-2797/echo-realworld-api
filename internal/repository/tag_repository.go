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

// FindOrCreateByNames resolves a slice of tag name strings to domain.Tag records,
// creating any that do not yet exist.
func (r *tagRepository) FindOrCreateByNames(ctx context.Context, names []string) ([]*domain.Tag, error) {
	tags := make([]*domain.Tag, 0, len(names))
	for _, name := range names {
		tag := &domain.Tag{Name: name}
		if err := r.db.WithContext(ctx).Where(domain.Tag{Name: name}).FirstOrCreate(tag).Error; err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}
