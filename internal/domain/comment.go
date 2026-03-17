package domain

import "gorm.io/gorm"

// Comment belongs to an Article and has an Author.
type Comment struct {
	gorm.Model
	Body string `gorm:"not null;type:text"`

	AuthorID uint `gorm:"not null;index"`
	Author   User

	ArticleID uint `gorm:"not null;index"`
}
