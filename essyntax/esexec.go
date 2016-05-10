package essyntax

import (

	"github.com/oldenbur/sql-parser/sql"
	. "github.com/mattbaird/elastigo/lib"
	"fmt"
)

// see: https://golang.org/pkg/text/tabwriter/

func ExecSearch(search *sql.SelectStatement, printcmd bool) error {

	c := NewConn()

	condStr, err := genCondClause(search.WhereCond)
	if err != nil {
		return err
	}
	queryStr := fmt.Sprintf(`{"filter": %s}`, condStr)

	if len(search.TableList) != 1 {
		return fmt.Errorf("only searching a single index is supported, got: %v", search.TableList)
	}
	url := fmt.Sprintf(`/%s/_search`, search.TableList[0])

	args := map[string]interface{}{"pretty":true}

	if printcmd {
		fmt.Printf(`POST %s\n%s\n`, url, queryStr)
	}

	r, err := c.DoCommand("POST", url, args, queryStr)
	if err != nil {
		return err
	}

	fmt.Printf("result: %v\n", string(r))

	return nil
}