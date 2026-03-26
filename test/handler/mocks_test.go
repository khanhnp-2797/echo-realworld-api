package handler_test

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/khanhnp-2797/echo-realworld-api/internal/domain"
	"github.com/khanhnp-2797/echo-realworld-api/internal/repository"
)

// ────────────────────────────────────────────────────────────────────
// MockUserService
// ────────────────────────────────────────────────────────────────────

type MockUserService struct{ mock.Mock }

func (m *MockUserService) Register(ctx context.Context, username, email, password string) (*domain.User, string, error) {
	args := m.Called(ctx, username, email, password)
	if u, ok := args.Get(0).(*domain.User); ok {
		return u, args.String(1), args.Error(2)
	}
	return nil, "", args.Error(2)
}

func (m *MockUserService) Login(ctx context.Context, email, password string) (*domain.User, string, error) {
	args := m.Called(ctx, email, password)
	if u, ok := args.Get(0).(*domain.User); ok {
		return u, args.String(1), args.Error(2)
	}
	return nil, "", args.Error(2)
}

func (m *MockUserService) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	args := m.Called(ctx, id)
	if u, ok := args.Get(0).(*domain.User); ok {
		return u, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserService) GetProfile(ctx context.Context, username string) (*domain.User, error) {
	args := m.Called(ctx, username)
	if u, ok := args.Get(0).(*domain.User); ok {
		return u, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserService) Follow(ctx context.Context, followerID uint, username string) (*domain.User, error) {
	args := m.Called(ctx, followerID, username)
	if u, ok := args.Get(0).(*domain.User); ok {
		return u, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserService) Unfollow(ctx context.Context, followerID uint, username string) (*domain.User, error) {
	args := m.Called(ctx, followerID, username)
	if u, ok := args.Get(0).(*domain.User); ok {
		return u, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserService) IsFollowing(ctx context.Context, followerID, followedID uint) bool {
	args := m.Called(ctx, followerID, followedID)
	return args.Bool(0)
}

// ────────────────────────────────────────────────────────────────────
// MockArticleService
// ────────────────────────────────────────────────────────────────────

type MockArticleService struct{ mock.Mock }

func (m *MockArticleService) GetBySlug(ctx context.Context, slug string) (*domain.Article, error) {
	args := m.Called(ctx, slug)
	if a, ok := args.Get(0).(*domain.Article); ok {
		return a, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockArticleService) List(ctx context.Context, filter repository.ArticleFilter) ([]*domain.Article, int64, error) {
	args := m.Called(ctx, filter)
	if a, ok := args.Get(0).([]*domain.Article); ok {
		return a, args.Get(1).(int64), args.Error(2)
	}
	return nil, 0, args.Error(2)
}

func (m *MockArticleService) Feed(ctx context.Context, userID uint, filter repository.ArticleFilter) ([]*domain.Article, int64, error) {
	args := m.Called(ctx, userID, filter)
	if a, ok := args.Get(0).([]*domain.Article); ok {
		return a, args.Get(1).(int64), args.Error(2)
	}
	return nil, 0, args.Error(2)
}

func (m *MockArticleService) Create(ctx context.Context, authorID uint, title, description, body string, tagList []string) (*domain.Article, error) {
	args := m.Called(ctx, authorID, title, description, body, tagList)
	if a, ok := args.Get(0).(*domain.Article); ok {
		return a, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockArticleService) Update(ctx context.Context, slug string, authorID uint, title, description, body *string) (*domain.Article, error) {
	args := m.Called(ctx, slug, authorID, title, description, body)
	if a, ok := args.Get(0).(*domain.Article); ok {
		return a, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockArticleService) DeleteBySlug(ctx context.Context, slug string, authorID uint) error {
	args := m.Called(ctx, slug, authorID)
	return args.Error(0)
}

func (m *MockArticleService) Favorite(ctx context.Context, userID uint, slug string) (*domain.Article, error) {
	args := m.Called(ctx, userID, slug)
	if a, ok := args.Get(0).(*domain.Article); ok {
		return a, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockArticleService) Unfavorite(ctx context.Context, userID uint, slug string) (*domain.Article, error) {
	args := m.Called(ctx, userID, slug)
	if a, ok := args.Get(0).(*domain.Article); ok {
		return a, args.Error(1)
	}
	return nil, args.Error(1)
}

// ────────────────────────────────────────────────────────────────────
// MockCommentService
// ────────────────────────────────────────────────────────────────────

type MockCommentService struct{ mock.Mock }

func (m *MockCommentService) AddComment(ctx context.Context, slug string, authorID uint, body string) (*domain.Comment, error) {
	args := m.Called(ctx, slug, authorID, body)
	if c, ok := args.Get(0).(*domain.Comment); ok {
		return c, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockCommentService) GetComments(ctx context.Context, slug string) ([]*domain.Comment, error) {
	args := m.Called(ctx, slug)
	if c, ok := args.Get(0).([]*domain.Comment); ok {
		return c, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockCommentService) DeleteComment(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ────────────────────────────────────────────────────────────────────
// MockTagService
// ────────────────────────────────────────────────────────────────────

type MockTagService struct{ mock.Mock }

func (m *MockTagService) GetAll(ctx context.Context) ([]*domain.Tag, error) {
	args := m.Called(ctx)
	if t, ok := args.Get(0).([]*domain.Tag); ok {
		return t, args.Error(1)
	}
	return nil, args.Error(1)
}
