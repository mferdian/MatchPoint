package helpers

import (
	"regexp"
	"strings"
)

func SnakeCase(str string) string {
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := re.ReplaceAllString(str, "${1}_${2}")
	re2 := regexp.MustCompile("([A-Z]+)([A-Z][a-z])")
	snake = re2.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
