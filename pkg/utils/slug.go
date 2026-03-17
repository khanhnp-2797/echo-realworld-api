package utils

import (
	"regexp"
	"strings"
)

var nonAlphanumeric = regexp.MustCompile(`[^a-z0-9]+`)

// Slugify converts a title into a URL-friendly slug.
// e.g. "How to Train Your Dragon" -> "how-to-train-your-dragon"
func Slugify(title string) string {
	slug := strings.ToLower(title)
	slug = nonAlphanumeric.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}
