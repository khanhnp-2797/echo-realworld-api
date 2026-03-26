package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/khanhnp-2797/echo-realworld-api/internal/dto"
	"github.com/khanhnp-2797/echo-realworld-api/internal/middleware"
	"github.com/khanhnp-2797/echo-realworld-api/internal/repository"
	"github.com/khanhnp-2797/echo-realworld-api/internal/service"
	"github.com/khanhnp-2797/echo-realworld-api/internal/ws"
)

// ArticleHandler handles article and comment endpoints.
type ArticleHandler struct {
	articleSvc service.ArticleService
	commentSvc service.CommentService
	userSvc    service.UserService // needed for following checks
	hub        *ws.Hub             // nil = WebSocket broadcast disabled
}

func NewArticleHandler(
	articleSvc service.ArticleService,
	commentSvc service.CommentService,
	userSvc service.UserService,
) *ArticleHandler {
	return &ArticleHandler{
		articleSvc: articleSvc,
		commentSvc: commentSvc,
		userSvc:    userSvc,
	}
}

// SetHub attaches the WebSocket hub so AddComment can broadcast events.
func (h *ArticleHandler) SetHub(hub *ws.Hub) { h.hub = hub }

// ──────────────────────────── Handlers ────────────────────────────

// API POST /api/articles — Create a new article (auth required)
//
// @Summary   Create article
// @Tags      articles
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param     body body dto.CreateArticleRequest true "Article data"
// @Success   201 {object} dto.ArticleResponse
// @Failure   401 {object} map[string]any "Unauthorized"
// @Failure   422 {object} map[string]any "Validation error"
// @Router    /articles [post]
func (h *ArticleHandler) CreateArticle(c echo.Context) error {
	var req dto.CreateArticleRequest
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	authorID := middleware.UserIDFromContext(c)
	article, err := h.articleSvc.Create(
		c.Request().Context(),
		authorID,
		req.Article.Title,
		req.Article.Description,
		req.Article.Body,
		req.Article.TagList,
	)
	if err != nil {
		return handleServiceError(err)
	}

	following := h.userSvc.IsFollowing(c.Request().Context(), authorID, article.AuthorID)
	return c.JSON(http.StatusCreated, dto.ArticleResponse{Article: dto.ToArticleBody(article, authorID, following)})
}

// API PUT /api/articles/:slug — Update own article (auth required)
//
// @Summary   Update article
// @Tags      articles
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param     slug path string                  true "Article slug"
// @Param     body body dto.UpdateArticleRequest true "Fields to update"
// @Success   200 {object} dto.ArticleResponse
// @Failure   401 {object} map[string]any "Unauthorized"
// @Failure   403 {object} map[string]any "Forbidden"
// @Failure   404 {object} map[string]any "Not found"
// @Router    /articles/{slug} [put]
func (h *ArticleHandler) UpdateArticle(c echo.Context) error {
	var req dto.UpdateArticleRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	slug := c.Param("slug")
	authorID := middleware.UserIDFromContext(c)

	article, err := h.articleSvc.Update(
		c.Request().Context(),
		slug,
		authorID,
		req.Article.Title,
		req.Article.Description,
		req.Article.Body,
	)
	if err != nil {
		return handleServiceError(err)
	}

	following := h.userSvc.IsFollowing(c.Request().Context(), authorID, article.AuthorID)
	return c.JSON(http.StatusOK, dto.ArticleResponse{Article: dto.ToArticleBody(article, authorID, following)})
}

