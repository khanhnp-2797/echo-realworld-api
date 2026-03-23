package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/khanhnp-2797/echo-realworld-api/internal/domain"
	"github.com/khanhnp-2797/echo-realworld-api/internal/dto"
	"github.com/khanhnp-2797/echo-realworld-api/internal/handler"
)

func TestListTags_Success(t *testing.T) {
	mockSvc := new(MockTagService)
	mockSvc.On("GetAll", mock.Anything).
		Return([]*domain.Tag{{Name: "go"}, {Name: "echo"}}, nil)

	e := newEcho()
	req := httptest.NewRequest(http.MethodGet, "/api/tags", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := handler.NewTagHandler(mockSvc)
	err := h.ListTags(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var body dto.TagsResponse
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, []string{"go", "echo"}, body.Tags)

	mockSvc.AssertExpectations(t)
}

func TestListTags_EmptyResult(t *testing.T) {
	mockSvc := new(MockTagService)
	mockSvc.On("GetAll", mock.Anything).
		Return([]*domain.Tag{}, nil)

	e := newEcho()
	req := httptest.NewRequest(http.MethodGet, "/api/tags", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := handler.NewTagHandler(mockSvc)
	err := h.ListTags(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var body dto.TagsResponse
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Empty(t, body.Tags)

	mockSvc.AssertExpectations(t)
}

func TestListTags_ServiceError(t *testing.T) {
	mockSvc := new(MockTagService)
	mockSvc.On("GetAll", mock.Anything).
		Return(nil, errors.New("db down"))

	e := newEcho()
	req := httptest.NewRequest(http.MethodGet, "/api/tags", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := handler.NewTagHandler(mockSvc)
	err := h.ListTags(c)

	assert.Error(t, err)
	he, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusInternalServerError, he.Code)

	mockSvc.AssertExpectations(t)
}
