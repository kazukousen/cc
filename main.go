package main

import (
	"fmt"
	"os"
)

var in string

func main() {
	if len(os.Args) != 2 {
		_, _ = fmt.Fprintf(os.Stderr, "the number of arguments is insufficient\n")
		os.Exit(1)
	}

	in = os.Args[1]

	fmt.Printf(`.intel_syntax noprefix
.globl main
main:
	mov rax, %d
`, toInt())

	for len(in) > 0 {
		switch in[0] {
		case '+':
			in = in[1:]
			fmt.Printf("	add rax, %d\n", toInt())
		case '-':
			in = in[1:]
			fmt.Printf("	sub rax, %d\n", toInt())
		}
	}

	fmt.Printf("	ret\n")
}

func toInt() int {
	var ret int
	for len(in) > 0 && in[0] >= '0' && in[0] <= '9' {
		ret = ret*10 + int(in[0]-'0')
		in = in[1:]
	}
	return ret
}