// API DELETE /api/articles/:slug — Delete own article (auth required)
//
// @Summary   Delete article
// @Tags      articles
// @Security  BearerAuth
// @Param     slug path string true "Article slug"
// @Success   204 "No Content"
// @Failure   401 {object} map[string]any "Unauthorized"
// @Failure   403 {object} map[string]any "Forbidden"
// @Failure   404 {object} map[string]any "Not found"
// @Router    /articles/{slug} [delete]
func (h *ArticleHandler) DeleteArticle(c echo.Context) error {
	slug := c.Param("slug")
	authorID := middleware.UserIDFromContext(c)

	if err := h.articleSvc.DeleteBySlug(c.Request().Context(), slug, authorID); err != nil {
		return handleServiceError(err)
	}
	return c.NoContent(http.StatusNoContent)
}

// API GET /api/articles — List articles with optional filters (public)
//
// @Summary   List articles
// @Tags      articles
// @Produce   json
// @Param     tag       query string false "Filter by tag"
// @Param     author    query string false "Filter by author username"
// @Param     favorited query string false "Filter by username who favorited"
// @Param     limit     query int    false "Limit (default 20)"
// @Param     offset    query int    false "Offset (default 0)"
// @Success   200 {object} dto.ArticlesResponse
// @Failure   500 {object} map[string]any "Internal server error"
// @Router    /articles [get]
func (h *ArticleHandler) ListArticles(c echo.Context) error {
	viewerID := middleware.OptionalUserIDFromContext(c)
	filter := repository.ArticleFilter{
		Tag:       c.QueryParam("tag"),
		Author:    c.QueryParam("author"),
		Favorited: c.QueryParam("favorited"),
		Limit:     queryInt(c, "limit", 20),
		Offset:    queryInt(c, "offset", 0),
	}

	articles, count, err := h.articleSvc.List(c.Request().Context(), filter)
	if err != nil {
		return handleServiceError(err)
	}

	bodies := make([]dto.ArticleBody, 0, len(articles))
	for _, a := range articles {
		following := h.userSvc.IsFollowing(c.Request().Context(), viewerID, a.AuthorID)
		bodies = append(bodies, dto.ToArticleBody(a, viewerID, following))
	}

	return c.JSON(http.StatusOK, dto.ArticlesResponse{Articles: bodies, ArticlesCount: count})
}

