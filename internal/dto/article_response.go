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
	Slug           string      `json:"slug"`
	Title          string      `json:"title"`
	Description    string      `json:"description"`
	Body           string      `json:"body"`
	TagList        []string    `json:"tagList"`
	CreatedAt      time.Time   `json:"createdAt"`
	UpdatedAt      time.Time   `json:"updatedAt"`
	Favorited      bool        `json:"favorited"`
	FavoritesCount int         `json:"favoritesCount"`
	Author         ProfileBody `json:"author"`
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
// viewerID is used to compute favorited; following is pre-computed by the caller.
func ToArticleBody(a *domain.Article, viewerID uint, following bool) ArticleBody {
	tags := make([]string, 0, len(a.Tags))
	for _, t := range a.Tags {
		tags = append(tags, t.Name)
	}
	favorited := false
	for _, u := range a.FavoritedBy {
		if u.ID == viewerID {
			favorited = true
			break
		}
	}
	return ArticleBody{
		Slug:           a.Slug,
		Title:          a.Title,
		Description:    a.Description,
		Body:           a.Body,
		TagList:        tags,
		CreatedAt:      a.CreatedAt,
		UpdatedAt:      a.UpdatedAt,
		Favorited:      favorited,
		FavoritesCount: len(a.FavoritedBy),
		Author:         ToProfileBody(&a.Author, following),
	}
}

// ToCommentBody maps a domain.Comment to a CommentBody DTO.
// following indicates whether the viewer follows the comment author.
func ToCommentBody(c *domain.Comment, following bool) CommentBody {
	return CommentBody{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Body:      c.Body,
		Author:    ToProfileBody(&c.Author, following),
	}
}
