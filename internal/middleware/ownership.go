package middleware

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/khanhnp-2797/echo-realworld-api/internal/repository"
)

// CommentOwner is a middleware that verifies the authenticated user is the
// author of the comment identified by the :id path parameter.
// It short-circuits with 403 Forbidden if ownership check fails, or 404 if
// the comment does not exist.
// The fetched comment is stored in context under the key "comment" so the
// handler can reuse it without an extra DB round-trip.
func CommentOwner(commentRepo repository.CommentRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			idStr := c.Param("id")
			id, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "invalid comment id")
			}

			comment, err := commentRepo.FindByID(c.Request().Context(), uint(id))
			if err != nil {
				return echo.NewHTTPError(http.StatusNotFound, "comment not found")
			}

			currentUserID := UserIDFromContext(c)
			if comment.AuthorID != currentUserID {
				return echo.NewHTTPError(http.StatusForbidden, "you can only delete your own comments")
			}

			c.Set("comment", comment)
			return next(c)
		}
	}
}
