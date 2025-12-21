package query

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/antonybholmes/go-sys"
)

type (
	MatchType int

	// Given a placeholder and match type create the sql to
	// perform the actual match. Placeholder is the numerical
	// index from 1...n of the variable to insert into the sql
	// statement, therefore for sql systems that support ?1 type
	// variables you can use that, otherwise it can be ignored
	// and the generic '?' used.
	//SqlClauseFunc func(placeholderIndex int, matchType MatchType, not bool) string
	SqlClauseFunc func(placeholderIndex int, value string, addParens bool) string

	Term struct {
		Value string
		Exact bool
	}
)

const (
	MatchTypeExact MatchType = iota
	MatchTypeStartsWith
	MatchTypeEndsWith
	MatchTypeContains
)

var (
	// matches any invalid characters for stripping out of queries,
	// not we allow parentheses for simple boolean query building if
	// desired
	InvalidCharsRegex = regexp.MustCompile(`[^a-zA-Z0-9\ \_\,\+\=\"\^\$\(\)\.\-\:\?\*\%\!]+`)

	AndTermRegex = regexp.MustCompile(`(=)?"([^"]+)"|(=)?([^"+,\s]+)`)
)

func SanitizeQuery(input string) string {
	return strings.TrimSpace(sys.NormalizeSpaces(InvalidCharsRegex.ReplaceAllString(input, "")))
}

// Parses a query into blocks of and tags using
// a simple query format of plus and space for AND and comma for OR
// Thus foo+bar,thing -> [['foo', 'bar'], ['thing']]
func ParseQuery(query string) (orTags []string, terms [][]*Term) {
	orGroups := strings.Split(query, ",") // comma separates OR groups

	terms = make([][]*Term, 0, len(orGroups))

	for _, group := range orGroups {
		group = strings.TrimSpace(group)

		if group == "" {
			continue
		}

		andTerms := make([]*Term, 0, len(orGroups))

		parts := AndTermRegex.FindAllStringSubmatch(group, -1)
		for _, m := range parts {
			exact := m[1] == "=" || m[3] == "="

			value := strings.TrimSpace(m[2])

			if value == "" {
				value = strings.TrimSpace(m[4])
			}

			if value != "" {
				andTerms = append(andTerms, &Term{Value: value, Exact: exact})
			}
		}

		if len(andTerms) > 0 {
			terms = append(terms, andTerms)
		}

		// trim each piece and replace spaces with + since we treat spaces as being AND
		//parts := strings.Split(SPACES_REGEX.ReplaceAllString(group, "+"), "+") // plus separates AND parts
		//terms = append(terms, parts)
	}
	return orGroups, terms
}

// Creates a boolean sql query from a text query using + for AND and comma for OR.
// User must supply a clause function that given a placeholder string, returns
// the core part of the query that matches to the placeholder item. This is done
// this way to provide flexibility when defining the query, e.g. we can check multiple
// table fields for the same placeholder if necessary
func BoolQuery(query string, clause func(placeholder string, exact bool) string) (string, []interface{}) {
	_, andTags := ParseQuery(query)

	andClauses := make([]string, 0, len(andTags))

	// required so that we can use it with sqlite params
	args := make([]interface{}, 0, len(andTags))

	for _, group := range andTags {
		tagClauses := make([]string, 0, len(group))
		for _, tag := range group {
			if tag.Exact {
				args = append(args, tag.Value)
			} else {
				args = append(args, "%"+tag.Value+"%")
			}

			//log.Debug().Msgf("args %v", args)
			//log.Debug().Msgf("args %v", tag)

			placeholder := fmt.Sprintf("?%d", len(args))
			//tagClauses = append(tagClauses, fmt.Sprintf("(gex.gene_symbol LIKE %s OR gex.ensembl_id LIKE %s)", placeholder, placeholder))
			tagClauses = append(tagClauses, clause(placeholder, tag.Exact))
		}
		andClauses = append(andClauses, "("+strings.Join(tagClauses, " AND ")+")")
	}

	finalSQL := strings.Join(andClauses, " OR ")

	return finalSQL, args
}
