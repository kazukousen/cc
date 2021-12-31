package main

import (
	"fmt"
	"os"
)

func consume(c string) bool {
	if len(tokens) > 0 && tokens[0].val == c {
		advance()
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
		advance()
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
	name      string
	locals    []*obj
	args      []*obj
	body      statement
	stackSize int
}

type expression interface {
	isExpr()
	getType() *typ
	setType(ty *typ)
}

type statement interface {
	isStmt()
}

type binaryNode struct {
	ty  *typ
	op  string
	lhs expression
	rhs expression
}

type assignNode binaryNode

type unaryNode struct {
	ty    *typ
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
	ty   *typ
	name string
	args []expression
}

func (*binaryNode) isExpr()     {}
func (*assignNode) isExpr()     {}
func (*unaryNode) isExpr()      {}
func (*exprStmtNode) isStmt()   {}
func (*returnStmtNode) isStmt() {}
func (*addrNode) isExpr()       {}
func (*derefNode) isExpr()      {}
func (*intLit) isExpr()         {}
func (*obj) isExpr()            {}
func (*funcCallNode) isExpr()   {}

func (n *binaryNode) getType() *typ     { return n.ty }
func (n *assignNode) getType() *typ     { return n.ty }
func (n *unaryNode) getType() *typ      { return n.ty }
func (n *exprStmtNode) getType() *typ   { return n.ty }
func (n *returnStmtNode) getType() *typ { return n.ty }
func (n *addrNode) getType() *typ       { return n.ty }
func (n *derefNode) getType() *typ      { return n.ty }
func (n *intLit) getType() *typ         { return n.ty }
func (n *obj) getType() *typ            { return n.ty }
func (n *funcCallNode) getType() *typ   { return n.ty }

func (n *binaryNode) setType(ty *typ)     { n.ty = ty }
func (n *assignNode) setType(ty *typ)     { n.ty = ty }
func (n *unaryNode) setType(ty *typ)      { n.ty = ty }
func (n *exprStmtNode) setType(ty *typ)   { n.ty = ty }
func (n *returnStmtNode) setType(ty *typ) { n.ty = ty }
func (n *addrNode) setType(ty *typ)       { n.ty = ty }
func (n *derefNode) setType(ty *typ)      { n.ty = ty }
func (n *intLit) setType(ty *typ)         { n.ty = ty }
func (n *obj) setType(ty *typ)            { n.ty = ty }
func (n *funcCallNode) setType(ty *typ)   { n.ty = ty }

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

func (*ifStmtNode) isStmt()    {}
func (*forStmtNode) isStmt()   {}
func (*blockStmtNode) isStmt() {}

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

// funcDecl = declspec declarator "{" compoundStmt
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
	addType(f.body)
	f.stackSize = calcStackSize(f.locals)

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
	val.setType(ty)

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

// compoundStmt = (declaration | stmt)* "}"
func compoundStmt() statement {
	ret := &blockStmtNode{code: []statement{}}
	for !consume("}") {
		if equalToken(tokenKindType) {
			ret.code = append(ret.code, declaration()...)
		} else {
			ret.code = append(ret.code, stmt())
		}
	}
	return ret
}

// declaration = declspec declarator ("=" expr)? ("," declarator ("=" expr)?)*)? ";"
func declaration() []statement {
	var ret []statement
	ty := declSpec()
	for i := 0; ; i++ {
		if i > 0 {
			expect(",")
		}
		v, _ := declarator(ty)
		if consume("=") {
			n := &assignNode{op: "=", lhs: v, rhs: expr()}
			ret = append(ret, &exprStmtNode{child: n})
		}
		if consume(";") {
			break
		}
	}
	return ret
}

func stmt() statement {
	if consume("return") {
		ret := &returnStmtNode{child: expr()}
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
		return &exprStmtNode{child: ret}
	}
}

func ifStmt() statement {
	expect("(")
	cond := expr()
	expect(")")
	then := stmt()
	if consume("else") {
		els := stmt()
		return &ifStmtNode{cond: cond, then: then, els: els}
	} else {
		return &ifStmtNode{cond: cond, then: then, els: nil}
	}
}

func whileStmt() statement {
	expect("(")
	cond := expr()
	expect(")")
	then := stmt()
	return &forStmtNode{ini: nil, cond: cond, step: nil, then: then}
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
	return &forStmtNode{ini: ini, cond: cond, step: step, then: then}
}

// expr = assign
func expr() expression {
	return assign()
}

// assign = equality ("=" assign)?
func assign() expression {
	ret := equality()
	if consume("=") {
		ret = &assignNode{op: "=", lhs: ret, rhs: assign()}
	}
	return ret
}

// equality = relational ("==" relational | "!=" relational)*
func equality() expression {
	ret := relational()
	for {
		switch {
		case consume("=="):
			ret = &binaryNode{op: "==", lhs: ret, rhs: relational()}
		case consume("!="):
			ret = &binaryNode{op: "!=", lhs: ret, rhs: relational()}
		default:
			return ret
		}
	}
}

// relational = add ("<" add | "<=" add | ">" add | ">=" add)*
func relational() expression {
	ret := add()
	for {
		switch {
		case consume("<"):
			ret = &binaryNode{op: "<", lhs: ret, rhs: add()}
		case consume("<="):
			ret = &binaryNode{op: "<=", lhs: ret, rhs: add()}
		case consume(">"):
			ret = &binaryNode{op: "<", lhs: add(), rhs: ret}
		case consume(">="):
			ret = &binaryNode{op: "<=", lhs: add(), rhs: ret}
		default:
			return ret
		}
	}
}

// add = mul ("+" mul | "-" mul)*
func add() expression {
	ret := mul()
	for {
		switch {
		case consume("+"):
			ret = newAddBinary(ret, mul())
		case consume("-"):
			ret = newSubBinary(ret, mul())
		default:
			return ret
		}
	}
}

// mul = unary ("*" unary | "/" unary)*
func mul() expression {
	ret := unary()
	for {
		switch {
		case consume("*"):
			ret = &binaryNode{op: "*", lhs: ret, rhs: unary()}
		case consume("/"):
			ret = &binaryNode{op: "/", lhs: ret, rhs: unary()}
		default:
			return ret
		}
	}
}

// unary = ("-" | "+" | "&" | "*") unary | primary
func unary() expression {
	switch {
	case consume("-"):
		return &binaryNode{op: "-", lhs: &intLit{val: 0}, rhs: unary()}
	case consume("+"):
		return unary()
	case consume("&"):
		return &addrNode{child: unary()}
	case consume("*"):
		return &derefNode{child: unary()}
	default:
		return primary()
	}
}

// primary = "(" expr ")" | ident ("(" callArgs)? | num
func primary() expression {
	if consume("(") {
		ret := expr()
		expect(")")
		return ret
	}

	if tok := consumeToken(tokenKindIdent); tok != nil {
		if consume("(") {
			return &funcCallNode{name: tok.val, args: callArgs(), ty: newLiteralType("int")}
		} else {
			lv := findLocal(tok.val)
			if lv == nil {
				_, _ = fmt.Fprintln(os.Stderr, "Undefined variable", tok.val)
				os.Exit(1)
			}
			return lv
		}
	}

	return num()
}

// callArgs = (assign ("," assign)*)? ")"
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
	ret := &intLit{val: tokens[0].num}
	advance()
	return ret
}

