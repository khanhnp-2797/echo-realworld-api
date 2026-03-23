package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/khanhnp-2797/echo-realworld-api/internal/domain"
	"github.com/khanhnp-2797/echo-realworld-api/internal/dto"
	"github.com/khanhnp-2797/echo-realworld-api/internal/handler"
	"github.com/khanhnp-2797/echo-realworld-api/internal/repository"
	"github.com/khanhnp-2797/echo-realworld-api/pkg/apperrors"
)

// ─── helpers ──────────────────────────────────────────────────────────────────

func sampleArticle(slug, title string, author domain.User) *domain.Article {
	return &domain.Article{
		Slug:        slug,
		Title:       title,
		Description: "desc",
		Body:        "body",
		Author:      author,
		Tags:        []*domain.Tag{},
		FavoritedBy: []*domain.User{},
	}
}

// ─── ListArticles ─────────────────────────────────────────────────────────────

func TestListArticles_Success(t *testing.T) {
	author := domain.User{Username: "alice", Email: "alice@example.com"}
	articles := []*domain.Article{
		sampleArticle("how-to-go", "How To Go", author),
		sampleArticle("echo-rocks", "Echo Rocks", author),
	}

	mockArticleSvc := new(MockArticleService)
	mockArticleSvc.On("List", mock.Anything, repository.ArticleFilter{Limit: 20, Offset: 0}).
		Return(articles, int64(2), nil)

	mockUserSvc := new(MockUserService)
	mockUserSvc.On("IsFollowing", mock.Anything, uint(0), uint(0)).Return(false)

	req := httptest.NewRequest(http.MethodGet, "/api/articles", nil)
	rec := httptest.NewRecorder()

	e := newEcho()
	c := e.NewContext(req, rec)

	h := handler.NewArticleHandler(mockArticleSvc, new(MockCommentService), mockUserSvc)
	err := h.ListArticles(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp dto.ArticlesResponse
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, int64(2), resp.ArticlesCount)
	assert.Len(t, resp.Articles, 2)
	assert.Equal(t, "how-to-go", resp.Articles[0].Slug)

	mockArticleSvc.AssertExpectations(t)
}

func TestListArticles_FilterByTag(t *testing.T) {
	author := domain.User{Username: "alice"}
	articles := []*domain.Article{sampleArticle("tagged", "Tagged", author)}

	mockArticleSvc := new(MockArticleService)
	mockArticleSvc.On("List", mock.Anything, repository.ArticleFilter{Tag: "go", Limit: 20, Offset: 0}).
		Return(articles, int64(1), nil)

	mockUserSvc := new(MockUserService)
	mockUserSvc.On("IsFollowing", mock.Anything, uint(0), uint(0)).Return(false)

	req := httptest.NewRequest(http.MethodGet, "/api/articles?tag=go", nil)
	rec := httptest.NewRecorder()

	e := newEcho()
	c := e.NewContext(req, rec)

	h := handler.NewArticleHandler(mockArticleSvc, new(MockCommentService), mockUserSvc)
	err := h.ListArticles(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp dto.ArticlesResponse
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, int64(1), resp.ArticlesCount)

	mockArticleSvc.AssertExpectations(t)
}

func TestListArticles_Empty(t *testing.T) {
	mockArticleSvc := new(MockArticleService)
	mockArticleSvc.On("List", mock.Anything, repository.ArticleFilter{Limit: 20, Offset: 0}).
		Return([]*domain.Article{}, int64(0), nil)

	mockUserSvc := new(MockUserService)

	req := httptest.NewRequest(http.MethodGet, "/api/articles", nil)
	rec := httptest.NewRecorder()

	e := newEcho()
	c := e.NewContext(req, rec)

	h := handler.NewArticleHandler(mockArticleSvc, new(MockCommentService), mockUserSvc)
	err := h.ListArticles(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp dto.ArticlesResponse
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, int64(0), resp.ArticlesCount)
	assert.Empty(t, resp.Articles)

	mockArticleSvc.AssertExpectations(t)
}

// ─── GetArticle ───────────────────────────────────────────────────────────────

func TestGetArticle_Success(t *testing.T) {
	author := domain.User{Username: "alice"}
	article := sampleArticle("how-to-go", "How To Go", author)

	mockArticleSvc := new(MockArticleService)
	mockArticleSvc.On("GetBySlug", mock.Anything, "how-to-go").Return(article, nil)

	mockUserSvc := new(MockUserService)
	mockUserSvc.On("IsFollowing", mock.Anything, uint(0), uint(0)).Return(false)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	e := newEcho()
	c := e.NewContext(req, rec)
	c.SetParamNames("slug")
	c.SetParamValues("how-to-go")

	h := handler.NewArticleHandler(mockArticleSvc, new(MockCommentService), mockUserSvc)
	err := h.GetArticle(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp dto.ArticleResponse
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "how-to-go", resp.Article.Slug)
	assert.Equal(t, "How To Go", resp.Article.Title)
	assert.Equal(t, "alice", resp.Article.Author.Username)

	mockArticleSvc.AssertExpectations(t)
}

func TestGetArticle_NotFound(t *testing.T) {
	mockArticleSvc := new(MockArticleService)
	mockArticleSvc.On("GetBySlug", mock.Anything, "missing").Return(nil, apperrors.ErrNotFound)

	mockUserSvc := new(MockUserService)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	e := newEcho()
	c := e.NewContext(req, rec)
	c.SetParamNames("slug")
	c.SetParamValues("missing")

	h := handler.NewArticleHandler(mockArticleSvc, new(MockCommentService), mockUserSvc)
	err := h.GetArticle(c)

	assert.Error(t, err)
	he, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusNotFound, he.Code)

	mockArticleSvc.AssertExpectations(t)
}

// ─── Feed ─────────────────────────────────────────────────────────────────────

func TestFeed_Success(t *testing.T) {
	author := domain.User{Username: "bob"}
	articles := []*domain.Article{sampleArticle("bobs-post", "Bob's Post", author)}

	mockArticleSvc := new(MockArticleService)
	mockArticleSvc.On("Feed", mock.Anything, uint(1), repository.ArticleFilter{Limit: 20, Offset: 0}).
		Return(articles, int64(1), nil)

	req := httptest.NewRequest(http.MethodGet, "/api/articles/feed", nil)
	rec := httptest.NewRecorder()

	e := newEcho()
	c := e.NewContext(req, rec)
	setAuth(c, 1)

	h := handler.NewArticleHandler(mockArticleSvc, new(MockCommentService), new(MockUserService))
	err := h.Feed(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp dto.ArticlesResponse
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, int64(1), resp.ArticlesCount)
	assert.Equal(t, "bobs-post", resp.Articles[0].Slug)

	mockArticleSvc.AssertExpectations(t)
}
