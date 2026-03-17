package repository

import (
	"context"

	"github.com/khanhnp-2797/echo-realworld-api/internal/domain"
	"gorm.io/gorm"
)

type articleRepository struct {
	db *gorm.DB
}

// NewArticleRepository returns a GORM-backed ArticleRepository.
func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db: db}
}

// FindBySlug fetches a single article with Author, Tags and Comments preloaded.
func (r *articleRepository) FindBySlug(ctx context.Context, slug string) (*domain.Article, error) {
	var article domain.Article
	err := r.db.WithContext(ctx).
		Preload("Author").
		Preload("Tags").
		Preload("Comments.Author").
		Where("slug = ?", slug).
		First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// List returns a paginated, optionally-filtered list of articles.
func (r *articleRepository) List(ctx context.Context, filter ArticleFilter) ([]*domain.Article, int64, error) {
	var articles []*domain.Article
	var count int64

	q := r.db.WithContext(ctx).
		Model(&domain.Article{}).
		Preload("Author").
		Preload("Tags")

	if filter.Tag != "" {
		q = q.Joins("JOIN article_tags ON article_tags.article_id = articles.id").
			Joins("JOIN tags ON tags.id = article_tags.tag_id").
			Where("tags.name = ?", filter.Tag)
	}
	if filter.Author != "" {
		q = q.Joins("JOIN users AS authors ON authors.id = articles.author_id").
			Where("authors.username = ?", filter.Author)
	}

	q.Count(&count)

	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}

	err := q.Order("articles.created_at DESC").Limit(limit).Offset(offset).Find(&articles).Error
	return articles, count, err
}
