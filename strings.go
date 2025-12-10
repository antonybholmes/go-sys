package sys

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	SpacesRegex = regexp.MustCompile(`\s+`)
)

// Converts any string to camelCase suitable for JS keys
func ToCamelCaseKey(s string) string {
	// Remove any characters that are not letters or numbers
	clean := make([]rune, 0, len(s))
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' || r == '-' || r == '_' {
			clean = append(clean, r)
		}
	}

	// Replace delimiters with spaces
	str := strings.ReplaceAll(string(clean), "-", " ")
	str = strings.ReplaceAll(str, "_", " ")

	// Replace multiple spaces with a single space
	str = SpacesRegex.ReplaceAllString(str, " ")

	words := strings.Fields(str)
	if len(words) == 0 {
		return ""
	}

	// First word lowercase
	runes := []rune(words[0])
	for i := range runes {
		runes[i] = unicode.ToLower(runes[i])
	}
	words[0] = string(runes)

	// Capitalize subsequent words
	for i := 1; i < len(words); i++ {
		runes := []rune(words[i])
		if len(runes) > 0 {
			runes[0] = unicode.ToUpper(runes[0])
			for j := 1; j < len(runes); j++ {
				runes[j] = unicode.ToLower(runes[j])
			}
			words[i] = string(runes)
		}
	}

	return strings.Join(words, "")
}
