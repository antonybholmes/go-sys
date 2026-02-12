package query

import (
	"strings"

	"github.com/antonybholmes/go-sys"
)

type (
	SqlClauseFunc func(placeholderIndex int, value string, addParens bool) string

	// Node represents an expression tree node
	Node interface {
		// Builds the sql representation of this node
		// using the given clause function to create
		// the sql for each search term
		// Param clause function to create sql clause for
		// each search term here it becomes user and database
		// specific so user supplies function to supply the
		// actual sql clause for each term
		// Param args is a pointer to a slice of strings
		// that we append the actual search term values to
		// as we build the sql so that the caller can use
		// them as query parameters.
		// Param addParens indicates whether to add parentheses
		// around the generated sql for this node, useful to
		// reduce excessively nested expressions, e.g. NOT operand
		// does not need parens around its child as it is self contained.
		// We delegate whether to add parens to the sub expression since
		// simple sql statements may not require them.
		BuildSql(clause SqlClauseFunc, addParens bool, args *[]string) string
	}

	// SearchTermNode for variables like A, B, etc.
	SearchTermNode struct {
		Term string
	}

	NotNode struct {
		Child Node
	}

	// AndNode for AND operations
	AndNode struct {
		Left  Node
		Right Node
	}

	// OrNode for OR operations
	OrNode struct {
		Left  Node
		Right Node
	}

	// Parser struct
	Parser struct {
		input string
		pos   int
	}
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

