package apperrors

import "errors"

// Sentinel errors used across the application.
var (
	ErrNotFound      = errors.New("resource not found")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrForbidden     = errors.New("forbidden")
	ErrEmailTaken    = errors.New("email already taken")
	ErrUsernameTaken = errors.New("username already taken")
	ErrInvalidCreds  = errors.New("invalid email or password")
	ErrInternal      = errors.New("internal server error")
)
