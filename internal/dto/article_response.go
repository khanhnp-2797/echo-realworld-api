package dto

import (
	"time"

	"github.com/khanhnp-2797/echo-realworld-api/internal/domain"
)

// ArticleResponse wraps ArticleBody for single-article endpoints.
type ArticleResponse struct {
	Article ArticleBody `json:"article"`
}

// ArticlesResponse wraps a list of articles with total count.
type ArticlesResponse struct {
	Articles      []ArticleBody `json:"articles"`
	ArticlesCount int64         `json:"articlesCount"`
}

// ArticleBody is the public DTO for an Article.
type ArticleBody struct {
	Slug        string      `json:"slug"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Body        string      `json:"body"`
	TagList     []string    `json:"tagList"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
	Author      ProfileBody `json:"author"`
}

// CommentResponse wraps CommentBody for single-comment endpoints.
type CommentResponse struct {
	Comment CommentBody `json:"comment"`
}

// CommentsResponse wraps a list of comments.
type CommentsResponse struct {
	Comments []CommentBody `json:"comments"`
}

// CommentBody is the public DTO for a Comment.
type CommentBody struct {
	ID        uint        `json:"id"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
	Body      string      `json:"body"`
	Author    ProfileBody `json:"author"`
}

// ToArticleBody maps a domain.Article to an ArticleBody DTO.
func ToArticleBody(a *domain.Article) ArticleBody {
	tags := make([]string, 0, len(a.Tags))
	for _, t := range a.Tags {
		tags = append(tags, t.Name)
	}
	return ArticleBody{
		Slug:        a.Slug,
		Title:       a.Title,
		Description: a.Description,
		Body:        a.Body,
		TagList:     tags,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
		Author:      ToProfileBody(&a.Author),
	}
}

// ToCommentBody maps a domain.Comment to a CommentBody DTO.
func ToCommentBody(c *domain.Comment) CommentBody {
	return CommentBody{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Body:      c.Body,
		Author:    ToProfileBody(&c.Author),
	}
}
