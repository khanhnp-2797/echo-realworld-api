package dto

// TagsResponse is the public DTO for the GET /api/tags endpoint.
type TagsResponse struct {
	Tags []string `json:"tags"`
}
