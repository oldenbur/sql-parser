package eslink

import (
	"testing"

	log "github.com/cihub/seelog"
	T "github.com/oldenbur/sql-parser/testutil"
	. "github.com/smartystreets/goconvey/convey"
)

func init() { T.ConfigureTestLogger() }

func TestESLink(t *testing.T) {

	defer log.Flush()

	Convey("Basic connection\n", t, func() {
		err := conn()
		So(err, ShouldBeNil)
	})

	Convey("Json stree\n", t, func() {
		err := bulkJson()
		So(err, ShouldBeNil)
	})

}
