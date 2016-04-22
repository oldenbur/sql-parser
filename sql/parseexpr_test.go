package sql

import (
	"fmt"
	"strings"
	"testing"

	log "github.com/cihub/seelog"
	T "github.com/oldenbur/sql-parser/testutil"
	. "github.com/smartystreets/goconvey/convey"
)

func init() { T.ConfigureTestLogger() }

func TestFuncCall(t *testing.T) {

	defer log.Flush()

	Convey("Test parsing a string Expr\n", t, func() {
		p := NewParser(strings.NewReader(`"thisIsAString"`))
		e, err := p.parseExpr()
		So(err, ShouldBeNil)
		So(e, ShouldResemble, &StringExpr{Val:`"thisIsAString"`})
		log.Debugf("stringExpr: %v", e)
	})

	Convey("Test parsing an invalid expression\n", t, func() {
		p := NewParser(strings.NewReader(`SELECT 123.456 "anotherString"`))
		_, err := p.parseExpr()
		So(err, ShouldResemble, fmt.Errorf(`parseExpr() expected expression (string, number or function call), got SELECT`))
	})

	Convey("Test parsing an integer\n", t, func() {
		p := NewParser(strings.NewReader(`8765`))
		i, err := p.parseExpr()
		So(err, ShouldBeNil)
		So(i, ShouldResemble, &NumExpr{Val: 8765.0})
		log.Debugf("NumExpr: %s", i)
	})

	Convey("Test parsing a float\n", t, func() {
		p := NewParser(strings.NewReader(`8765.432`))
		f, err := p.parseExpr()
		So(err, ShouldBeNil)
		So(f, ShouldResemble, &NumExpr{Val: 8765.432})
		log.Debugf("NumExpr: %s", f)
	})

	Convey("Test parsing function call without args\n", t, func() {
		p := NewParser(strings.NewReader(`FuncName()`))
		f, err := p.parseExpr()
		So(err, ShouldBeNil)
		So(f, ShouldResemble, &FuncCallExpr{Name:"FuncName", Args: []Expr{}})
		log.Debugf("cond: %s", f)
	})

	Convey("Test parsing function call with one string argument\n", t, func() {
		p := NewParser(strings.NewReader(`FuncName("stringArg")`))
		f, err := p.parseExpr()
		So(err, ShouldBeNil)
		So(f, ShouldResemble, &FuncCallExpr{Name:"FuncName", Args: []Expr{&StringExpr{Val: `"stringArg"`}}})
		log.Debugf("cond: %s", f)
	})

	Convey("Test parsing function call with many arguments\n", t, func() {
		p := NewParser(strings.NewReader(`FuncName("stringArg", -43.21, InnerFunc("innerArg"))`))
		f, err := p.parseFuncCall()
		So(err, ShouldBeNil)
		So(f, ShouldResemble, &FuncCallExpr{Name: "FuncName", Args: []Expr{
			&StringExpr{Val: `"stringArg"`},
			&NumExpr{Val: -43.21},
			&FuncCallExpr{Name: "InnerFunc", Args: []Expr{&StringExpr{Val: `"innerArg"`}}},
		}})
		log.Debugf("cond: %s", f)
	})

	Convey("Test parsing function call with one string argument\n", t, func() {
		p := NewParser(strings.NewReader(`FuncName 123`))
		_, err := p.parseExpr()
		So(err, ShouldResemble, fmt.Errorf(`expected '(', got '123'`))
	})

	Convey("Test parsing function call with one string argument\n", t, func() {
		p := NewParser(strings.NewReader(`FuncName(123`))
		_, err := p.parseExpr()
		So(err, ShouldResemble, fmt.Errorf(`expected COMMA or PAREN_R after FuncName arg 1, got EOF`))
	})

	Convey("Test parsing function call with one string argument\n", t, func() {
		p := NewParser(strings.NewReader(`FuncName(123 "strang")`))
		f, err := p.parseExpr()
		log.Debugf("f: %v", f)
		So(err, ShouldResemble, fmt.Errorf(`expected COMMA or PAREN_R after FuncName arg 1, got "strang"`))
	})



}
