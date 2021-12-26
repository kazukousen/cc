package main

import (
	"fmt"
	"os"
	"strings"
)

var userIn string
var in string

type tokenKind int

const (
	tokenKindReserved tokenKind = iota
	tokenKindNumber
)

type token struct {
	kind tokenKind
	val  string
	num  int
}

func tokenize() []*token {
	var tokens []*token
	for len(in) > 0 {

		if in[0] == ' ' {
			in = in[1:]
			continue
		}

		if in[0] == '+' || in[0] == '-' {
			tokens = append(tokens, &token{kind: tokenKindReserved, val: string(in[0])})
			in = in[1:]
			continue
		}

		if in[0] >= '0' && in[0] <= '9' {
			n := toInt()
			tokens = append(tokens, &token{kind: tokenKindNumber, num: n})
			continue
		}

		errorAt("unexpected character: " + string(in[0]))
	}
	return tokens
}

func errorAt(msg string) {
	n := len(userIn) - len(in)
	_, _ = fmt.Fprintln(os.Stderr, userIn)
	_, _ = fmt.Fprintln(os.Stderr, strings.Repeat(" ", n), "^")
	_, _ = fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func main() {
	if len(os.Args) != 2 {
		_, _ = fmt.Fprintf(os.Stderr, "the number of arguments is insufficient\n")
		os.Exit(1)
	}

	userIn = os.Args[1]
	in = os.Args[1]

	tokens := tokenize()

	fmt.Printf(`.intel_syntax noprefix
.globl main
main:
	mov rax, %d
`, tokens[0].num)
	tokens = tokens[1:]

	for len(tokens) > 0 {
		switch tokens[0].val {
		case "+":
			fmt.Printf("	add rax, %d\n", tokens[1].num)
			tokens = tokens[2:]
		case "-":
			fmt.Printf("	sub rax, %d\n", tokens[1].num)
			tokens = tokens[2:]
		default:
			_, _ = fmt.Fprintln(os.Stderr, tokens)
			_, _ = fmt.Fprintln(os.Stderr, "[Error] Unexpected token:", tokens[0])
			os.Exit(1)
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
