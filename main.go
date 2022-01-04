package main

import (
	"fmt"
	"os"
)

var tokens []*token
var locals []*obj

func main() {
	if len(os.Args) != 2 {
		_, _ = fmt.Fprintf(os.Stderr, "the number of arguments is insufficient\n")
		os.Exit(1)
	}

	tokenizeFile(os.Args[1])
	prog := parse()

	codegen(prog)
}
