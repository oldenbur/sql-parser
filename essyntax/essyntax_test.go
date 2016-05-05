package essyntax

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	T "github.com/oldenbur/sql-parser/testutil"
	. "github.com/oldenbur/sql-parser/sql"
	log "github.com/cihub/seelog"
)

func init() { T.ConfigureTestLogger() }

func TestQuery(t *testing.T) {

	defer log.Flush()

	Convey("Test getRangeOp\n", t, func() {
		So(genRangeOp(LT), ShouldEqual, "lt")
		So(genRangeOp(LE), ShouldEqual, "lte")
		So(genRangeOp(GT), ShouldEqual, "gt")
		So(genRangeOp(GE), ShouldEqual, "gte")
	})

	Convey("Test ES comparisons\n", t, func() {

		es, err := genCompClause(&CondComp{Ident:"numLT", CondOp: LT, Val: &NumExpr{Val: 12.3}})
		So(err, ShouldBeNil)
		So(es, ShouldEqual, `{"range": {"lt": {"numLT": 12.3}}}`)
		log.Debug(es)

		es, err = genCompClause(&CondComp{Ident:"strEQ", CondOp: EQ, Val: &NumExpr{Val: 23.4}})
		So(err, ShouldBeNil)
		So(es, ShouldEqual, `{"term": {"strEQ": 23.4}}`)
		log.Debug(es)

		es, err = genCompClause(&CondComp{Ident:"strNE", CondOp: NE, Val: &NumExpr{Val: 34.5}})
		So(err, ShouldBeNil)
		So(es, ShouldEqual, `{"bool": {"must_not": {"term": {"strNE": 34.5}}}}`)
		log.Debug(es)

		_, err = genCompClause(&CondComp{Ident:"strP", CondOp: PAREN_R, Val: &NumExpr{Val: 45.6}})
		So(err, ShouldResemble, fmt.Errorf("unexpected comparison token generating number comparison: PAREN_R"))

		es, err = genCompClause(&CondComp{Ident:"strEQ", CondOp: EQ, Val: &StringExpr{Val: `"strEQval"`}})
		So(err, ShouldBeNil)
		So(es, ShouldEqual, `{"term": {"strEQ": "strEQval"}}`)
		log.Debug(es)

		es, err = genCompClause(&CondComp{Ident:"strNE", CondOp: NE, Val: &StringExpr{Val: `"strNEval"`}})
		So(err, ShouldBeNil)
		So(es, ShouldEqual, `{"bool": {"must_not": {"term": {"strNE": "strNEval"}}}}`)
		log.Debug(es)

		_, err = genCompClause(&CondComp{Ident:"strGT", CondOp: GT, Val: &StringExpr{Val: `"strGTval"`}})
		So(err, ShouldResemble, fmt.Errorf("unexpected comparison token generating string comparison: GT"))

	})
}