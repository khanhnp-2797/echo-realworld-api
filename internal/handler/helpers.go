package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"github.com/khanhnp-2797/echo-realworld-api/pkg/apperrors"
)

// bindAndValidate binds the request body and runs validation.
func bindAndValidate(c echo.Context, v any) error {
	if err := c.Bind(v); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	if err := c.Validate(v); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			msgs := make(map[string][]string)
			for _, fe := range ve {
				field := fe.Field()
				msgs[field] = append(msgs[field], fe.Tag())
			}
			return echo.NewHTTPError(http.StatusUnprocessableEntity, map[string]any{"errors": msgs})
		}
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	return nil
}

// handleServiceError maps domain/app errors to HTTP errors.
func handleServiceError(err error) *echo.HTTPError {
	switch {
	case errors.Is(err, apperrors.ErrNotFound):
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	case errors.Is(err, apperrors.ErrUnauthorized), errors.Is(err, apperrors.ErrInvalidCreds):
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	case errors.Is(err, apperrors.ErrForbidden):
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	case errors.Is(err, apperrors.ErrEmailTaken), errors.Is(err, apperrors.ErrUsernameTaken):
		return echo.NewHTTPError(http.StatusUnprocessableEntity, map[string]any{
			"errors": map[string]any{"body": []string{err.Error()}},
		})
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
}

// queryInt reads an integer query param, returning def when absent or invalid.
func queryInt(c echo.Context, key string, def int) int {
	v, err := strconv.Atoi(c.QueryParam(key))
	if err != nil || v < 0 {
		return def
	}
	return v
}
