package handler_test

import (
	"github.com/labstack/echo/v4"

	appvalidator "github.com/khanhnp-2797/echo-realworld-api/pkg/validator"
)

// newEcho returns a pre-configured Echo instance matching production setup.
func newEcho() *echo.Echo {
	e := echo.New()
	e.Validator = appvalidator.New()
	return e
}

// setAuth injects the authenticated user ID (and a stub token) into an Echo context,
// simulating a successfully verified JWT without running the actual middleware.
func setAuth(c echo.Context, userID uint) {
	c.Set("currentUserID", userID)
	c.Set("currentToken", "test-token")
}
