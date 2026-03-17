package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/khanhnp-2797/echo-realworld-api/internal/dto"
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

// ──────────────────────────── Handlers ────────────────────────────

// API POST /api/users — Register a new user (public)
//
// @Summary  Register a new user
// @Tags     users
// @Accept   json
// @Produce  json
// @Param    body body dto.RegisterRequest true "User credentials"
// @Success  201 {object} dto.UserResponse
// @Failure  422 {object} map[string]any "Validation error"
// @Failure  500 {object} map[string]any "Internal server error"
// @Router   /users [post]
func (h *UserHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	user, token, err := h.userSvc.Register(c.Request().Context(),
		req.User.Username, req.User.Email, req.User.Password)
	if err != nil {
		return handleServiceError(err)
	}

	return c.JSON(http.StatusCreated, dto.UserResponse{User: dto.ToUserBody(user, token)})
}

// API POST /api/users/login — Authenticate and receive a JWT token (public)
//
// @Summary  Login and receive a JWT token
// @Tags     users
// @Accept   json
// @Produce  json
// @Param    body body dto.LoginRequest true "Login credentials"
// @Success  200 {object} dto.UserResponse
// @Failure  401 {object} map[string]any "Invalid credentials"
// @Failure  422 {object} map[string]any "Validation error"
// @Router   /users/login [post]
func (h *UserHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	user, token, err := h.userSvc.Login(c.Request().Context(), req.User.Email, req.User.Password)
	if err != nil {
		return handleServiceError(err)
	}

	return c.JSON(http.StatusOK, dto.UserResponse{User: dto.ToUserBody(user, token)})
}

// API GET /api/user — Get the currently authenticated user (auth required)
//
// @Summary   Get current user
// @Tags      users
// @Security  BearerAuth
// @Produce   json
// @Success   200 {object} dto.UserResponse
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

	return c.JSON(http.StatusOK, dto.UserResponse{User: dto.ToUserBody(user, token)})
}

// API GET /api/profiles/:username — Get a user's public profile (public)
//
// @Summary   Get profile
// @Tags      profiles
// @Produce   json
// @Param     username path string true "Username"
// @Success   200 {object} dto.ProfileResponse
// @Failure   404 {object} map[string]any "User not found"
// @Router    /profiles/{username} [get]
func (h *UserHandler) GetProfile(c echo.Context) error {
	username := c.Param("username")

	user, err := h.userSvc.GetProfile(c.Request().Context(), username)
	if err != nil {
		return handleServiceError(err)
	}

	return c.JSON(http.StatusOK, dto.ProfileResponse{Profile: dto.ToProfileBody(user)})
}
