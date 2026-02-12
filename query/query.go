package query

import (
	"strings"

	"github.com/antonybholmes/go-sys"
)

var (
	// O(1) lookup for allowed chars to strip out invalid chars from queries
	allowedChar [256]bool
)

func init() {
	// load the chars we are allowed into array once upon init

	for c := byte('a'); c <= 'z'; c++ {
		allowedChar[c] = true
	}

	for c := byte('A'); c <= 'Z'; c++ {
		allowedChar[c] = true
	}

	for c := byte('0'); c <= '9'; c++ {
		allowedChar[c] = true
	}

	for _, c := range []byte(" _,+=\"^$().-:?*%!") {
		allowedChar[c] = true
	}
}

func StripInvalid(s string) string {
	b := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {

		if allowedChar[s[i]] {
			b = append(b, s[i])
		}
	}
	return string(b)
}

func SanitizeQuery(input string) string {
	return strings.TrimSpace(sys.NormalizeSpaces(StripInvalid(input)))
}

// func normalizeImplicitAnd(input string) string {

// 	var b strings.Builder
// 	inQuotes := false
// 	last := rune(0)
// 	lastWasSpace := false
// 	i := 0
// 	n := len(input)

// 	for i < n {
// 		ch := rune(input[i])

// 		switch ch {
// 		case '"':
// 			// strings between quotes are used as is
// 			inQuotes = !inQuotes
// 			b.WriteRune(ch)
// 			i++
// 		case ' ':
// 			if inQuotes {
// 				// preserve spaces exactly inside quotes
// 				b.WriteRune(ch)
// 				i++
// 				continue
// 			}

// 			if isWordChar(last) {
// 				lastWasSpace = true
// 			}

// 			// consume run of spaces
// 			// j := i
// 			// for j < n && input[j] == ' ' {
// 			// 	j++
// 			// }

// 			// // this will be the character after the spaces
// 			// next := rune(0)
// 			// if j < n {
// 			// 	next = rune(input[j])
// 			// }

// 			// if isWordChar(last) && isWordChar(next) {
// 			// 	b.WriteRune('+') // implicit AND
// 			// 	last = '+'
// 			// }

// 			// // skip all spaces
// 			// i = j
// 		case '+', ',', '(':
// 			b.WriteRune(ch)
// 			i++

// 			if inQuotes {
// 				// preserve spaces exactly inside quotes
// 				continue
// 			}

// 			// the last non-space character we've seen
// 			last = ch
// 			lastWasSpace = false

// 			// consume run of spaces
// 			// for i < n && input[i] == ' ' {
// 			// 	i++
// 			// }
// 		default:
// 			b.WriteRune(ch)
// 			last = ch
// 			i++
// 		}

// 	}

// 	return b.String()
// }
