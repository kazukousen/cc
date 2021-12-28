package main

import (
	"fmt"
	"os"
)

var userIn string
var in string
var tokens []*token
var locals []*obj

func main() {
	if len(os.Args) != 2 {
		_, _ = fmt.Fprintf(os.Stderr, "the number of arguments is insufficient\n")
		os.Exit(1)
	}

	userIn = os.Args[1]

	in = os.Args[1]

	tokenize()
	funcs := program()

	codegen(funcs)
}
