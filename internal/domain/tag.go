package domain

import "gorm.io/gorm"

// Tag is a label that can be applied to many articles.
type Tag struct {
	gorm.Model
	Name     string     `gorm:"uniqueIndex;not null;size:100"`
	Articles []*Article `gorm:"many2many:article_tags"`
}
