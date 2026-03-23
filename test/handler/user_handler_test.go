package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/khanhnp-2797/echo-realworld-api/internal/domain"
	"github.com/khanhnp-2797/echo-realworld-api/internal/dto"
	"github.com/khanhnp-2797/echo-realworld-api/internal/handler"
	"github.com/khanhnp-2797/echo-realworld-api/pkg/apperrors"
)

// ─── helpers ─────────────────────────────────────────────────────────────────

func sampleUser(id uint, username, email string) *domain.User {
	return &domain.User{Username: username, Email: email}
}

// ─── Register ────────────────────────────────────────────────────────────────

func TestRegister_Success(t *testing.T) {
	mockSvc := new(MockUserService)
	mockSvc.On("Register", mock.Anything, "alice", "alice@example.com", "secret123").
		Return(sampleUser(1, "alice", "alice@example.com"), "jwt-token", nil)

	body := `{"user":{"username":"alice","email":"alice@example.com","password":"secret123"}}`
	req := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := newEcho()
	c := e.NewContext(req, rec)

	h := handler.NewUserHandler(mockSvc)
	err := h.Register(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp dto.UserResponse
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "alice", resp.User.Username)
	assert.Equal(t, "alice@example.com", resp.User.Email)
	assert.Equal(t, "jwt-token", resp.User.Token)

	mockSvc.AssertExpectations(t)
}

func TestRegister_ValidationError_MissingFields(t *testing.T) {
	mockSvc := new(MockUserService)

	body := `{"user":{"username":"a","email":"not-an-email","password":"short"}}` // validation fails
	req := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := newEcho()
	c := e.NewContext(req, rec)

	h := handler.NewUserHandler(mockSvc)
	err := h.Register(c)

	assert.Error(t, err)
	he, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusUnprocessableEntity, he.Code)

	mockSvc.AssertNotCalled(t, "Register", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestRegister_EmailAlreadyTaken(t *testing.T) {
	mockSvc := new(MockUserService)
	mockSvc.On("Register", mock.Anything, "alice", "alice@example.com", "secret123").
		Return(nil, "", apperrors.ErrEmailTaken)

	body := `{"user":{"username":"alice","email":"alice@example.com","password":"secret123"}}`
	req := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := newEcho()
	c := e.NewContext(req, rec)

	h := handler.NewUserHandler(mockSvc)
	err := h.Register(c)

	assert.Error(t, err)
	he, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusUnprocessableEntity, he.Code)

	mockSvc.AssertExpectations(t)
}

// ─── Login ────────────────────────────────────────────────────────────────────

func TestLogin_Success(t *testing.T) {
	mockSvc := new(MockUserService)
	mockSvc.On("Login", mock.Anything, "alice@example.com", "secret123").
		Return(sampleUser(1, "alice", "alice@example.com"), "jwt-token", nil)

	body := `{"user":{"email":"alice@example.com","password":"secret123"}}`
	req := httptest.NewRequest(http.MethodPost, "/api/users/login", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := newEcho()
	c := e.NewContext(req, rec)

	h := handler.NewUserHandler(mockSvc)
	err := h.Login(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp dto.UserResponse
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "alice", resp.User.Username)
	assert.Equal(t, "jwt-token", resp.User.Token)

	mockSvc.AssertExpectations(t)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	mockSvc := new(MockUserService)
	mockSvc.On("Login", mock.Anything, "alice@example.com", "wrongpass").
		Return(nil, "", apperrors.ErrInvalidCreds)

	body := `{"user":{"email":"alice@example.com","password":"wrongpass"}}`
	req := httptest.NewRequest(http.MethodPost, "/api/users/login", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := newEcho()
	c := e.NewContext(req, rec)

	h := handler.NewUserHandler(mockSvc)
	err := h.Login(c)

	assert.Error(t, err)
	he, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusUnauthorized, he.Code)

	mockSvc.AssertExpectations(t)
}

// ─── GetProfile ───────────────────────────────────────────────────────────────

func TestGetProfile_Success(t *testing.T) {
	mockSvc := new(MockUserService)
	mockSvc.On("GetProfile", mock.Anything, "alice").
		Return(sampleUser(2, "alice", "alice@example.com"), nil)
	mockSvc.On("IsFollowing", mock.Anything, uint(0), uint(0)).
		Return(false)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	e := newEcho()
	c := e.NewContext(req, rec)
	c.SetParamNames("username")
	c.SetParamValues("alice")

	h := handler.NewUserHandler(mockSvc)
	err := h.GetProfile(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp dto.ProfileResponse
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "alice", resp.Profile.Username)
	assert.False(t, resp.Profile.Following)

	mockSvc.AssertExpectations(t)
}

func TestGetProfile_NotFound(t *testing.T) {
	mockSvc := new(MockUserService)
	mockSvc.On("GetProfile", mock.Anything, "ghost").
		Return(nil, apperrors.ErrNotFound)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	e := newEcho()
	c := e.NewContext(req, rec)
	c.SetParamNames("username")
	c.SetParamValues("ghost")

	h := handler.NewUserHandler(mockSvc)
	err := h.GetProfile(c)

	assert.Error(t, err)
	he, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusNotFound, he.Code)

	mockSvc.AssertExpectations(t)
}

// ─── GetCurrentUser ───────────────────────────────────────────────────────────

func TestGetCurrentUser_Success(t *testing.T) {
	mockSvc := new(MockUserService)
	mockSvc.On("GetByID", mock.Anything, uint(42)).
		Return(sampleUser(42, "bob", "bob@example.com"), nil)

	req := httptest.NewRequest(http.MethodGet, "/api/user", nil)
	rec := httptest.NewRecorder()

	e := newEcho()
	c := e.NewContext(req, rec)
	setAuth(c, 42)

	h := handler.NewUserHandler(mockSvc)
	err := h.GetCurrentUser(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp dto.UserResponse
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "bob", resp.User.Username)
	assert.Equal(t, "test-token", resp.User.Token)

	mockSvc.AssertExpectations(t)
}