// API GET /api/articles/feed — Feed from followed users (auth required)
//
// @Summary   Get feed
// @Tags      articles
// @Security  BearerAuth
// @Produce   json
// @Param     limit  query int false "Limit (default 20)"
// @Param     offset query int false "Offset (default 0)"
// @Success   200 {object} dto.ArticlesResponse
// @Failure   401 {object} map[string]any "Unauthorized"
// @Router    /articles/feed [get]
func (h *ArticleHandler) Feed(c echo.Context) error {
	viewerID := middleware.UserIDFromContext(c)
	filter := repository.ArticleFilter{
		Limit:  queryInt(c, "limit", 20),
		Offset: queryInt(c, "offset", 0),
	}

	articles, count, err := h.articleSvc.Feed(c.Request().Context(), viewerID, filter)
	if err != nil {
		return handleServiceError(err)
	}

	bodies := make([]dto.ArticleBody, 0, len(articles))
	for _, a := range articles {
		// All articles in the feed are authored by followed users → following=true
		bodies = append(bodies, dto.ToArticleBody(a, viewerID, true))
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
	viewerID := middleware.OptionalUserIDFromContext(c)

	article, err := h.articleSvc.GetBySlug(c.Request().Context(), slug)
	if err != nil {
		return handleServiceError(err)
	}

	following := h.userSvc.IsFollowing(c.Request().Context(), viewerID, article.AuthorID)
	return c.JSON(http.StatusOK, dto.ArticleResponse{Article: dto.ToArticleBody(article, viewerID, following)})
}

// API POST /api/articles/:slug/favorite — Favorite an article (auth required)
//
// @Summary   Favorite article
// @Tags      articles
// @Security  BearerAuth
// @Produce   json
// @Param     slug path string true "Article slug"
// @Success   200 {object} dto.ArticleResponse
// @Failure   401 {object} map[string]any "Unauthorized"
// @Failure   404 {object} map[string]any "Article not found"
// @Router    /articles/{slug}/favorite [post]
func (h *ArticleHandler) FavoriteArticle(c echo.Context) error {
	slug := c.Param("slug")
	viewerID := middleware.UserIDFromContext(c)

	article, err := h.articleSvc.Favorite(c.Request().Context(), viewerID, slug)
	if err != nil {
		return handleServiceError(err)
	}

	following := h.userSvc.IsFollowing(c.Request().Context(), viewerID, article.AuthorID)
	return c.JSON(http.StatusOK, dto.ArticleResponse{Article: dto.ToArticleBody(article, viewerID, following)})
}

// API DELETE /api/articles/:slug/favorite — Unfavorite an article (auth required)
//
// @Summary   Unfavorite article
// @Tags      articles
// @Security  BearerAuth
// @Produce   json
// @Param     slug path string true "Article slug"
// @Success   200 {object} dto.ArticleResponse
// @Failure   401 {object} map[string]any "Unauthorized"
// @Failure   404 {object} map[string]any "Article not found"
// @Router    /articles/{slug}/favorite [delete]
func (h *ArticleHandler) UnfavoriteArticle(c echo.Context) error {
	slug := c.Param("slug")
	viewerID := middleware.UserIDFromContext(c)

	article, err := h.articleSvc.Unfavorite(c.Request().Context(), viewerID, slug)
	if err != nil {
		return handleServiceError(err)
	}

	following := h.userSvc.IsFollowing(c.Request().Context(), viewerID, article.AuthorID)
	return c.JSON(http.StatusOK, dto.ArticleResponse{Article: dto.ToArticleBody(article, viewerID, following)})
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

	body := dto.ToCommentBody(comment, false)

	// Broadcast to all WebSocket subscribers watching this article's comments.
	if h.hub != nil {
		type wsEvent struct {
			Type    string          `json:"type"`
			Comment dto.CommentBody `json:"comment"`
		}
		if payload, err := json.Marshal(wsEvent{Type: "new_comment", Comment: body}); err == nil {
			log.Printf("[ws] broadcasting new_comment slug=%s", slug)
			h.hub.Broadcast(slug, payload)
		} else {
			log.Printf("[ws] failed to marshal comment event: %v", err)
		}
	} else {
		log.Printf("[ws] hub is nil, skipping broadcast")
	}

	// Comment author is the current user — viewers don't follow themselves
	return c.JSON(http.StatusCreated, dto.CommentResponse{Comment: body})
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
	viewerID := middleware.OptionalUserIDFromContext(c)

	comments, err := h.commentSvc.GetComments(c.Request().Context(), slug)
	if err != nil {
		return handleServiceError(err)
	}

	bodies := make([]dto.CommentBody, 0, len(comments))
	for _, cm := range comments {
		following := h.userSvc.IsFollowing(c.Request().Context(), viewerID, cm.AuthorID)
		bodies = append(bodies, dto.ToCommentBody(cm, following))
	}

	return c.JSON(http.StatusOK, dto.CommentsResponse{Comments: bodies})
}

// API DELETE /api/articles/:slug/comments/:id — Delete own comment (auth + owner required)
//
// @Summary   Delete comment
// @Tags      articles
// @Security  BearerAuth
// @Param     slug path string true "Article slug"
// @Param     id   path int    true "Comment ID"
// @Success   204
// @Failure   401 {object} map[string]any "Unauthorized"
// @Failure   403 {object} map[string]any "Forbidden"
// @Failure   404 {object} map[string]any "Comment not found"
// @Router    /articles/{slug}/comments/{id} [delete]
func (h *ArticleHandler) DeleteComment(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid comment id")
	}

	// Ownership already verified by CommentOwner middleware
	if err := h.commentSvc.DeleteComment(c.Request().Context(), uint(id)); err != nil {
		return handleServiceError(err)
	}

	return c.NoContent(http.StatusNoContent)
}
