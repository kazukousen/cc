package main

import (
	"fmt"
	"os"
)

var userIn string
var in string
var tokens []*token
var code []*node

func main() {
	if len(os.Args) != 2 {
		_, _ = fmt.Fprintf(os.Stderr, "the number of arguments is insufficient\n")
		os.Exit(1)
	}

	userIn = os.Args[1]

	in = os.Args[1]

	tokenize()
	program()

	fmt.Printf(`.intel_syntax noprefix
.globl main
main:
	push rbp
	mov rbp, rsp
	sub rsp, 208
`)

	for _, c := range code {
		gen(c)
		fmt.Printf("	pop rax\n")
	}

	fmt.Printf(`	mov rsp, rbp
	pop rbp
	ret
`)
}
