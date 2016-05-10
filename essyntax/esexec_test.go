package essyntax


import (
	"testing"

	log "github.com/cihub/seelog"
	T "github.com/oldenbur/sql-parser/testutil"
	. "github.com/smartystreets/goconvey/convey"
	. "github.com/mattbaird/elastigo/lib"
)

func init() { T.ConfigureTestLogger() }

func TestEsExec(t *testing.T) {

	defer log.Flush()

	c := NewConn()
	PopulateTestDB(t, c)
//	defer TearDownTestDB(c)

	Convey("Fields String method\n", t, func() {

	})
}