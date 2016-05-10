package main

import (
//	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/oldenbur/sql-parser/sql"
	"github.com/oldenbur/sql-parser/essyntax"
)

var printcmd, norun bool

func main() {

	app := cli.NewApp()
	app.Name = "essql"
	app.Usage = "SQL-like interpreter for elasticsearch queries"
	app.Action = run
	app.Flags = []cli.Flag {
		cli.BoolFlag{
			Name: "print, p",
			Destination: &printcmd,
			Usage: "if true, print the generated elasticsearch REST API call",
		},
		cli.BoolFlag{
			Name: "norun, n",
			Destination: &norun,
			Usage: "if true, do not execute the command",
		},
	}
	app.Run(os.Args)
}

func run(c *cli.Context) error {

	fmt.Printf("essql\nArgs: [%v]\nprint: %v\nnorun: %v\n", strings.Join(c.Args(), ","), printcmd, norun)

	sel, err := sql.NewParser(strings.NewReader(c.Args()[0])).Parse()
	if err != nil {
		return err
	}

	err = essyntax.ExecSearch(sel, printcmd)
	if err != nil {
		return err
	}

	return nil
}
