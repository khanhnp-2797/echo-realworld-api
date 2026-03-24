package repository

import (
	"context"
	"strings"

	"github.com/khanhnp-2797/echo-realworld-api/internal/domain"
	"github.com/khanhnp-2797/echo-realworld-api/pkg/apperrors"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository returns a GORM-backed UserRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		e := err.Error()
		if strings.Contains(e, "idx_users_email") {
			return apperrors.ErrEmailTaken
		}
		if strings.Contains(e, "idx_users_username") {
			return apperrors.ErrUsernameTaken
		}
		return err
	}
	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Follow creates a follower_id → followed_id row in user_follows.
func (r *userRepository) Follow(ctx context.Context, followerID, followedID uint) error {
	follower := domain.User{}
	follower.ID = followerID
	followed := domain.User{}
	followed.ID = followedID
	return r.db.WithContext(ctx).Model(&follower).Association("Following").Append(&followed)
}

// Unfollow removes the follow relationship.
func (r *userRepository) Unfollow(ctx context.Context, followerID, followedID uint) error {
	follower := domain.User{}
	follower.ID = followerID
	followed := domain.User{}
	followed.ID = followedID
	return r.db.WithContext(ctx).Model(&follower).Association("Following").Delete(&followed)
}

// IsFollowing returns true if followerID follows followedID.
func (r *userRepository) IsFollowing(ctx context.Context, followerID, followedID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("user_follows").
		Where("follower_id = ? AND followed_id = ?", followerID, followedID).
		Count(&count).Error
	return count > 0, err
}
