package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/khanhnp-2797/echo-realworld-api/internal/domain"
	"github.com/khanhnp-2797/echo-realworld-api/internal/middleware"
	"github.com/khanhnp-2797/echo-realworld-api/internal/service"
)

// UserHandler handles user auth and profile endpoints.
type UserHandler struct {
	userSvc service.UserService
}

func NewUserHandler(userSvc service.UserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

// ──────────────────────────── Request / Response DTOs ────────────────────────────

type registerRequest struct {
	User struct {
		Username string `json:"username" validate:"required,min=2,max=100"`
		Email    string `json:"email"    validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	} `json:"user"`
}

type loginRequest struct {
	User struct {
		Email    string `json:"email"    validate:"required,email"`
		Password string `json:"password" validate:"required"`
	} `json:"user"`
}

type userResponse struct {
	User userBody `json:"user"`
}

// userBody is the public DTO for a User — password is never included.
type userBody struct {
	Email    string  `json:"email"`
	Token    string  `json:"token"`
	Username string  `json:"username"`
	Bio      *string `json:"bio"`
	Image    *string `json:"image"`
}

type profileResponse struct {
	Profile profileBody `json:"profile"`
}

// profileBody is the public DTO for a User profile.
type profileBody struct {
	Username string  `json:"username"`
	Bio      *string `json:"bio"`
	Image    *string `json:"image"`
}

// ──────────────────────────── Mappers ────────────────────────────

func toUserBody(u *domain.User, token string) userBody {
	return userBody{
		Email:    u.Email,
		Token:    token,
		Username: u.Username,
		Bio:      u.Bio,
		Image:    u.Image,
	}
}

func toProfileBody(u *domain.User) profileBody {
	return profileBody{
		Username: u.Username,
		Bio:      u.Bio,
		Image:    u.Image,
	}
}

// ──────────────────────────── Handlers ────────────────────────────

// API POST /api/users — Register a new user (public)
//
// @Summary  Register a new user
// @Tags     users
// @Accept   json
// @Produce  json
// @Param    body body registerRequest true "User credentials"
// @Success  201 {object} userResponse
// @Failure  422 {object} map[string]any "Validation error"
// @Failure  500 {object} map[string]any "Internal server error"
// @Router   /users [post]
func (h *UserHandler) Register(c echo.Context) error {
	var req registerRequest
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	user, token, err := h.userSvc.Register(c.Request().Context(),
		req.User.Username, req.User.Email, req.User.Password)
	if err != nil {
		return handleServiceError(err)
	}

	return c.JSON(http.StatusCreated, userResponse{User: toUserBody(user, token)})
}

// API POST /api/users/login — Authenticate and receive a JWT token (public)
//
// @Summary  Login and receive a JWT token
// @Tags     users
// @Accept   json
// @Produce  json
// @Param    body body loginRequest true "Login credentials"
// @Success  200 {object} userResponse
// @Failure  401 {object} map[string]any "Invalid credentials"
// @Failure  422 {object} map[string]any "Validation error"
// @Router   /users/login [post]
func (h *UserHandler) Login(c echo.Context) error {
	var req loginRequest
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	user, token, err := h.userSvc.Login(c.Request().Context(), req.User.Email, req.User.Password)
	if err != nil {
		return handleServiceError(err)
	}

	return c.JSON(http.StatusOK, userResponse{User: toUserBody(user, token)})
}

// API GET /api/user — Get the currently authenticated user (auth required)
//
// @Summary   Get current user
// @Tags      users
// @Security  BearerAuth
// @Produce   json
// @Success   200 {object} userResponse
// @Failure   401 {object} map[string]any "Unauthorized"
// @Failure   404 {object} map[string]any "User not found"
// @Router    /user [get]
func (h *UserHandler) GetCurrentUser(c echo.Context) error {
	userID := middleware.UserIDFromContext(c)
	token := middleware.TokenFromContext(c)

	user, err := h.userSvc.GetByID(c.Request().Context(), userID)
	if err != nil {
		return handleServiceError(err)
	}

	return c.JSON(http.StatusOK, userResponse{User: toUserBody(user, token)})
}

// API GET /api/profiles/:username — Get a user's public profile (public)
//
// @Summary   Get profile
// @Tags      profiles
// @Produce   json
// @Param     username path string true "Username"
// @Success   200 {object} profileResponse
// @Failure   404 {object} map[string]any "User not found"
// @Router    /profiles/{username} [get]
func (h *UserHandler) GetProfile(c echo.Context) error {
	username := c.Param("username")

	user, err := h.userSvc.GetProfile(c.Request().Context(), username)
	if err != nil {
		return handleServiceError(err)
	}

	return c.JSON(http.StatusOK, profileResponse{Profile: toProfileBody(user)})
}

// UserHandler handles all user & profile endpoints.
