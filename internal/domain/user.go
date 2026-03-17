package domain

import "gorm.io/gorm"

// User represents an application user.
type User struct {
	gorm.Model
	Username string  `gorm:"uniqueIndex;not null;size:255"`
	Email    string  `gorm:"uniqueIndex;not null;size:255"`
	Password string  `gorm:"not null"`
	Bio      *string `gorm:"size:1000"`
	Image    *string `gorm:"size:500"`
}
