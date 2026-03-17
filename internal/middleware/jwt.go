package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

const (
	userIDKey = "currentUserID"
	tokenKey  = "currentToken"
)

// JWTAuth returns an Echo middleware that validates Bearer JWT tokens.
// When called with optional=true the middleware allows unauthenticated requests
// through but still parses the token if present.
func JWTAuth(secret string, optional bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenStr := extractBearer(c)
			if tokenStr == "" {
				if optional {
					return next(c)
				}
				return echo.NewHTTPError(http.StatusUnauthorized, "missing or malformed JWT")
			}

			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, echo.NewHTTPError(http.StatusUnauthorized, "unexpected signing method")
				}
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				if optional {
					return next(c)
				}
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired JWT")
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid JWT claims")
			}

			sub, ok := claims["sub"].(float64)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid JWT subject")
			}

			c.Set(userIDKey, uint(sub))
			c.Set(tokenKey, tokenStr)

			return next(c)
		}
	}
}

// UserIDFromContext retrieves the authenticated user's ID (required routes).
func UserIDFromContext(c echo.Context) uint {
	v, _ := c.Get(userIDKey).(uint)
	return v
}

// OptionalUserIDFromContext returns 0 if no authenticated user is present.
func OptionalUserIDFromContext(c echo.Context) uint {
	v, _ := c.Get(userIDKey).(uint)
	return v
}

// TokenFromContext retrieves the raw JWT string stored by the middleware.
func TokenFromContext(c echo.Context) string {
	v, _ := c.Get(tokenKey).(string)
	return v
}

func extractBearer(c echo.Context) string {
	auth := c.Request().Header.Get("Authorization")
	if strings.HasPrefix(auth, "Token ") {
		return strings.TrimPrefix(auth, "Token ")
	}
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	return ""
}
