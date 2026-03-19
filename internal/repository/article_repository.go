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

// FindBySlug fetches a single article with Author, Tags, Comments and FavoritedBy preloaded.
func (r *articleRepository) FindBySlug(ctx context.Context, slug string) (*domain.Article, error) {
	var article domain.Article
	err := r.db.WithContext(ctx).
		Preload("Author").
		Preload("Tags").
		Preload("Comments.Author").
		Preload("FavoritedBy").
		Where("slug = ?", slug).
		First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// List returns a paginated, optionally-filtered list of articles.
// Uses subqueries for tag/favorited filters to avoid duplicate rows.
func (r *articleRepository) List(ctx context.Context, filter ArticleFilter) ([]*domain.Article, int64, error) {
	var articles []*domain.Article
	var count int64

	q := r.db.WithContext(ctx).Model(&domain.Article{})

	if filter.Tag != "" {
		q = q.Where("id IN (?)",
			r.db.Table("article_tags").
				Select("article_id").
				Joins("JOIN tags ON tags.id = article_tags.tag_id").
				Where("tags.name = ?", filter.Tag))
	}
	if filter.Author != "" {
		q = q.Joins("JOIN users ON users.id = articles.author_id").
			Where("users.username = ?", filter.Author)
	}
	if filter.Favorited != "" {
		q = q.Where("id IN (?)",
			r.db.Table("article_favorites").
				Select("article_id").
				Joins("JOIN users ON users.id = article_favorites.user_id").
				Where("users.username = ?", filter.Favorited))
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

	err := q.
		Preload("Author").
		Preload("Tags").
		Preload("FavoritedBy").
		Order("articles.created_at DESC").
		Limit(limit).Offset(offset).
		Find(&articles).Error
	return articles, count, err
}

// Feed returns paginated articles from users that viewerID follows.
func (r *articleRepository) Feed(ctx context.Context, userID uint, filter ArticleFilter) ([]*domain.Article, int64, error) {
	var articles []*domain.Article
	var count int64

	q := r.db.WithContext(ctx).
		Model(&domain.Article{}).
		Where("author_id IN (?)",
			r.db.Table("user_follows").
				Select("followed_id").
				Where("follower_id = ?", userID))

	q.Count(&count)

	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}

	err := q.
		Preload("Author").
		Preload("Tags").
		Preload("FavoritedBy").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&articles).Error
	return articles, count, err
}

// Favorite adds userID to the article's FavoritedBy list.
func (r *articleRepository) Favorite(ctx context.Context, userID, articleID uint) error {
	user := domain.User{}
	user.ID = userID
	article := domain.Article{}
	article.ID = articleID
	return r.db.WithContext(ctx).Model(&article).Association("FavoritedBy").Append(&user)
}

// Unfavorite removes userID from the article's FavoritedBy list.
func (r *articleRepository) Unfavorite(ctx context.Context, userID, articleID uint) error {
	user := domain.User{}
	user.ID = userID
	article := domain.Article{}
	article.ID = articleID
	return r.db.WithContext(ctx).Model(&article).Association("FavoritedBy").Delete(&user)
}
