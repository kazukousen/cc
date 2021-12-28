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
		_, _ = fmt.Fprintln(os.Stderr, "Unexpected token:", tokens[0].val, "want:", c)
		os.Exit(1)
	}
}

func advance() {
	tokens = tokens[1:]
}

type function struct {
	name   string
	locals []*obj
	args   []*obj
	body   *node
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
	nodeKindVar
	nodeKindNum
	nodeKindBlock
	nodeKindIf
	nodeKindFor
	nodeKindCall
	nodeKindAddr
	nodeKindDeref
	nodeKindReturn
)

type node struct {
	kind nodeKind

	lhs *node
	rhs *node

	num      int
	variable *obj

	ini  *node
	cond *node
	step *node
	then *node
	els  *node

	code []*node

	name string
	args []*node
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

func newNodeIf(cond *node, then *node, els *node) *node {
	return &node{
		kind: nodeKindIf,
		cond: cond,
		then: then,
		els:  els,
	}
}

func newNodeFor(ini *node, cond *node, step *node, then *node) *node {
	return &node{
		kind: nodeKindFor,
		ini:  ini,
		cond: cond,
		step: step,
		then: then,
	}
}

type obj struct {
	name   string
	offset int
}

func findLocal(name string) *obj {
	for i := range locals {
		lv := locals[i]
		if lv.name == name {
			return lv
		}
	}
	return nil
}

func newNodeLocal(name string) *node {

	if lv := findLocal(name); lv != nil {
		return &node{
			kind:     nodeKindVar,
			variable: lv,
		}
	}

	variable := &obj{
		name: name,
	}
	locals = append(locals, variable)

	return &node{
		kind:     nodeKindVar,
		variable: variable,
	}
}

func program() (funcs []*function) {
	for len(tokens) > 0 {
		funcs = append(funcs, funcDecl())
	}
	return
}

func funcDecl() *function {

	locals = []*obj{}

	declSpec()
	v, args := declarator()

	f := &function{
		name: v.name,
		args: args,
	}

	expect("{")
	f.body = &node{kind: nodeKindBlock, code: []*node{}}
	for !consume("}") {
		f.body.code = append(f.body.code, stmt())
	}
	f.locals = locals

	return f
}

func declSpec() {
	expect("int")
}

func declarator() (*obj, []*obj) {

	tok := consumeIdent()
	if tok == nil {
		_, _ = fmt.Fprintln(os.Stderr, "Expect an identifier in declarator", tokens[0].val)
		os.Exit(1)
	}
	val := newNodeLocal(tok.val).variable

	var args []*obj
	if consume("(") {
		if consume(")") {
			return val, args
		}

		for {
			declSpec()
			v, _ := declarator()
			args = append(args, v)
			if !consume(",") {
				break
			}
		}
		expect(")")
	}

	return val, args
}

func funcArgs() (args []*obj) {
	if !consume(")") {
		return
	}
	expect(")")
	return
}

func stmt() *node {
	if consume("return") {
		ret := newNode(nodeKindReturn, expr(), nil)
		expect(";")
		return ret
	} else if consume("{") {
		ret := &node{kind: nodeKindBlock, code: []*node{}}
		for !consume("}") {
			ret.code = append(ret.code, stmt())
		}
		return ret
	} else if consume("if") {
		return ifStmt()
	} else if consume("while") {
		return whileStmt()
	} else if consume("for") {
		return forStmt()
	} else {
		ret := expr()
		expect(";")
		return ret
	}
}

func ifStmt() *node {
	expect("(")
	cond := expr()
	expect(")")
	then := stmt()
	if consume("else") {
		els := stmt()
		return newNodeIf(cond, then, els)
	} else {
		return newNodeIf(cond, then, nil)
	}
}

func whileStmt() *node {
	expect("(")
	cond := expr()
	expect(")")
	then := stmt()
	return newNodeFor(nil, cond, nil, then)
}

func forStmt() *node {
	expect("(")
	var ini *node
	for !consume(";") {
		ini = expr()
	}
	var cond *node
	for !consume(";") {
		cond = expr()
	}
	var step *node
	for !consume(")") {
		step = expr()
	}
	then := stmt()
	return newNodeFor(ini, cond, step, then)
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
	case consume("&"):
		ret := unary()
		return newNode(nodeKindAddr, ret, nil)
	case consume("*"):
		ret := unary()
		return newNode(nodeKindDeref, ret, nil)
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
		if consume("(") {
			return &node{kind: nodeKindCall, name: tok.val, args: callArgs()}
		} else {
			return newNodeLocal(tok.val)
		}
	}

	return num()
}

func callArgs() (args []*node) {
	if consume(")") {
		return
	}
	args = append(args, assign())
	for consume(",") {
		args = append(args, assign())
	}
	expect(")")
	return
}

func num() *node {
	ret := newNodeNum(tokens[0].num)
	advance()
	return ret
}
