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

type program struct {
	funcs []*function
}

var globals = make(map[string]*obj)

type function struct {
	name      string
	params    []*typ
	body      statement
	locals    []*obj
	stackSize int
}

type obj struct {
	ty   *typ
	name string

	isLocal bool

	// local variable
	offset int

	// global variable
	initData []byte
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

type funcCallNode struct {
	ty   *typ
	name string
	args []expression
}

type memberNode struct {
	unaryNode
	member *member
}

func (n *memberNode) isExpr() {}

func (n *memberNode) getType() *typ { return n.ty }

func (n *memberNode) setType(ty *typ) { n.ty = ty }

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

type scope struct {
	vars []*obj
}

var scopes []*scope

func enterScope() {
	scopes = append(scopes, &scope{vars: []*obj{}})
}

func leaveScope() {
	scopes = scopes[:len(scopes)-1]
}

func pushScope(v *obj) {
	sc := scopes[len(scopes)-1]
	sc.vars = append(sc.vars, v)
}

func (f *function) getLocal(name string) *obj {
	for i := range f.locals {
		lv := f.locals[i]
		if lv.name == name {
			return lv
		}
	}

	return nil
}

func findLocalInScope(name string) *obj {
	for i := len(scopes) - 1; i >= 0; i-- {
		sc := scopes[i]
		for vi := range sc.vars {
			lv := sc.vars[vi]
			if lv.name == name {
				return lv
			}
		}
	}

	if gv, ok := globals[name]; ok {
		return gv
	}

	return nil
}

func newNodeLocal(ty *typ) *obj {

	lv := &obj{
		name:    ty.name,
		ty:      ty,
		isLocal: true,
	}

	pushScope(lv)
	locals = append(locals, lv)

	return lv
}

func newGlobalVariable(ty *typ) *obj {
	if gv, ok := globals[ty.name]; ok {
		return gv
	}
	gv := &obj{
		ty:   ty,
		name: ty.name,
	}
	globals[ty.name] = gv
	return gv
}

var uniqueID = 0

func newUniqueName() string {
	s := fmt.Sprintf(".L..%d", uniqueID)
	uniqueID++
	return s
}

func newStringLiteral(s string) *obj {
	ty := newLiteralType("char")
	ty = arrayOf(ty, len(s))
	ty.name = newUniqueName()
	gv := newGlobalVariable(ty)
	gv.initData = []byte(s)
	return gv
}

// program = decl*
// decl = declspec declarator ("{" funcDecl | varDecl)
func parse() *program {
	prog := &program{
		funcs: []*function{},
	}
	for len(tokens) > 0 {
		ty := declSpec()
		ty = declarator(ty)
		if consume("{") {
			prog.funcs = append(prog.funcs, funcDecl(ty))
			continue
		}
		varDecl(ty)
	}
	return prog
}

// varDecl = ("," declarator)* ";"
func varDecl(ty *typ) {
	_ = newGlobalVariable(ty)
	for consume(",") {
		ty = declarator(ty)
		_ = newGlobalVariable(ty)
	}
	expect(";")
}

// funcDecl = compoundStmt
func funcDecl(ty *typ) *function {

	locals = []*obj{}
	scopes = []*scope{}

	f := &function{
		name:   ty.name,
		params: ty.params,
	}

	enterScope()

	for _, p := range f.params {
		newNodeLocal(p)
	}

	f.body = compoundStmt()
	f.stackSize = assignLVarOffsets()
	addType(f.body)

	leaveScope()

	f.locals = locals

	return f
}

// declspec = "int" | "char" | struct-decl
func declSpec() *typ {
	tok := consumeToken(tokenKindType)
	if tok == nil {
		_, _ = fmt.Fprintln(os.Stderr, "Expect an type:", tokens[0].val)
		os.Exit(1)
	}

	if tok.val == "struct" {
		return structDecl()
	}

	return newLiteralType(tok.val)
}

// struct-decl = "{" (declspec declarator ("," declarator)* ";")* "}"
func structDecl() *typ {

	expect("{")

	var members []*member
	for !consume("}") {
		baseTy := declSpec()

		for i := 0; !consume(";"); i++ {
			if i > 0 {
				expect(",")
			}
			ty := declarator(baseTy)
			members = append(members, &member{
				ty:   ty,
				name: ty.name,
			})
		}
	}

	return newStructType(members)
}

// declarator = "*"* ident type-suffix
func declarator(baseTy *typ) *typ {

	ty := new(typ)
	*ty = *baseTy

	for consume("*") {
		ty = pointerTo(ty)
	}

	tok := consumeToken(tokenKindIdent)
	if tok == nil {
		_, _ = fmt.Fprintln(os.Stderr, "Expect an identifier in declarator:", tokens[0].val)
		os.Exit(1)
	}

	ty.name = tok.val

	ty = typeSuffix(ty)

	return ty
}

// type-suffix = "(" func-params | "[" num "]" type-suffix | ε
func typeSuffix(ty *typ) *typ {
	if consume("(") {
		return funcParams(ty)
	}

	if consume("[") {
		length := num().val
		expect("]")
		ty = typeSuffix(ty)
		ty = arrayOf(ty, length)
		return ty
	}

	return ty
}

// func-params = (param ("," param)*)? ")"
// param = declspec declarator
func funcParams(ty *typ) *typ {
	var params []*typ
	for i := 0; !consume(")"); i++ {
		if i > 0 {
			expect(",")
		}
		p := declSpec()
		p = declarator(p)
		params = append(params, p)
	}
	ty.params = params
	return ty
}

// compoundStmt = (declaration | stmt)* "}"
func compoundStmt() statement {

	enterScope()

	ret := &blockStmtNode{code: []statement{}}
	for !consume("}") {
		if equalToken(tokenKindType) {
			ret.code = append(ret.code, declaration()...)
		} else {
			ret.code = append(ret.code, stmt())
		}
	}

	leaveScope()

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
		ty = declarator(ty)
		lv := newNodeLocal(ty)
		if consume("=") {
			n := &assignNode{op: "=", lhs: lv, rhs: expr()}
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

// unary = ("-" | "+" | "&" | "*") unary | postfix
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
		return postfix()
	}
}

// postfix = primary ("[" expr "]" | "." ident)*
func postfix() expression {
	ret := primary()

	for {
		if consume("[") {
			ret = &derefNode{child: newAddBinary(ret, expr())}
			expect("]")
			continue
		}

		if consume(".") {
			ret = structRef(ret)
			continue
		}

		return ret
	}
}

func structRef(n expression) expression {
	addType(n)
	ty := n.getType()
	if ty.kind != typeKindStruct {
		panic("expected struct type")
	}

	tok := consumeToken(tokenKindIdent)
	if tok == nil {
		_, _ = fmt.Fprintln(os.Stderr, "ident expected", tok.val)
		os.Exit(1)
	}

	var mem *member
	for i := range ty.members {
		m := ty.members[i]
		if m.name == tok.val {
			mem = m
		}
	}

	return &memberNode{unaryNode: unaryNode{child: n}, member: mem}
}

// primary = "(" expr ")" | "sizeof" unary | ident ("(" callArgs)? | num
func primary() expression {
	if consume("(") {
		ret := expr()
		expect(")")
		return ret
	}

	if consume("sizeof") {
		n := unary()
		addType(n)
		return &intLit{val: n.getType().size}
	}

	if tok := consumeToken(tokenKindIdent); tok != nil {
		if consume("(") {
			return &funcCallNode{name: tok.val, args: callArgs(), ty: newLiteralType("int")}
		} else {
			lv := findLocalInScope(tok.val)
			if lv == nil {
				_, _ = fmt.Fprintln(os.Stderr, "Undefined variable", tok.val)
				os.Exit(1)
			}

			return lv
		}
	}

	if tok := consumeToken(tokenKindStringLiteral); tok != nil {
		return newStringLiteral(tok.str)
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

func num() *intLit {
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
	if lhs.getType().isInteger() && rhs.getType().hasBase() {
		tmp := lhs
		lhs = rhs
		rhs = tmp
	}

	// ptr + num
	if lhs.getType().hasBase() && rhs.getType().isInteger() {
		rhs = &binaryNode{op: "*", lhs: rhs, rhs: &intLit{val: lhs.getType().base.size}}
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
	if lhs.getType().hasBase() && rhs.getType().isInteger() {
		rhs = &binaryNode{op: "*", lhs: rhs, rhs: &intLit{val: lhs.getType().base.size}}
		return &binaryNode{op: "-", lhs: lhs, rhs: rhs}
	}

	// ptr - ptr, which returns how many elements are between the two.
	if lhs.getType().hasBase() && rhs.getType().hasBase() {
		n := &binaryNode{op: "-", lhs: lhs, rhs: rhs, ty: newLiteralType("int")}
		return &binaryNode{op: "/", lhs: n, rhs: &intLit{val: lhs.getType().base.size}}
	}

	panic("invalid operands")
}
