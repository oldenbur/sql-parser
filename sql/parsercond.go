package sql

import (
	"fmt"
)

type Cond interface {
	String() string
}

// CondComp represents a single comparison, e.g. f = 'bucky'
type CondComp struct {
	Ident string
	CondOp Token  // e.g. =, <=
	Val Expr
}

func (c CondComp) String() string {
	return fmt.Sprintf("%s %s %s", c.Ident, c.CondOp, c.Val)
}

// CondConj represents a single level of ANDed or ORed statements,
// e.g. f1 = "v1" AND myNum >= 12.34 AND (f2 != "v2" OR id = 12)
// There is an AND node with two Conds and a single Node, which is
// itself an OR node with two conds and no Nodes.
type CondConj struct {
	Op Token  // AND or OR
	Left Cond
	Right Cond
}

func (c CondConj) String() string {
	return fmt.Sprintf("(%s %s %s)", c.Left, c.Op, c.Right)
}

// parseCondTree assumes that the scanner is in position to parse a potentially
// compound boolean expression, e.g.
//   t1.field1 = "val1" AND (t2.field1 <= -12.34 OR t1.field2 != "val2")
// If parsing is successful, a populated Cond tree structure representing the
// parsed expression is returned, otherwise error.
func (p *Parser) parseCondTree() (Cond, error) {

	var left, right Cond
	var err error

	tok, lit := p.scanIgnoreWhitespace()
	if tok == PAREN_L {
		left, err = p.parseCondTree()
		if err != nil {
			return nil, err
		}
	} else if tok == IDENT {
		p.unscan()
		left, err = p.parseCondComp()
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf(`expeected PAREN_L or IDENT, got "%s"`, lit)
	}

	tok, lit = p.scanIgnoreWhitespace()
	if tok == AND || tok == OR {

		condConj := &CondConj{Op: tok, Left: left}
		right, err = p.parseCondTree()
		if err != nil {
			return nil, err
		}
		condConj.Right = right

		if tok, lit = p.scanIgnoreWhitespace(); tok != PAREN_R {
			p.unscan()
		}

		return condConj, nil
	} else if tok == PAREN_R {
		return left, nil
	} else if tok != EOF {
		return nil, fmt.Errorf(`expected AND or OR, got "%s"`, lit)
	}

	return left, nil
}

// parseCondComp assumes that the scanner is in the position to parse a condition
// expression, e.g. t1.field1 = "stringval". If parsing is successful, a populated
// CondComp structure is returned, otherwise an error.
func (p *Parser) parseCondComp() (*CondComp, error) {

	tok, ident := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf(`expected IDENT, got "%s"`, ident)
	}

	op, lit := p.scanIgnoreWhitespace()
	if !isOperator(op) {
		return nil, fmt.Errorf(`expected operator, got "%s"`, lit)
	}

	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	return &CondComp{Ident: ident, CondOp: op, Val: expr}, nil
}

func isOperator(tok Token) bool {
	return tok == EQ || tok == NE || tok == LT || tok == GT || tok == LE || tok == GE
}