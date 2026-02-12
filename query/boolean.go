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

// Replace boolean words AND/OR with +/, but only when they are separate words
// and not with quotes
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
			// find run of letters
			for i < len(input) && sys.IsLetter(input[i]) {
				i++
			}

			// since a run must start with a letter, then
			// came before (and after) must be non-letter e.g.
			// space, therefore if we see AND and OR they
			// must be surrounded by non-letter chars and can
			// be treated as boolean operators
			word := input[start:i]
			switch word {
			case "AND":
				b.WriteByte('+')
			case "OR":
				b.WriteByte(',')
			default:
				// not a boolean word, write as is
				b.WriteString(word)
			}
			continue
		}

		// not a letter, write as is
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

			// if the last non-space character was a
			// search term char e.g. a letter, we could
			// be in the middle of an implicit AND i.e.
			// a space between two search terms
			// so we flag that we could be starting a
			// run of spaces that represent an implicit AND
			if isSearchTermChar(last) {
				isRunOfSpaces = true
			}

		case '+', ',': //, '(':
			b.WriteRune(ch)

			// if we are in a run of spaces
			// cancel the run since are either
			// AND/OR or
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
				// convert ? to _ for SQL single char wildcard
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
