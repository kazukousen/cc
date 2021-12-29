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

func equalToken(kind tokenKind) bool {
	return len(tokens) > 0 && tokens[0].kind == kind
}

func consumeToken(kind tokenKind) *token {
	if equalToken(kind) {
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
	body   statement
}

type expression interface {
	isExpr()
}

type statement interface {
	isStmt()
}

type binaryNode struct {
	op  string
	lhs expression
	rhs expression
}

type assignNode binaryNode

type unaryNode struct {
	child expression
}

type exprStmtNode unaryNode
type returnStmtNode unaryNode
type addrNode unaryNode
type derefNode unaryNode

type intLit struct {
	ty  *typ
	val int
}

type obj struct {
	ty     *typ
	name   string
	offset int
}

type funcCallNode struct {
	name string
	args []expression
}

func (binaryNode) isExpr()     {}
func (assignNode) isExpr()     {}
func (unaryNode) isExpr()      {}
func (exprStmtNode) isStmt()   {}
func (returnStmtNode) isStmt() {}
func (addrNode) isExpr()       {}
func (derefNode) isExpr()      {}
func (intLit) isExpr()         {}
func (obj) isExpr()            {}
func (funcCallNode) isExpr()   {}

type ifStmtNode struct {
	cond expression
	then statement
	els  statement
}

type forStmtNode struct {
	ini  expression
	cond expression
	step expression
	then statement
}

type blockStmtNode struct {
	code []statement
}

func (ifStmtNode) isStmt()    {}
func (forStmtNode) isStmt()   {}
func (blockStmtNode) isStmt() {}

func findLocal(name string) *obj {
	for i := range locals {
		lv := locals[i]
		if lv.name == name {
			return lv
		}
	}
	return nil
}

func newNodeLocal(name string) *obj {

	if lv := findLocal(name); lv != nil {
		return lv
	}

	lv := &obj{
		name: name,
	}
	locals = append(locals, lv)

	return lv
}

func program() (funcs []*function) {
	for len(tokens) > 0 {
		funcs = append(funcs, funcDecl())
	}
	return
}

func funcDecl() *function {

	locals = []*obj{}

	ty := declSpec()
	v, args := declarator(ty)

	f := &function{
		name: v.name,
		args: args,
	}

	expect("{")
	f.body = compoundStmt()
	f.locals = locals

	return f
}

func declSpec() *typ {
	tok := consumeToken(tokenKindType)
	if tok == nil {
		_, _ = fmt.Fprintln(os.Stderr, "Expect an type:", tokens[0].val)
		os.Exit(1)
	}

	return newLiteralType(tok.val)
}

func declarator(ty *typ) (*obj, []*obj) {

	for consume("*") {
		ty = pointerTo(ty)
	}

	tok := consumeToken(tokenKindIdent)
	if tok == nil {
		_, _ = fmt.Fprintln(os.Stderr, "Expect an identifier in declarator:", tokens[0].val)
		os.Exit(1)
	}
	val := newNodeLocal(tok.val)
	val.ty = ty

	var args []*obj
	if consume("(") {
		for i := 0; !consume(")"); i++ {
			if i > 0 {
				expect(",")
			}
			ty := declSpec()
			v, ps := declarator(ty)
			args = append(args, v)
			args = append(args, ps...)
		}
	}

	return val, args
}

func stmt() statement {
	if consume("return") {
		ret := returnStmtNode{child: expr()}
		expect(";")
		return ret
	} else if consume("{") {
		return compoundStmt()
	} else if consume("if") {
		return ifStmt()
	} else if consume("while") {
		return whileStmt()
	} else if consume("for") {
		return forStmt()
	} else {
		ret := expr()
		expect(";")
		return exprStmtNode{child: ret}
	}
}

func compoundStmt() statement {
	ret := blockStmtNode{code: []statement{}}
	for !consume("}") {
		if equalToken(tokenKindType) {
			ty := declSpec()
			for i := 0; !consume(";"); i++ {
				if i > 0 {
					expect(",")
				}
				v, _ := declarator(ty)
				if consume("=") {
					n := assignNode{op: "=", lhs: v, rhs: assign()}
					ret.code = append(ret.code, exprStmtNode{child: n})
				}
			}
		} else {
			ret.code = append(ret.code, stmt())
		}
	}
	return ret
}

func ifStmt() statement {
	expect("(")
	cond := expr()
	expect(")")
	then := stmt()
	if consume("else") {
		els := stmt()
		return ifStmtNode{cond: cond, then: then, els: els}
	} else {
		return ifStmtNode{cond: cond, then: then, els: nil}
	}
}

func whileStmt() statement {
	expect("(")
	cond := expr()
	expect(")")
	then := stmt()
	return forStmtNode{ini: nil, cond: cond, step: nil, then: then}
}

func forStmt() statement {
	expect("(")
	var ini expression
	for !consume(";") {
		ini = expr()
	}
	var cond expression
	for !consume(";") {
		cond = expr()
	}
	var step expression
	for !consume(")") {
		step = expr()
	}
	then := stmt()
	return forStmtNode{ini: ini, cond: cond, step: step, then: then}
}

func expr() expression {
	return assign()
}

func assign() expression {
	ret := equality()
	if consume("=") {
		ret = assignNode{op: "=", lhs: ret, rhs: assign()}
	}
	return ret
}

func equality() expression {
	ret := relational()
	for {
		switch {
		case consume("=="):
			ret = binaryNode{op: "==", lhs: ret, rhs: relational()}
		case consume("!="):
			ret = binaryNode{op: "!=", lhs: ret, rhs: relational()}
		default:
			return ret
		}
	}
}

func relational() expression {
	ret := add()
	for {
		switch {
		case consume("<"):
			ret = binaryNode{op: "<", lhs: ret, rhs: add()}
		case consume("<="):
			ret = binaryNode{op: "<=", lhs: ret, rhs: add()}
		case consume(">"):
			ret = binaryNode{op: "<", lhs: add(), rhs: ret}
		case consume(">="):
			ret = binaryNode{op: "<=", lhs: add(), rhs: ret}
		default:
			return ret
		}
	}
}

func add() expression {
	ret := mul()
	for {
		switch {
		case consume("+"):
			ret = binaryNode{op: "+", lhs: ret, rhs: mul()}
		case consume("-"):
			ret = binaryNode{op: "-", lhs: ret, rhs: mul()}
		default:
			return ret
		}
	}
}

func mul() expression {
	ret := unary()
	for {
		switch {
		case consume("*"):
			ret = binaryNode{op: "*", lhs: ret, rhs: unary()}
		case consume("/"):
			ret = binaryNode{op: "/", lhs: ret, rhs: unary()}
		default:
			return ret
		}
	}
}

func unary() expression {
	switch {
	case consume("-"):
		return binaryNode{op: "-", lhs: intLit{val: 0}, rhs: unary()}
	case consume("+"):
		return unary()
	case consume("&"):
		return addrNode{child: unary()}
	case consume("*"):
		return derefNode{child: unary()}
	default:
		return primary()
	}
}

func primary() expression {
	if consume("(") {
		ret := expr()
		expect(")")
		return ret
	}

	if tok := consumeToken(tokenKindIdent); tok != nil {
		if consume("(") {
			return funcCallNode{name: tok.val, args: callArgs()}
		} else {
			return newNodeLocal(tok.val)
		}
	}

	return num()
}

func callArgs() (args []expression) {
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

func num() expression {
	ret := intLit{val: tokens[0].num}
	advance()
	return ret
}
