package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/khanhnp-2797/echo-realworld-api/internal/domain"
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

// ──────────────────────────── Request / Response DTOs ────────────────────────────

type addCommentRequest struct {
	Comment struct {
		Body string `json:"body" validate:"required"`
	} `json:"comment"`
}

type articleResponse struct {
	Article articleBody `json:"article"`
}

type articlesResponse struct {
	Articles      []articleBody `json:"articles"`
	ArticlesCount int64         `json:"articlesCount"`
}

// articleBody is the public DTO for an Article (no sensitive fields).
type articleBody struct {
	Slug        string      `json:"slug"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Body        string      `json:"body"`
	TagList     []string    `json:"tagList"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
	Author      profileBody `json:"author"`
}

type commentResponse struct {
	Comment commentBody `json:"comment"`
}

type commentsResponse struct {
	Comments []commentBody `json:"comments"`
}

type commentBody struct {
	ID        uint        `json:"id"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
	Body      string      `json:"body"`
	Author    profileBody `json:"author"` // Eager-loaded via GORM Preload
}

// ──────────────────────────── Mappers ────────────────────────────

func toArticleBody(a *domain.Article) articleBody {
	tags := make([]string, 0, len(a.Tags))
	for _, t := range a.Tags {
		tags = append(tags, t.Name)
	}
	return articleBody{
		Slug:        a.Slug,
		Title:       a.Title,
		Description: a.Description,
		Body:        a.Body,
		TagList:     tags,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
		Author:      toProfileBody(&a.Author),
	}
}

func toCommentBody(c *domain.Comment) commentBody {
	return commentBody{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Body:      c.Body,
		Author:    toProfileBody(&c.Author),
	}
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
// @Success   200 {object} articlesResponse
// @Failure   500 {object} map[string]any "Internal server error"
// @Router    /articles [get]
func (h *ArticleHandler) ListArticles(c echo.Context) error {
	// c.Bind() reads query params declared via struct tags
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

	bodies := make([]articleBody, 0, len(articles))
	for _, a := range articles {
		bodies = append(bodies, toArticleBody(a))
	}

	// c.JSON() serialises the DTO and writes Content-Type: application/json
	return c.JSON(http.StatusOK, articlesResponse{Articles: bodies, ArticlesCount: count})
}

// API GET /api/articles/:slug — Get a single article by slug (public)
//
// @Summary   Get article
// @Tags      articles
// @Produce   json
// @Param     slug path string true "Article slug"
// @Success   200 {object} articleResponse
// @Failure   404 {object} map[string]any "Article not found"
// @Router    /articles/{slug} [get]
func (h *ArticleHandler) GetArticle(c echo.Context) error {
	slug := c.Param("slug") // path param via Echo routing :slug

	article, err := h.articleSvc.GetBySlug(c.Request().Context(), slug)
	if err != nil {
		return handleServiceError(err)
	}

	return c.JSON(http.StatusOK, articleResponse{Article: toArticleBody(article)})
}

// API POST /api/articles/:slug/comments — Add a comment to an article (auth required)
//
// @Summary   Add comment
// @Tags      articles
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param     slug path string           true "Article slug"
// @Param     body body addCommentRequest true "Comment body"
// @Success   201 {object} commentResponse
// @Failure   401 {object} map[string]any "Unauthorized"
// @Failure   404 {object} map[string]any "Article not found"
// @Failure   422 {object} map[string]any "Validation error"
// @Router    /articles/{slug}/comments [post]
func (h *ArticleHandler) AddComment(c echo.Context) error {
	var req addCommentRequest
	// bindAndValidate: c.Bind() + c.Validate() in one step
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	slug := c.Param("slug")
	authorID := middleware.UserIDFromContext(c)

	comment, err := h.commentSvc.AddComment(c.Request().Context(), slug, authorID, req.Comment.Body)
	if err != nil {
		return handleServiceError(err)
	}

	return c.JSON(http.StatusCreated, commentResponse{Comment: toCommentBody(comment)})
}

// API GET /api/articles/:slug/comments — Get all comments for an article (public)
//
// @Summary   Get comments
// @Tags      articles
// @Produce   json
// @Param     slug path string true "Article slug"
// @Success   200 {object} commentsResponse
// @Failure   404 {object} map[string]any "Article not found"
// @Router    /articles/{slug}/comments [get]
func (h *ArticleHandler) GetComments(c echo.Context) error {
	slug := c.Param("slug")

	comments, err := h.commentSvc.GetComments(c.Request().Context(), slug)
	if err != nil {
		return handleServiceError(err)
	}

	bodies := make([]commentBody, 0, len(comments))
	for _, cm := range comments {
		bodies = append(bodies, toCommentBody(cm))
	}

	return c.JSON(http.StatusOK, commentsResponse{Comments: bodies})
}

// ──────────────────────────── Utility ────────────────────────────

func queryInt(c echo.Context, key string, def int) int {
	v, err := strconv.Atoi(c.QueryParam(key))
	if err != nil || v < 0 {
		return def
	}
	return v
}

// ArticleHandler handles all article and comment endpoints.
