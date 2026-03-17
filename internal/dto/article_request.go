package dto

// AddCommentRequest is the request body for POST /api/articles/:slug/comments.
type AddCommentRequest struct {
	Comment struct {
		Body string `json:"body" validate:"required"`
	} `json:"comment"`
}
