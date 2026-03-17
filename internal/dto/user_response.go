package dto

import "github.com/khanhnp-2797/echo-realworld-api/internal/domain"

// UserResponse wraps UserBody for single-user endpoints.
type UserResponse struct {
	User UserBody `json:"user"`
}

// UserBody is the public DTO for an authenticated user (password excluded).
type UserBody struct {
	Email    string  `json:"email"`
	Token    string  `json:"token"`
	Username string  `json:"username"`
	Bio      *string `json:"bio"`
	Image    *string `json:"image"`
}

// ProfileResponse wraps ProfileBody for profile endpoints.
type ProfileResponse struct {
	Profile ProfileBody `json:"profile"`
}

// ProfileBody is the public DTO for a user profile.
// Also embedded inside ArticleBody and CommentBody.
type ProfileBody struct {
	Username string  `json:"username"`
	Bio      *string `json:"bio"`
	Image    *string `json:"image"`
}

// ToUserBody maps a domain.User + JWT token to a UserBody DTO.
func ToUserBody(u *domain.User, token string) UserBody {
	return UserBody{
		Email:    u.Email,
		Token:    token,
		Username: u.Username,
		Bio:      u.Bio,
		Image:    u.Image,
	}
}

// ToProfileBody maps a domain.User to a ProfileBody DTO.
func ToProfileBody(u *domain.User) ProfileBody {
	return ProfileBody{
		Username: u.Username,
		Bio:      u.Bio,
		Image:    u.Image,
	}
}
