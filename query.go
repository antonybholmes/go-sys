package sys

import (
	"fmt"
	"strings"
	"unicode"
)

//var SPACES_REGEX = regexp.MustCompile(`\s+`)

func normalizeImplicitAnd(input string) string {

	var b strings.Builder
	inQuotes := false
	last := rune(0)

	for i, ch := range input {
		switch ch {
		case '"':
			// strings between quotes are used as is
			inQuotes = !inQuotes
			b.WriteRune(ch)
		case ' ':
			// we are in the middle of two words so add a + between them
			// to indicate and
			if !inQuotes && isWordChar(last) && isWordChar(peek(input, i+1)) {
				b.WriteRune('+') // implicit AND
			} else {
				b.WriteRune(ch) // preserve space inside quotes
			}
		default:
			b.WriteRune(ch)
		}

		if !unicode.IsSpace(ch) {
			last = ch
		}
	}
	return b.String()
}

func isSearchTermChar(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsDigit(c) || c == '-' || c == '_' || c == '=' || c == '^' || c == '$'
}

func isWordChar(c rune) bool {
	return isSearchTermChar(c) || c == '(' || c == ')'
}

func peek(s string, i int) rune {
	if i >= len(s) {
		return 0
	}
	return rune(s[i])
}

// Node represents an expression tree node
type Node interface {
	//Eval(vars map[string]bool) bool
	BuildSql(clause SqlClauseFunc, args *[]interface{}) string
}

// VarNode for variables like A, B, etc.
type VarNode struct {
	Value     string
	MatchType MatchType
}

func makeVarNode(raw string) VarNode {

	switch {
	case strings.HasPrefix(raw, "=") && len(raw) > 1:
		return VarNode{Value: raw[1:], MatchType: Exact}
	case strings.HasPrefix(raw, "^") && strings.HasSuffix(raw, "$") && len(raw) > 2:
		return VarNode{Value: raw[1 : len(raw)-1], MatchType: Exact}
	case strings.HasPrefix(raw, "^") && len(raw) > 1:
		return VarNode{Value: raw[1:], MatchType: StartsWith}
	case strings.HasSuffix(raw, "$") && len(raw) > 1:
		return VarNode{Value: raw[:len(raw)-1], MatchType: EndsWith}
	default:
		return VarNode{Value: raw, MatchType: Contains}
	}
}

// func (v VarNode) Eval(vars map[string]bool) bool {
// 	return vars[v.Value]
// }

func (v VarNode) BuildSql(clause SqlClauseFunc, args *[]interface{}) string {
	switch v.MatchType {
	case Exact:
		*args = append(*args, v.Value)
	case StartsWith:
		*args = append(*args, v.Value+"%")
	case EndsWith:
		*args = append(*args, "%"+v.Value)
	default:
		*args = append(*args, "%"+v.Value+"%")
	}

	//log.Debug().Msgf("args %v", args)

	placeholder := uint(len(*args)) //fmt.Sprintf("?%d", len(*args))
	//tagClauses = append(tagClauses, fmt.Sprintf("(gex.gene_symbol LIKE %s OR gex.ensembl_id LIKE %s)", placeholder, placeholder))
	return clause(placeholder, v.MatchType)
}

// AndNode for AND operations
type AndNode struct {
	Left  Node
	Right Node
}

// func (a AndNode) Eval(vars map[string]bool) bool {
// 	return a.Left.Eval(vars) && a.Right.Eval(vars)
// }

func (a AndNode) BuildSql(clause SqlClauseFunc, args *[]interface{}) string {
	return "(" + a.Left.BuildSql(clause, args) + " AND " + a.Right.BuildSql(clause, args) + ")"
}

// OrNode for OR operations
type OrNode struct {
	Left  Node
	Right Node
}

// func (o OrNode) Eval(vars map[string]bool) bool {
// 	return o.Left.Eval(vars) || o.Right.Eval(vars)
// }

func (o OrNode) BuildSql(clause SqlClauseFunc, args *[]interface{}) string {
	return "(" + o.Left.BuildSql(clause, args) + " OR " + o.Right.BuildSql(clause, args) + ")"
}

// Parser struct
type Parser struct {
	input string
	pos   int
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

func (p *Parser) skipWhitespace() {
	for unicode.IsSpace(p.peek()) {
		p.next()
	}
}

// Entry point: parse expression with OR as lowest precedence
func (p *Parser) ParseExpr() (Node, error) {
	left, err := p.parseTerm()

	if err != nil {
		return nil, err
	}

	for {
		p.skipWhitespace()
		if p.peek() == ',' {
			p.next()
			right, err := p.parseTerm()

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

// Handle ANDs
func (p *Parser) parseTerm() (Node, error) {
	left, err := p.parseFactor()

	if err != nil {
		return nil, err
	}

	for {
		p.skipWhitespace()
		if p.peek() == '+' {
			p.next()
			right, err := p.parseFactor()

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

func (p *Parser) parseFactor() (Node, error) {
	p.skipWhitespace()
	ch := p.peek()

	if ch == '(' {
		p.next()
		expr, err := p.ParseExpr()

		if err != nil {
			return nil, err
		}

		if p.peek() != ')' {
			return nil, fmt.Errorf("missing closing parenthesis")
		}

		p.next()
		return expr, nil
	}

	// Handle quoted variable names
	if ch == '"' {
		p.next()
		start := p.pos
		for {
			if p.peek() == '"' || p.peek() == 0 {
				break
			}
			p.next()
		}

		if p.peek() != '"' {
			return nil, fmt.Errorf("unterminated quoted variable")
		}

		value := p.input[start:p.pos]
		p.next() // consume closing quote

		ret := makeVarNode(value)
		return &ret, nil
	}

	// Unquoted variable
	start := p.pos

	// keep advancing until we reach the end of the
	// search term
	for isSearchTermChar(p.peek()) {
		p.next()
	}

	if start == p.pos {
		return nil, fmt.Errorf("expected variable")
	}

	value := strings.TrimSpace(p.input[start:p.pos])
	ret := makeVarNode(value)
	return &ret, nil
}

type SqlBoolQueryResp struct {
	Sql  string
	Args []interface{}
}

func SqlBoolQuery(query string, clause SqlClauseFunc) (*SqlBoolQueryResp, error) {

	// first normalize query to replace spaces with + to be treated as ands
	query = normalizeImplicitAnd(query)

	// required so that we can use it with sqlite params
	args := make([]interface{}, 0, 20)

	parser := NewParser(query)
	tree, err := parser.ParseExpr()

	if err != nil {
		return nil, err
	}

	sql := tree.BuildSql(clause, &args)

	return &SqlBoolQueryResp{Sql: sql, Args: args}, nil

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
