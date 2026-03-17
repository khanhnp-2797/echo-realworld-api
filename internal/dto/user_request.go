package dto

// RegisterRequest is the request body for POST /api/users.
type RegisterRequest struct {
	User struct {
		Username string `json:"username" validate:"required,min=2,max=100"`
		Email    string `json:"email"    validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	} `json:"user"`
}

// LoginRequest is the request body for POST /api/users/login.
type LoginRequest struct {
	User struct {
		Email    string `json:"email"    validate:"required,email"`
		Password string `json:"password" validate:"required"`
	} `json:"user"`
}
