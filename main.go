package main

import (
	"fmt"
	"os"
)

var userIn string
var in string
var tokens []*token
var code []*node
var locals []*obj

func main() {
	if len(os.Args) != 2 {
		_, _ = fmt.Fprintf(os.Stderr, "the number of arguments is insufficient\n")
		os.Exit(1)
	}

	userIn = os.Args[1]

	in = os.Args[1]

	tokenize()
	program()

	offset := 0
	for i := len(locals) - 1; i >= 0; i-- {
		v := locals[i]
		offset += 8
		v.offset = -offset
	}

	fmt.Printf(`.intel_syntax noprefix
.globl main
main:
	push rbp
	mov rbp, rsp
	sub rsp, %d
`, len(locals)*8)

	for _, c := range code {
		gen(c)
	}

	fmt.Printf(`.Lreturn:
	mov rsp, rbp
	pop rbp
	ret
`)
}
