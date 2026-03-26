package dto

// CreateArticleRequest is the request body for POST /api/articles.
type CreateArticleRequest struct {
	Article struct {
		Title       string   `json:"title"       validate:"required,min=1,max=500"`
		Description string   `json:"description" validate:"required,min=1,max=1000"`
		Body        string   `json:"body"        validate:"required"`
		TagList     []string `json:"tagList"`
	} `json:"article"`
}

// UpdateArticleRequest is the request body for PUT /api/articles/:slug.
// All fields are optional — only non-empty ones are applied.
type UpdateArticleRequest struct {
	Article struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Body        *string `json:"body"`
	} `json:"article"`
}

// AddCommentRequest is the request body for POST /api/articles/:slug/comments.
type AddCommentRequest struct {
	Comment struct {
		Body string `json:"body" validate:"required"`
	} `json:"comment"`
}
