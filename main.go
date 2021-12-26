package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		_, _ = fmt.Fprintf(os.Stderr, "the number of arguments is insufficient\n")
		os.Exit(1)
	}
	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "the argument must be number\n")
	}

	fmt.Printf(`.intel_syntax noprefix
.globl main
main:
	mov rax, %d
	ret
`, n)
}