func newAddBinary(lhs, rhs expression) expression {
	addType(lhs)
	addType(rhs)
	// num + num
	if lhs.getType().isInteger() && rhs.getType().isInteger() {
		return &binaryNode{op: "+", lhs: lhs, rhs: rhs}
	}

	// canonicalize num + ptr to ptr + num
	if lhs.getType().isInteger() && rhs.getType().isPointer() {
		tmp := lhs
		lhs = rhs
		rhs = tmp
	}

	// ptr + num
	if lhs.getType().isPointer() && rhs.getType().isInteger() {
		rhs = &binaryNode{op: "*", lhs: rhs, rhs: &intLit{val: 8}}
		return &binaryNode{op: "+", lhs: lhs, rhs: rhs}
	}

	panic("invalid operands")
}

func newSubBinary(lhs, rhs expression) expression {
	addType(lhs)
	addType(rhs)
	// num - num
	if lhs.getType().isInteger() && rhs.getType().isInteger() {
		return &binaryNode{op: "-", lhs: lhs, rhs: rhs}
	}

	// ptr - num
	if lhs.getType().isPointer() && rhs.getType().isInteger() {
		rhs = &binaryNode{op: "*", lhs: rhs, rhs: &intLit{val: 8}}
		return &binaryNode{op: "-", lhs: lhs, rhs: rhs}
	}

	// ptr - ptr, which returns how many elements are between the two.
	if lhs.getType().isPointer() && rhs.getType().isPointer() {
		lhs = &binaryNode{op: "-", lhs: lhs, rhs: rhs, ty: newLiteralType("int")}
		return &binaryNode{op: "/", lhs: lhs, rhs: &intLit{val: 8}}
	}

	panic("invalid operands")
}
