package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/khanhnp-2797/echo-realworld-api/internal/dto"
	"github.com/khanhnp-2797/echo-realworld-api/internal/middleware"
	"github.com/khanhnp-2797/echo-realworld-api/internal/repository"
	"github.com/khanhnp-2797/echo-realworld-api/internal/service"
)

// ArticleHandler handles article and comment endpoints.
type ArticleHandler struct {
	articleSvc service.ArticleService
	commentSvc service.CommentService
}

func NewArticleHandler(
	articleSvc service.ArticleService,
	commentSvc service.CommentService,
) *ArticleHandler {
	return &ArticleHandler{articleSvc: articleSvc, commentSvc: commentSvc}
}

// ──────────────────────────── Handlers ────────────────────────────

// API GET /api/articles — List articles with optional filters (public)
//
// @Summary   List articles
// @Tags      articles
// @Produce   json
// @Param     tag    query string false "Filter by tag"
// @Param     author query string false "Filter by author username"
// @Param     limit  query int    false "Limit (default 20)"
// @Param     offset query int    false "Offset (default 0)"
// @Success   200 {object} dto.ArticlesResponse
// @Failure   500 {object} map[string]any "Internal server error"
// @Router    /articles [get]
func (h *ArticleHandler) ListArticles(c echo.Context) error {
	filter := repository.ArticleFilter{
		Tag:    c.QueryParam("tag"),
		Author: c.QueryParam("author"),
		Limit:  queryInt(c, "limit", 20),
		Offset: queryInt(c, "offset", 0),
	}

	articles, count, err := h.articleSvc.List(c.Request().Context(), filter)
	if err != nil {
		return handleServiceError(err)
	}

	bodies := make([]dto.ArticleBody, 0, len(articles))
	for _, a := range articles {
		bodies = append(bodies, dto.ToArticleBody(a))
	}

	return c.JSON(http.StatusOK, dto.ArticlesResponse{Articles: bodies, ArticlesCount: count})
}

// API GET /api/articles/:slug — Get a single article by slug (public)
//
// @Summary   Get article
// @Tags      articles
// @Produce   json
// @Param     slug path string true "Article slug"
// @Success   200 {object} dto.ArticleResponse
// @Failure   404 {object} map[string]any "Article not found"
// @Router    /articles/{slug} [get]
func (h *ArticleHandler) GetArticle(c echo.Context) error {
	slug := c.Param("slug")

	article, err := h.articleSvc.GetBySlug(c.Request().Context(), slug)
	if err != nil {
		return handleServiceError(err)
	}

	return c.JSON(http.StatusOK, dto.ArticleResponse{Article: dto.ToArticleBody(article)})
}

// API POST /api/articles/:slug/comments — Add a comment to an article (auth required)
//
// @Summary   Add comment
// @Tags      articles
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param     slug path string           true "Article slug"
// @Param     body body dto.AddCommentRequest true "Comment body"
// @Success   201 {object} dto.CommentResponse
// @Failure   401 {object} map[string]any "Unauthorized"
// @Failure   404 {object} map[string]any "Article not found"
// @Failure   422 {object} map[string]any "Validation error"
// @Router    /articles/{slug}/comments [post]
func (h *ArticleHandler) AddComment(c echo.Context) error {
	var req dto.AddCommentRequest
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	slug := c.Param("slug")
	authorID := middleware.UserIDFromContext(c)

	comment, err := h.commentSvc.AddComment(c.Request().Context(), slug, authorID, req.Comment.Body)
	if err != nil {
		return handleServiceError(err)
	}

	return c.JSON(http.StatusCreated, dto.CommentResponse{Comment: dto.ToCommentBody(comment)})
}

// API GET /api/articles/:slug/comments — Get all comments for an article (public)
//
// @Summary   Get comments
// @Tags      articles
// @Produce   json
// @Param     slug path string true "Article slug"
// @Success   200 {object} dto.CommentsResponse
// @Failure   404 {object} map[string]any "Article not found"
// @Router    /articles/{slug}/comments [get]
func (h *ArticleHandler) GetComments(c echo.Context) error {
	slug := c.Param("slug")

	comments, err := h.commentSvc.GetComments(c.Request().Context(), slug)
	if err != nil {
		return handleServiceError(err)
	}

	bodies := make([]dto.CommentBody, 0, len(comments))
	for _, cm := range comments {
		bodies = append(bodies, dto.ToCommentBody(cm))
	}

	return c.JSON(http.StatusOK, dto.CommentsResponse{Comments: bodies})
}
