package query

import (
	"database/sql"
	"errors"
	"fmt"

	"strings"
	"unicode"

	"github.com/antonybholmes/go-sys"
	"github.com/antonybholmes/go-sys/log"
)

type (
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

func normalizeBooleanWords(input string) string {
	var b strings.Builder

	inQuotes := false

	for i := 0; i < len(input); {
		ch := input[i]

		if ch == '"' {
			inQuotes = !inQuotes
			b.WriteByte(ch)
			i++
			continue
		}

		// look for runs of letters that could be boolean words
		// and either replace or leaveas as is e.g x AND y -> x+y
		// but cANDy stays as is
		if !inQuotes && sys.IsLetter(ch) {
			start := i
			for i < len(input) && sys.IsLetter(input[i]) {
				i++
			}

			word := input[start:i]
			switch word {
			case "AND":
				b.WriteByte('+')
			case "OR":
				b.WriteByte(',')
			default:
				b.WriteString(word)
			}
			continue
		}

		b.WriteByte(ch)
		i++
	}

	return b.String()
}

// Deals with implicit ANDs represented by spaces between words
// and also normalizes boolean words like AND/OR to +/,
func normalizeQuery(input string) string {
	var b strings.Builder
	inQuotes := false
	last := rune(0)
	isRunOfSpaces := false

	// replace of AND/OR words first
	input = normalizeBooleanWords(input)

	for _, ch := range input {
		switch ch {
		case '"':
			inQuotes = !inQuotes
			b.WriteRune(ch)
			last = ch
			isRunOfSpaces = false

		case ' ':
			if inQuotes {
				b.WriteRune(ch)
				isRunOfSpaces = false
				continue
			}

			// we last ended on a word character
			// so this could be an implicit AND if
			// it's a run of spaces between word characters
			if isSearchTermChar(last) {
				isRunOfSpaces = true
			}

		case '+', ',', '(':
			b.WriteRune(ch)

			if !inQuotes {
				last = ch
				isRunOfSpaces = false
			}
		case '*':
			if !inQuotes {
				// convert * to % for SQL wildcard
				ch = '%'
			}

			b.WriteRune(ch)
			last = ch
			isRunOfSpaces = false
		case '?':
			if !inQuotes {
				// convert _ to _ for SQL single char wildcard
				ch = '_'
			}

			b.WriteRune(ch)
			last = ch
			isRunOfSpaces = false

		default:
			// the last was a word char and we had a
			// run of spaces, so insert a + as implicit AND
			if isRunOfSpaces && isSearchTermChar(last) && isSearchTermChar(ch) {
				b.WriteRune('+')
			}

			b.WriteRune(ch)
			last = ch
			isRunOfSpaces = false
		}
	}

	return b.String()
}

// represents a search token that is not part of a boolean expression
func isWordChar(c rune) bool {
	return unicode.IsLetter(c) ||
		unicode.IsDigit(c) ||
		c == '-' ||
		c == '_' ||
		c == '=' ||
		c == '^' ||
		c == '$' ||
		c == '.' ||
		c == ':' ||
		c == '%' ||
		c == '*' ||
		c == '?'
}

func isSearchTermChar(c rune) bool {
	return isWordChar(c) || c == '(' || c == ')'
}

// func peek(s string, i int) rune {
// 	if i >= len(s) {
// 		return 0
// 	}
// 	return rune(s[i])
// }

// Creates a new SearchTermNode from the raw string.
// This is the atomic unit of a search expression tree.
func newSearchTermNode(raw string, isExact bool, hasWildcards bool) (*SearchTermNode, error) {

	if len(raw) == 0 {
		return nil, errors.New("empty search term")
	}

	ret := SearchTermNode{Term: raw} //, MatchType: MatchTypeContains}

	// if strings.HasPrefix(raw, "-") {
	// 	ret.Not = true
	// 	raw = raw[1:]

	// 	if len(raw) == 0 {
	// 		return nil, errors.New("empty search term")
	// 	}
	// }

	// Here wildcards refers only to %. If the user has specified
	// single character wildcards with ?, we will still add % around
	// the term unless it's an exact match
	if !isExact && !hasWildcards {
		// if not exact match and user has not specified wildcards,
		// we default to contains by adding % around the term since
		// this the most intuitive behavior to look for anything
		// that contains the term we asked for
		ret.Term = "%" + ret.Term + "%"
	}

	// switch {
	// case strings.HasPrefix(raw, "="):
	// 	if len(raw) == 1 {
	// 		return nil, errors.New("empty search term")
	// 	}

	// 	ret.Value = raw[1:]
	// 	ret.MatchType = MatchTypeExact

	// case strings.HasPrefix(raw, "^") && strings.HasSuffix(raw, "$"):
	// 	if len(raw) == 2 {
	// 		return nil, errors.New("empty search term")
	// 	}

	// 	ret.Value = raw[1 : len(raw)-1]
	// 	ret.MatchType = MatchTypeExact

	// case strings.HasPrefix(raw, "^"):
	// 	if len(raw) == 1 {
	// 		return nil, errors.New("empty search term")
	// 	}

	// 	ret.Value = raw[1:]
	// 	ret.MatchType = MatchTypeStartsWith

	// case strings.HasSuffix(raw, "$"):
	// 	if len(raw) == 1 {
	// 		return nil, errors.New("empty search term")
	// 	}

	// 	ret.Value = raw[:len(raw)-1]
	// 	ret.MatchType = MatchTypeEndsWith

	// default:
	// 	// default to starts with search since
	// 	// it's the most useful for gene symbols etc.
	// 	ret.Value = raw
	// 	//ret.MatchType = MatchTypeStartsWith
	// }

	return &ret, nil
}

// func (v VarNode) Eval(vars map[string]bool) bool {
// 	return vars[v.Value]
// }

func (v SearchTermNode) BuildSql(clause SqlClauseFunc, addParens bool, args *[]string) string {
	// switch v.MatchType {
	// case MatchTypeExact:
	// 	*args = append(*args, v.Value)
	// case MatchTypeStartsWith:
	// 	*args = append(*args, v.Value+"%")
	// case MatchTypeEndsWith:
	// 	*args = append(*args, "%"+v.Value)
	// default:
	// 	*args = append(*args, "%"+v.Value+"%")
	// }

	*args = append(*args, v.Term)

	//log.Debug().Msgf("args %v", args)

	// as we parse the tree in order, the placeholder index is just
	// the current length of the args slice as we add to it
	placeholderIndex := len(*args) //fmt.Sprintf("?%d", len(*args))
	//tagClauses = append(tagClauses, fmt.Sprintf("(gex.gene_symbol LIKE %s OR gex.ensembl_id LIKE %s)", placeholder, placeholder))
	return clause(placeholderIndex, v.Term, addParens)
}

// func (a AndNode) Eval(vars map[string]bool) bool {
// 	return a.Left.Eval(vars) && a.Right.Eval(vars)
// }

// highest precedence is to negate a term
func (n NotNode) BuildSql(clause SqlClauseFunc, addParens bool, args *[]string) string {
	return "NOT (" + n.Child.BuildSql(clause, false, args) + ")"
}

func (a AndNode) BuildSql(clause SqlClauseFunc, addParens bool, args *[]string) string {
	return AddParens(a.Left.BuildSql(clause, true, args)+" AND "+a.Right.BuildSql(clause, true, args), addParens)
}

// func (o OrNode) Eval(vars map[string]bool) bool {
// 	return o.Left.Eval(vars) || o.Right.Eval(vars)
// }

func (o OrNode) BuildSql(clause SqlClauseFunc, addParens bool, args *[]string) string {
	return AddParens(o.Left.BuildSql(clause, true, args)+" OR "+o.Right.BuildSql(clause, true, args), addParens)
}

func NewParser(input string) *Parser {
	return &Parser{input: input}
}

func (p *Parser) peek() rune {
	if p.pos >= len(p.input) {
		return 0
	}
	return rune(p.input[p.pos])
}

func (p *Parser) next() rune {
	if p.pos >= len(p.input) {
		return 0
	}
	ch := rune(p.input[p.pos])
	p.pos++
	return ch
}

// Skip whitespace characters to trim
// input so tokens cannot have leading/trailing spaces
func (p *Parser) skipWhitespace() {
	for unicode.IsSpace(p.peek()) {
		p.next()
	}
}

// Entry point: parse expression with OR as lowest precedence
// e.g. A + B, C would be (A AND B) OR C
func (p *Parser) ParseExpr() (Node, error) {
	return p.parseOrSubClause()
}

// Handle ORs
func (p *Parser) parseOrSubClause() (Node, error) {
	// and takes precedence over or
	left, err := p.parseAndSubClause()

	if err != nil {
		return nil, err
	}

	// as we scan, if we see commas, we have ORs otherwise
	// we assume the spaces between terms represent ANDs
	for {
		p.skipWhitespace()

		if p.peek() == ',' {
			p.next()
			right, err := p.parseAndSubClause()

			if err != nil {
				return nil, err
			}

			left = OrNode{Left: left, Right: right}
		} else {
			break
		}
	}
	return left, nil
}

func (p *Parser) parseNotSubClause() (Node, error) {
	p.skipWhitespace()

	// check for NOT operator
	if p.peek() == '-' || p.peek() == '!' {
		p.next()
		child, err := p.parseSearchTerm()

		if err != nil {
			return nil, err
		}

		return NotNode{Child: child}, nil
	}

	// otherwise parse as normal search term
	return p.parseSearchTerm()
}

// Handle ANDs
func (p *Parser) parseAndSubClause() (Node, error) {
	left, err := p.parseNotSubClause() //  p.parseSearchTerm()

	if err != nil {
		return nil, err
	}

	for {
		p.skipWhitespace()
		if p.peek() == '+' {
			p.next()
			right, err := p.parseSearchTerm()

			if err != nil {
				return nil, err
			}

			left = AndNode{Left: left, Right: right}
		} else {
			break
		}
	}

	return left, nil
}

// Parse variables and parentheses
// func (p *Parser) parseFactor() Node {
// 	p.skipWhitespace()
// 	ch := p.peek()

// 	if ch == '(' {
// 		p.next()
// 		expr := p.parseExpr()
// 		if p.peek() != ')' {
// 			panic("missing closing parenthesis")
// 		}
// 		p.next()
// 		return expr
// 	}

// 	// Read variable
// 	start := p.pos
// 	for unicode.IsLetter(p.peek()) {
// 		p.next()
// 	}
// 	if start == p.pos {
// 		panic("expected variable")
// 	}
// 	name := strings.TrimSpace(p.input[start:p.pos])
// 	return VarNode{Name: name}
// }

func (p *Parser) parseSearchTerm() (Node, error) {
	p.skipWhitespace()
	ch := p.peek()

	// if we see a '(', we have a sub-expression
	// so we parse that recursively
	if ch == '(' {
		p.next()
		expr, err := p.parseOrSubClause()

		if err != nil {
			return nil, err
		}

		if p.peek() != ')' {
			return nil, errors.New("missing closing parenthesis")
		}

		p.next()
		return expr, nil
	}

	// if we see a quote, we have a quoted variable so
	// parse until the closing quote as is an keep all chars
	if ch == '"' {
		p.next()
		start := p.pos
		for {
			// we want to keep all chars until we see the closing quote
			// or the end of the input
			if p.peek() == '"' || p.peek() == 0 {
				break
			}
			p.next()
		}

		if p.peek() != '"' {
			return nil, errors.New("unterminated quoted variable")
		}

		word := p.input[start:p.pos]

		// consume closing quote
		p.next()

		ret, err := newSearchTermNode(word, true, false)

		if err != nil {
			return nil, err
		}

		return ret, nil
	}

	// Unquoted variable
	start := p.pos

	// keep advancing until we reach the end of the
	// search term i.e. a space or other expression char
	for isWordChar(p.peek()) {
		p.next()
	}

	// if we did not advance, it's an error
	if start == p.pos {
		return nil, errors.New("expected variable")
	}

	// extract the token without leading/trailing spaces
	value := strings.TrimSpace(p.input[start:p.pos])

	isExact := strings.HasPrefix(value, "=") || (strings.HasPrefix(value, "^") && strings.HasSuffix(value, "$"))
	hasWildcards := strings.ContainsAny(value, "%")
	hasSingleWildcards := strings.ContainsAny(value, "?")

	if isExact && (hasWildcards || hasSingleWildcards) {
		return nil, errors.New("cannot have wildcards in exact match")
	}

	if isExact {
		// strip exact match markers
		value = strings.TrimPrefix(value, "=")
		value = strings.TrimPrefix(value, "^")
		value = strings.TrimSuffix(value, "$")
	}

	// make it into a SearchNode which also determines the match type
	ret, err := newSearchTermNode(value, isExact, hasWildcards)

	if err != nil {
		return nil, err
	}

	return ret, nil
}

type SqlBoolQueryResp struct {
	Sql  string
	Args []string
}

func SqlBoolTree(query string) (Node, error) {

	// first normalize query to replace spaces with + to be treated as ands
	query = normalizeQuery(query)

	log.Debug().Msgf("normalized query: %s", query)

	parser := NewParser(query)

	// create the expression tree
	tree, err := parser.ParseExpr()

	if err != nil {
		return nil, err
	}

	return tree, nil
}

func SqlBoolQueryFromTree(tree Node, clause SqlClauseFunc) (*SqlBoolQueryResp, error) {

	// required so that we can use it with sqlite params
	args := make([]string, 0, 20)

	// build the sql from the tree
	sql := tree.BuildSql(clause, false, &args)

	return &SqlBoolQueryResp{Sql: sql, Args: args}, nil

}

func SqlBoolQuery(query string, clause SqlClauseFunc) (*SqlBoolQueryResp, error) {

	// first normalize query to replace spaces with + to be treated as ands
	query = normalizeQuery(query)

	parser := NewParser(query)

	// create the expression tree
	tree, err := parser.ParseExpr()

	if err != nil {
		return nil, err
	}

	// As we parse, we build up the args slice of search terms
	// in order that they appear in the sql
	args := make([]string, 0, 20)

	// build the sql from the tree
	sql := tree.BuildSql(clause, false, &args)

	return &SqlBoolQueryResp{Sql: sql, Args: args}, nil

}

// Adds parentheses around the sql if addParens is true
func AddParens(sql string, addParens bool) string {
	if addParens {
		return "(" + sql + ")"
	}
	return sql
}

func IndexedPlaceholder(index int) string {
	return fmt.Sprintf("p%d", index)
}

// Creates a standard sql named parameter like :p1, :p2, etc.
func IndexedParam(index int) string {
	return ":" + IndexedPlaceholder(index)
}

// Defines named args for sql queries with indexed parameters
func IndexedNamedArgs(args []string) []any {
	ret := make([]any, 0, len(args))

	for i, arg := range args {
		ret = append(ret, sql.Named(IndexedPlaceholder(i+1), arg))
	}
	return ret
}

// func main() {
// 	expr := "A+B,C+(D+E)"
// 	vars := map[string]bool{
// 		"A": true,
// 		"B": true,
// 		"C": false,
// 		"D": true,
// 		"E": false,
// 	}

// 	parser := NewParser(expr)
// 	tree := parser.parseExpr()

// 	fmt.Println("Evaluated:", tree.Eval(vars)) // true
// }
