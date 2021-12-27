package main

import (
	"fmt"
	"os"
	"strings"
)

var userIn string
var in string
var tokens []*token

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

		if in[0] == '+' || in[0] == '-' || in[0] == '*' || in[0] == '/' || in[0] == '(' || in[0] == ')' {
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

func consume(c string) bool {
	if len(tokens) > 0 && tokens[0].val == c {
		tokens = tokens[1:]
		return true
	}
	return false
}

func expect(c string) {
	if !consume(c) {
		_, _ = fmt.Fprintln(os.Stderr, "Unexpected token:", tokens[0].val)
		os.Exit(1)
	}
}

func advance() {
	tokens = tokens[1:]
}

type nodeKind int

const (
	nodeKindAdd nodeKind = iota
	nodeKindSub
	nodeKindMul
	nodeKindDiv
	nodeKindNum
)

type node struct {
	kind nodeKind
	lhs  *node
	rhs  *node
	num  int
}

func newNode(kind nodeKind, left *node, right *node) *node {
	return &node{
		kind: kind,
		lhs:  left,
		rhs:  right,
	}
}

func newNodeNum(num int) *node {
	return &node{
		kind: nodeKindNum,
		num:  num,
	}
}

func expr() *node {
	ret := mul()
	for {
		switch {
		case consume("+"):
			ret = newNode(nodeKindAdd, ret, mul())
		case consume("-"):
			ret = newNode(nodeKindSub, ret, mul())
		default:
			return ret
		}
	}
}

func mul() *node {
	ret := primary()
	for {
		switch {
		case consume("*"):
			ret = newNode(nodeKindMul, ret, primary())
		case consume("/"):
			ret = newNode(nodeKindDiv, ret, primary())
		default:
			return ret
		}
	}
}

func primary() *node {
	if consume("(") {
		ret := expr()
		expect(")")
		return ret
	}
	return num()
}

func num() *node {
	ret := newNodeNum(tokens[0].num)
	advance()
	return ret
}

func gen(n *node) {

	if n.kind == nodeKindNum {
		fmt.Printf("	push %d\n", n.num)
		return
	}

	if n.lhs != nil {
		gen(n.lhs)
	}
	if n.rhs != nil {
		gen(n.rhs)
	}

	fmt.Printf("	pop rdi\n")
	fmt.Printf("	pop rax\n")

	switch n.kind {
	case nodeKindAdd:
		fmt.Printf("	add rax,rdi\n")
	case nodeKindSub:
		fmt.Printf("	sub rax,rdi\n")
	case nodeKindMul:
		fmt.Printf("	imul rax,rdi\n")
	case nodeKindDiv:
		fmt.Printf("	cqo\n")
		fmt.Printf("	idiv rdi\n")
	}

	fmt.Printf("	push rax\n")
}

func main() {
	if len(os.Args) != 2 {
		_, _ = fmt.Fprintf(os.Stderr, "the number of arguments is insufficient\n")
		os.Exit(1)
	}

	userIn = os.Args[1]
	in = os.Args[1]

	tokens = tokenize()

	fmt.Printf(`.intel_syntax noprefix
.globl main
main:
`)

	gen(expr())

	fmt.Printf("	pop rax\n")
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
