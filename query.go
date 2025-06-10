package sys

import (
	"fmt"
	"strings"
	"unicode"
)

// Node represents an expression tree node
type Node interface {
	Eval(vars map[string]bool) bool
	BuildSql(args *[]interface{}, clause ClauseFunc) string
}

// VarNode for variables like A, B, etc.
type VarNode struct {
	Value string
	Exact bool
}

func (v VarNode) Eval(vars map[string]bool) bool {
	return vars[v.Value]
}

func (v VarNode) BuildSql(args *[]interface{}, clause ClauseFunc) string {
	if v.Exact {
		*args = append(*args, v.Value)
	} else {
		*args = append(*args, "%"+v.Value+"%")
	}

	//log.Debug().Msgf("args %v", args)
	//log.Debug().Msgf("args %v", tag)

	placeholder := fmt.Sprintf("?%d", len(*args))
	//tagClauses = append(tagClauses, fmt.Sprintf("(gex.gene_symbol LIKE %s OR gex.ensembl_id LIKE %s)", placeholder, placeholder))
	return clause(placeholder, v.Exact)
}

// AndNode for AND operations
type AndNode struct {
	Left  Node
	Right Node
}

func (a AndNode) Eval(vars map[string]bool) bool {
	return a.Left.Eval(vars) && a.Right.Eval(vars)
}

func (a AndNode) BuildSql(args *[]interface{}, clause ClauseFunc) string {
	return "(" + a.Left.BuildSql(args, clause) + " AND " + a.Right.BuildSql(args, clause) + ")"
}

// OrNode for OR operations
type OrNode struct {
	Left  Node
	Right Node
}

func (o OrNode) Eval(vars map[string]bool) bool {
	return o.Left.Eval(vars) || o.Right.Eval(vars)
}

func (o OrNode) BuildSql(args *[]interface{}, clause ClauseFunc) string {
	return "(" + o.Left.BuildSql(args, clause) + " OR " + o.Right.BuildSql(args, clause) + ")"
}

func isVariable(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsDigit(c) || c == '-' || c == '_'
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
func (p *Parser) ParseExpr() Node {
	left := p.parseTerm()
	for {
		p.skipWhitespace()
		if p.peek() == ',' {
			p.next()
			right := p.parseTerm()
			left = OrNode{Left: left, Right: right}
		} else {
			break
		}
	}
	return left
}

// Handle ANDs
func (p *Parser) parseTerm() Node {
	left := p.parseFactor()
	for {
		p.skipWhitespace()
		if p.peek() == '+' {
			p.next()
			right := p.parseFactor()
			left = AndNode{Left: left, Right: right}
		} else {
			break
		}
	}
	return left
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

func (p *Parser) parseFactor() Node {
	p.skipWhitespace()
	ch := p.peek()

	if ch == '(' {
		p.next()
		expr := p.ParseExpr()
		if p.peek() != ')' {
			panic("missing closing parenthesis")
		}
		p.next()
		return expr
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
			panic("unterminated quoted variable")
		}
		name := p.input[start:p.pos]
		p.next() // consume closing quote
		return VarNode{Value: name}
	}

	// Unquoted variable
	start := p.pos

	for isVariable(p.peek()) {
		p.next()
	}

	if start == p.pos {
		panic("expected variable")
	}

	name := strings.TrimSpace(p.input[start:p.pos])

	return VarNode{Value: name}
}

func SqlBoolQuery(query string, clause ClauseFunc) (string, []interface{}) {

	args := make([]interface{}, 0, 20)

	parser := NewParser(query)
	tree := parser.ParseExpr()

	sql := tree.BuildSql(&args, clause)

	return sql, args

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
