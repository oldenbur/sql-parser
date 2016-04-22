package sql

import (
	"fmt"
	"strconv"
)

type Expr interface {
	String() string
}

func (p *Parser) parseExpr() (Expr, error) {

	tok, arg := p.scanIgnoreWhitespace()
	switch(tok) {
	case STRING:
		return &StringExpr{Val: arg}, nil
	case NUMBER:
		numVal, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			return nil, fmt.Errorf("ParseExpr() error in ParseFloat('%s'): %v", arg, err)
		}
		return &NumExpr{Val: numVal}, nil
	case IDENT:
		p.unscan()
		return p.parseFuncCall()
	default:
		return nil, fmt.Errorf(`parseExpr() expected expression (string, number or function call), got %v`, tok)
	}
}

type FuncCallExpr struct {
	Name string
	Args []Expr
}

func (f *FuncCallExpr) String() string {

	var argList string = ""
	for i, a := range f.Args {
		sep := ""
		if i > 0 {
			sep = ", "
		}
		argList = fmt.Sprintf("%s%s%s", argList, sep, a)
	}

	return fmt.Sprintf("%s(%s)", f.Name, argList)
}

type StringExpr struct {
	Val string
}

func (s *StringExpr) String() string {
	return s.Val
}

type NumExpr struct {
	Val float64
}

func (n *NumExpr) String() string {
	return fmt.Sprintf("%f", n.Val)
}

func (p *Parser) parseFuncCall() (Expr, error) {

	var funcName string
	var args []Expr = make([]Expr, 0)

	tok, ident := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf(`expected IDENT, got '%s'`, ident)
	}
	funcName = ident

	tok, arg := p.scanIgnoreWhitespace()
	if tok != PAREN_L {
		return nil, fmt.Errorf(`expected '(', got '%s'`, arg)
	}

	tok, _ = p.scanIgnoreWhitespace()
	i := 1
	for tok != EOF && tok != PAREN_R {

		p.unscan()
		e, err := p.parseExpr()
		if err != nil {
			return nil, fmt.Errorf(`error parsing %s argument %d: %v`, funcName, i, err)
		}
		args = append(args, e)

		tok, arg = p.scanIgnoreWhitespace()
		if tok == COMMA {
			tok, _ = p.scanIgnoreWhitespace()
		} else if tok != PAREN_R {
			return nil, fmt.Errorf(`expected COMMA or PAREN_R after %s arg %d, got %v`, funcName, i, arg)
		}

		i += 1
	}

	if tok == EOF {
		return nil, fmt.Errorf(`expected PAREN_R in function %s, got EOF`, funcName)
	}

	return &FuncCallExpr{Name: funcName, Args: args}, nil
}
