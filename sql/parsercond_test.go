package sql

import (
	"strings"
	"testing"

	log "github.com/cihub/seelog"
	T "github.com/oldenbur/sql-parser/testutil"
	. "github.com/smartystreets/goconvey/convey"
)

func init() { T.ConfigureTestLogger() }

func TestParserCond(t *testing.T) {

	defer log.Flush()

	Convey("Test parsing basic conditions\n", t, func() {
		p := NewParser(strings.NewReader(`A = "a"`))
		c, err := p.parseCondTree()
		So(err, ShouldBeNil)
		So(c, ShouldResemble, &CondComp{Ident:"A", CondOp: EQ, Val: `"a"`})
		log.Debugf("cond: %s", c)

		p = NewParser(strings.NewReader(`t1.A != "a" AND t2.B >= -2345`))
		c, err = p.parseCondTree()
		So(err, ShouldBeNil)
		So(c, ShouldResemble, &CondConj{
			Left: &CondComp{Ident:"t1.A", CondOp: NE, Val: `"a"`},
			Op: AND,
			Right: &CondComp{Ident:"t2.B", CondOp: GE, Val: `-2345`}})
		log.Debugf("cond: %s", c)

		p = NewParser(strings.NewReader(`t1.A = "aa aa" AND t2.B <= -.23 AND C = "c" AND t1.t2.D = -09`))
		c, err = p.parseCondTree()
		So(err, ShouldBeNil)
		chk := &CondConj{
			Left: &CondComp{Ident:"t1.A", CondOp: EQ, Val: `"aa aa"`}, Op: AND,
			Right: &CondConj{
				Left: &CondComp{Ident:"t2.B", CondOp: LE, Val: `-.23`}, Op: AND,
				Right: &CondConj{
					Left: &CondComp{Ident:"C", CondOp: EQ, Val: `"c"`}, Op: AND,
					Right: &CondComp{Ident:"t1.t2.D", CondOp: EQ, Val: `-09`}}}}
		So(c.String(), ShouldEqual, chk.String())
		log.Debugf("cond: %s", c)

		p = NewParser(strings.NewReader(`(t1.A != "a" AND t2.B >= -2345)`))
		c, err = p.parseCondTree()
		So(err, ShouldBeNil)
		So(c, ShouldResemble, &CondConj{
			Left: &CondComp{Ident:"t1.A", CondOp: NE, Val: `"a"`},
			Op: AND,
			Right: &CondComp{Ident:"t2.B", CondOp: GE, Val: `-2345`}})
		log.Debugf("cond: %s", c)

		p = NewParser(strings.NewReader(`(t1.A != "a" AND t2.B >= -2345) OR t3.C = "cccc  "`))
		c, err = p.parseCondTree()
		So(err, ShouldBeNil)
		chk = &CondConj{
			Left: &CondConj{
				Left: &CondComp{Ident:"t1.A", CondOp: NE, Val: `"a"`},
				Op: AND,
				Right: &CondComp{Ident:"t2.B", CondOp: GE, Val: `-2345`}},
			Op: OR,
			Right: &CondComp{Ident:"t3.C", CondOp: EQ, Val: `"cccc  "`}}
		So(c.String(), ShouldEqual, chk.String())
		log.Debugf("cond: %s", c)

		p = NewParser(strings.NewReader(`(t1.A != "a" AND t2.B >= -2345) OR (C < 5 AND D = 'd')`))
		c, err = p.parseCondTree()
		So(err, ShouldBeNil)
		chk = &CondConj{
			Left: &CondConj{
				Left: &CondComp{Ident:"t1.A", CondOp: NE, Val: `"a"`},
				Op: AND,
				Right: &CondComp{Ident:"t2.B", CondOp: GE, Val: `-2345`}},
			Op: OR,
			Right: &CondConj{
				Left: &CondComp{Ident:"C", CondOp: LT, Val: `5`},
				Op: AND,
				Right: &CondComp{Ident:"D", CondOp: EQ, Val: `'d'`}}}
		So(c.String(), ShouldEqual, chk.String())
		log.Debugf("cond: %s", c)

		p = NewParser(strings.NewReader(`(t1.A != "a" AND t2.B >= -2345) OR (C < 5 AND (D = 'd' OR E = 'e') AND F = 'f')`))
		c, err = p.parseCondTree()
		So(err, ShouldBeNil)
		chk = &CondConj{
			Left: &CondConj{
				Left: &CondComp{Ident:"t1.A", CondOp: NE, Val: `"a"`},
				Op: AND,
				Right: &CondComp{Ident:"t2.B", CondOp: GE, Val: `-2345`}},
			Op: OR,
			Right: &CondConj{
				Left: &CondComp{Ident:"C", CondOp: LT, Val: `5`},
				Op: AND,
				Right: &CondConj{
					Left: &CondConj{
						Left: &CondComp{Ident:"D", CondOp: EQ, Val: `'d'`},
						Op: OR,
						Right: &CondComp{Ident:"E", CondOp: EQ, Val: `'e'`}},
					Op: AND,
					Right: &CondComp{Ident:"F", CondOp: EQ, Val: `'f'`}}}}
		So(c.String(), ShouldEqual, chk.String())
		log.Debugf("cond: %s", c)

	})



}