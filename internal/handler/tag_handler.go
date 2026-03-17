package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/khanhnp-2797/echo-realworld-api/internal/service"
)

// TagHandler handles /api/tags endpoint.
type TagHandler struct {
	tagSvc service.TagService
}

func NewTagHandler(tagSvc service.TagService) *TagHandler {
	return &TagHandler{tagSvc: tagSvc}
}

type tagsResponse struct {
	Tags []string `json:"tags"`
}

// API GET /api/tags — Get all tags (public)
//
// @Summary   List all tags
// @Tags      tags
// @Produce   json
// @Success   200 {object} tagsResponse
// @Failure   500 {object} map[string]any "Internal server error"
// @Router    /tags [get]
func (h *TagHandler) ListTags(c echo.Context) error {
	tags, err := h.tagSvc.GetAll(c.Request().Context())
	if err != nil {
		return handleServiceError(err)
	}

	names := make([]string, 0, len(tags))
	for _, t := range tags {
		names = append(names, t.Name)
	}

	return c.JSON(http.StatusOK, tagsResponse{Tags: names})
}
