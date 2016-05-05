package essyntax

import (
	"fmt"

	"github.com/oldenbur/sql-parser/sql"
)

func ElasticSearchQuery(s *sql.SelectStatement) (string, error) {



	return "", nil
}

func singleIndexQuery(where sql.Cond) (string, error) {

	switch where := where.(type) {
	case *sql.CondComp:
		return genCompClause(where)
	case *sql.CondConj:
		if where.Left == nil {
			return singleIndexQuery(where.Right)
		} else if where.Right == nil {
			return singleIndexQuery(where.Left)
		} else {
			
		}

		return "", nil
	default:
		return "", fmt.Errorf("unexpected singleIndexQuery condition type: %t", where)
	}

}

// genCompClause creates an elasticsearch filter clause for the specified comparison
func genCompClause(where *sql.CondComp) (string, error) {

	switch val := where.Val.(type) {
	case *sql.NumExpr:

		op := where.CondOp
		if op == sql.LT || op == sql.LE || op == sql.GT || op == sql.GE {
			return fmt.Sprintf(`{"range": {"%s": {"%s": %v}}}`, genRangeOp(op), where.Ident, val.Val), nil
		} else if op == sql.EQ {
			return fmt.Sprintf(`{"term": {"%s": %v}}`, where.Ident, val.Val), nil
		} else if op == sql.NE {
			return fmt.Sprintf(`{"bool": {"must_not": {"term": {"%s": %v}}}}`, where.Ident, val.Val), nil
		} else {
			return "", fmt.Errorf("unexpected comparison token generating number comparison: %v", op)
		}

	case *sql.StringExpr:

		op := where.CondOp
		if op == sql.EQ {
			return fmt.Sprintf(`{"term": {"%s": %v}}`, where.Ident, val.Val), nil
		} else if op == sql.NE {
			return fmt.Sprintf(`{"bool": {"must_not": {"term": {"%s": %v}}}}`, where.Ident, val.Val), nil
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