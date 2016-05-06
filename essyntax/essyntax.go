package essyntax

import (
	"fmt"

	"github.com/oldenbur/sql-parser/sql"
)

func ElasticSearchQuery(s *sql.SelectStatement) (string, error) {



	return "", nil
}

// genCondClause returns an elasticsearch query clause generated from the specified clause,
// which can either be a conjuction or a comparison
func genCondClause(where sql.Cond) (string, error) {

	switch where := where.(type) {
	case *sql.CondComp:
		return genCompClause(where)

	case *sql.CondConj:
		if where.Left != nil && where.Right != nil {
			return genConjClause(where)
		} else if where.Left != nil {
			return genCondClause(where.Left)
		} else if where.Right != nil {
			return genCondClause(where.Right)
		} else {
			return "", fmt.Errorf("unexpected emtpy logical conjunction")
		}

	default:
		return "", fmt.Errorf("unexpected singleIndexQuery condition type: %t", where)
	}

}

// genConjClause returns an elasticsearch bool should or must clause generated from the
// specified conjunction clause
func genConjClause(conj *sql.CondConj) (string, error) {

	leftClause, err := genCondClause(conj.Left)
	if err != nil {
		return "", err
	}

	rightClause, err := genCondClause(conj.Right)
	if err != nil {
		return "", err
	}

	var esConj string
	switch conj.Op {
	case sql.AND: esConj = "must"
	case sql.OR: esConj = "should"
	default: return "", fmt.Errorf("unexpected operator generating conjuction: %v", conj.Op)
	}

	return fmt.Sprintf(`{"bool": {"%s": [%s, %s]}}`, esConj, leftClause, rightClause), nil
}

// genCompClause creates an elasticsearch term or range clause for the specified comparison
func genCompClause(comp *sql.CondComp) (string, error) {

	switch val := comp.Val.(type) {
	case *sql.NumExpr:

		op := comp.CondOp
		if op == sql.LT || op == sql.LE || op == sql.GT || op == sql.GE {
			return fmt.Sprintf(`{"range": {"%s": {"%s": %v}}}`, comp.Ident, genRangeOp(op), val.Val), nil
		} else if op == sql.EQ {
			return fmt.Sprintf(`{"term": {"%s": %v}}`, comp.Ident, val.Val), nil
		} else if op == sql.NE {
			return fmt.Sprintf(`{"bool": {"must_not": {"term": {"%s": %v}}}}`, comp.Ident, val.Val), nil
		} else {
			return "", fmt.Errorf("unexpected comparison token generating number comparison: %v", op)
		}

	case *sql.StringExpr:

		op := comp.CondOp
		if op == sql.EQ {
			return fmt.Sprintf(`{"term": {"%s": %v}}`, comp.Ident, val.Val), nil
		} else if op == sql.NE {
			return fmt.Sprintf(`{"bool": {"must_not": {"term": {"%s": %v}}}}`, comp.Ident, val.Val), nil
		} else {
			return "", fmt.Errorf("unexpected comparison token generating string comparison: %v", op)
		}

	case *sql.FuncCallExpr:
		return "", fmt.Errorf("function call comparisons not yet supported for: %s", val.Name)

	default:
		return "", fmt.Errorf("unexpected expression type in comparison: %t", val)
	}
}

// genRangeOp expects a comparison token (LT, LE, GT, GE) and returns the
// elasticsearch range operator equivalent, otherwise an empty string.
func genRangeOp(tok sql.Token) string {
	switch tok {
	case sql.LT: return "lt"
	case sql.LE: return "lte"
	case sql.GT: return "gt"
	case sql.GE: return "gte"
	default: return ""
	}
}