package utils

import (
	"regexp"
	"strings"
)

var (
	slugNonAlnum = regexp.MustCompile(`[^a-z0-9]+`)
	slugTrim     = regexp.MustCompile(`^-+|-+$`)
)

func Slugify(s string) string {

	s = strings.ToLower(s)
	s = slugNonAlnum.ReplaceAllString(s, "-")
	s = slugTrim.ReplaceAllString(s, "")

	return s
}
