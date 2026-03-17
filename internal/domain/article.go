package domain

import "gorm.io/gorm"

// Article is the core content entity.
type Article struct {
	gorm.Model
	Slug        string `gorm:"uniqueIndex;not null;size:255"`
	Title       string `gorm:"not null;size:500"`
	Description string `gorm:"size:1000"`
	Body        string `gorm:"not null;type:text"`

	AuthorID uint `gorm:"not null;index"`
	Author   User

	// Many2Many: article_tags join table
	Tags     []*Tag    `gorm:"many2many:article_tags"`
	Comments []Comment // HasMany
}
