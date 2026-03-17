package service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/khanhnp-2797/echo-realworld-api/internal/config"
	"github.com/khanhnp-2797/echo-realworld-api/internal/domain"
	"github.com/khanhnp-2797/echo-realworld-api/internal/repository"
	"github.com/khanhnp-2797/echo-realworld-api/pkg/apperrors"
	"golang.org/x/crypto/bcrypt"
)

// UserService defines the business operations on User.
type UserService interface {
	// Task 4: Auth
	Register(ctx context.Context, username, email, password string) (*domain.User, string, error)
	Login(ctx context.Context, email, password string) (*domain.User, string, error)
	GetByID(ctx context.Context, id uint) (*domain.User, error)

	// Task 3: Profiles
	GetProfile(ctx context.Context, username string) (*domain.User, error)
}

type userService struct {
	repo   repository.UserRepository
	jwtCfg config.JWTConfig
}

func NewUserService(repo repository.UserRepository, jwtCfg config.JWTConfig) UserService {
	return &userService{repo: repo, jwtCfg: jwtCfg}
}

// Register hashes the password and persists a new user, returning a JWT.
func (s *userService) Register(ctx context.Context, username, email, password string) (*domain.User, string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	user := &domain.User{
		Username: username,
		Email:    email,
		Password: string(hashed),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, "", err
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// Login validates credentials and returns the user with a fresh JWT.
func (s *userService) Login(ctx context.Context, email, password string) (*domain.User, string, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, "", apperrors.ErrInvalidCreds
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", apperrors.ErrInvalidCreds
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// GetByID fetches the current authenticated user by ID.
func (s *userService) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	return user, nil
}

// GetProfile fetches a user's public profile by username.
func (s *userService) GetProfile(ctx context.Context, username string) (*domain.User, error) {
	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	return user, nil
}

// generateToken creates a signed HS256 JWT for the given user ID.
func (s *userService) generateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * time.Duration(s.jwtCfg.ExpireHours)).Unix(),
		"iat": time.Now().Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.jwtCfg.Secret))
}

// UserUpdateInput carries optional fields for profile updates.
