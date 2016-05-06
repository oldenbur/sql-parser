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
		So(es, ShouldEqual, `{"range": {"numLT": {"lt": 12.3}}}`)
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

	Convey("Test ES conjuctions\n", t, func() {
		es, err := genCondClause(&CondConj{
			Left: &CondComp{Ident:"condAnd1", CondOp: EQ, Val: &StringExpr{Val: `"condAndVal"`}}, Op: AND,
			Right: &CondComp{Ident:"condAnd2", CondOp: EQ, Val: &NumExpr{Val: -9}}})
		So(err, ShouldBeNil)
		So(es, ShouldEqual, `{"bool": {"must": [{"term": {"condAnd1": "condAndVal"}}, {"term": {"condAnd2": -9}}]}}`)
		log.Debug(es)

		es, err = genCondClause(&CondConj{
			Left: &CondComp{Ident:"condOr1", CondOp: EQ, Val: &StringExpr{Val: `"condOrVal"`}}, Op: OR,
			Right: &CondComp{Ident:"condOr2", CondOp: EQ, Val: &NumExpr{Val: 23}}})
		So(err, ShouldBeNil)
		So(es, ShouldEqual, `{"bool": {"should": [{"term": {"condOr1": "condOrVal"}}, {"term": {"condOr2": 23}}]}}`)
		log.Debug(es)

		es, err = genCondClause(&CondConj{
			Left: &CondConj{
				Left: &CondComp{Ident:"c1", CondOp: NE, Val: &StringExpr{Val: `"c1val"`}},
				Op: AND,
				Right: &CondComp{Ident:"c2", CondOp: GE, Val: &NumExpr{Val: 2}}},
			Op: OR,
			Right: &CondConj{
				Left: &CondComp{Ident:"c3", CondOp: LT, Val: &NumExpr{Val: 3}},
				Op: AND,
				Right: &CondConj{
					Left: &CondConj{
						Left: &CondComp{Ident:"c4", CondOp: EQ, Val: &StringExpr{Val: `"c4val"`}},
						Op: OR,
						Right: &CondComp{Ident:"c5", CondOp: EQ, Val: &StringExpr{Val: `"c4val"`}}},
					Op: AND,
					Right: &CondComp{Ident:"c6", CondOp: EQ, Val: &StringExpr{Val: `"c4val"`}}}}})
		So(err, ShouldBeNil)
		So(es, ShouldEqual,
			`{"bool": {"should": [` +
				`{"bool": {"must": [{"bool": {"must_not": {"term": {"c1": "c1val"}}}}, {"range": {"c2": {"gte": 2}}}]}}, ` +
				`{"bool": {"must": [` +
					`{"range": {"c3": {"lt": 3}}}, ` +
					`{"bool": {"must": [` +
						`{"bool": {"should": [{"term": {"c4": "c4val"}}, {"term": {"c5": "c4val"}}]}}, ` +
						`{"term": {"c6": "c4val"}}]}}]}}]}}`)
		log.Debug(es)

	})
}