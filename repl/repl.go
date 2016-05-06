package repl

import (
	"github.com/abiosoft/ishell"
)

func initRepl() {

	shell := ishell.New()
	shell.Println("Sample Interactive Shell")
}
