package main

import (
	"fmt"
	"os"
)

func consume(c string) bool {
	if len(tokens) > 0 && tokens[0].val == c {
		tokens = tokens[1:]
		return true
	}
	return false
}

func consumeIdent() *token {
	if len(tokens) > 0 && tokens[0].kind == tokenKindIdent {
		tok := tokens[0]
		tokens = tokens[1:]
		return tok
	}
	return nil
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
	nodeKindEq     // ==
	nodeKindNe     // !=
	nodeKindLt     // <
	nodeKindLe     // <=
	nodeKindAssign // =
	nodeKindLocal
	nodeKindNum
)

type node struct {
	kind   nodeKind
	lhs    *node
	rhs    *node
	num    int
	offset int
}

func newNode(kind nodeKind, left *node, right *node) *node {
	return &node{
		kind: kind,
		lhs:  left,
		rhs:  right,
	}
}

func newNodeLocal(name string) *node {
	return &node{
		kind:   nodeKindLocal,
		offset: int(name[0]-'a'+1) * 8,
	}
}

func newNodeNum(num int) *node {
	return &node{
		kind: nodeKindNum,
		num:  num,
	}
}

func program() {
	for len(tokens) > 0 {
		code = append(code, stmt())
	}
}

func stmt() *node {
	ret := expr()
	expect(";")
	return ret
}

func expr() *node {
	return assign()
}

func assign() *node {
	ret := equality()
	if consume("=") {
		ret = newNode(nodeKindAssign, ret, assign())
	}
	return ret
}

func equality() *node {
	ret := relational()
	for {
		switch {
		case consume("=="):
			ret = newNode(nodeKindEq, ret, relational())
		case consume("!="):
			ret = newNode(nodeKindNe, ret, relational())
		default:
			return ret
		}
	}
}

func relational() *node {
	ret := add()
	for {
		switch {
		case consume("<"):
			ret = newNode(nodeKindLt, ret, add())
		case consume("<="):
			ret = newNode(nodeKindLe, ret, add())
		case consume(">"):
			ret = newNode(nodeKindLt, add(), ret)
		case consume(">="):
			ret = newNode(nodeKindLe, add(), ret)
		default:
			return ret
		}
	}
}

func add() *node {
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
	ret := unary()
	for {
		switch {
		case consume("*"):
			ret = newNode(nodeKindMul, ret, unary())
		case consume("/"):
			ret = newNode(nodeKindDiv, ret, unary())
		default:
			return ret
		}
	}
}

func unary() *node {
	switch {
	case consume("-"):
		return newNode(nodeKindSub, newNodeNum(0), unary())
	case consume("+"):
		return unary()
	default:
		return primary()
	}
}

func primary() *node {
	if consume("(") {
		ret := expr()
		expect(")")
		return ret
	}

	if tok := consumeIdent(); tok != nil {
		return newNodeLocal(tok.val)
	}

	return num()
}

func num() *node {
	ret := newNodeNum(tokens[0].num)
	advance()
	return ret
}
