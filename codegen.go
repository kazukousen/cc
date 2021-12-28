package main

import (
	"fmt"
	"os"
)

var label = 0
var funcName string

var argRegisters = []string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}

func codegen(funcs []*function) {
	fmt.Printf(".intel_syntax noprefix\n")
	for _, f := range funcs {
		offset := 0
		for i := len(f.locals) - 1; i >= 0; i-- {
			v := f.locals[i]
			offset += 8
			v.offset = -offset
		}

		funcName = f.name

		fmt.Printf(`.globl %[1]s
%[1]s:
	push rbp
	mov rbp, rsp
	sub rsp, %[2]d
`, funcName, len(f.locals)*8)

		for i, arg := range f.args {
			fmt.Printf("	mov [rbp%d], %s\n", arg.offset, argRegisters[i])
		}

		gen(f.body)

		fmt.Printf(`.Lreturn.%s:
	mov rsp, rbp
	pop rbp
	ret
`, funcName)
	}
}

func gen(n *node) {
	switch n.kind {
	case nodeKindReturn:
		gen(n.lhs)
		fmt.Printf("	pop rax\n")
		fmt.Printf("	jmp .Lreturn.%s\n", funcName)
		return
	case nodeKindNum:
		fmt.Printf("	push %d\n", n.num)
		return
	case nodeKindVar:
		genAddr(n)
		load()
		return
	case nodeKindAssign:
		genAddr(n.lhs)
		gen(n.rhs)
		store()
		return
	case nodeKindIf:
		gen(n.cond)
		fmt.Printf("	pop rax\n")
		fmt.Printf("	cmp rax, 0\n")

		if n.els != nil {
			fmt.Printf("	je .Lelse%d\n", label)
			gen(n.then)
			fmt.Printf("	jmp .Lend%d\n", label)
			fmt.Printf(".Lelse%d:\n", label)
			gen(n.els)
			fmt.Printf(".Lend%d:\n", label)
		} else {
			fmt.Printf("	je .Lend%d\n", label)
			gen(n.then)
			fmt.Printf(".Lend%d:\n", label)
		}
		label++
		return
	case nodeKindFor:
		if n.ini != nil {
			gen(n.ini)
		}
		fmt.Printf(".Lbegin%d:\n", label)
		if n.cond != nil {
			gen(n.cond)
			fmt.Printf("	pop rax\n")
			fmt.Printf("	cmp rax, 0\n")
			fmt.Printf("	je .Lend%d\n", label)
		}
		gen(n.then)
		if n.step != nil {
			gen(n.step)
		}
		fmt.Printf("	jmp .Lbegin%d\n", label)
		fmt.Printf(".Lend%d:\n", label)
	case nodeKindBlock:
		for _, s := range n.code {
			gen(s)
		}
		return
	case nodeKindCall:

		for _, arg := range n.args {
			gen(arg)
		}

		for i := len(n.args) - 1; i >= 0; i-- {
			fmt.Printf("	pop %s\n", argRegisters[i])
		}

		fmt.Printf("	call %s\n", n.name)
		fmt.Printf("	push rax\n")
		return
	case nodeKindAddr:
		genAddr(n.lhs)
		return
	case nodeKindDeref:
		gen(n.lhs)
		load()
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
		fmt.Printf("	add rax, rdi\n")
	case nodeKindSub:
		fmt.Printf("	sub rax, rdi\n")
	case nodeKindMul:
		fmt.Printf("	imul rax, rdi\n")
	case nodeKindDiv:
		fmt.Printf("	cqo\n")
		fmt.Printf("	idiv rdi\n")
	case nodeKindLt:
		fmt.Printf("	cmp rax, rdi\n")
		fmt.Printf("	setl al\n")
		fmt.Printf("	movzb rax, al\n")
	case nodeKindLe:
		fmt.Printf("	cmp rax, rdi\n")
		fmt.Printf("	setle al\n")
		fmt.Printf("	movzb rax, al\n")
	case nodeKindEq:
		fmt.Printf("	cmp rax, rdi\n")
		fmt.Printf("	sete al\n")
		fmt.Printf("	movzb rax, al\n")
	case nodeKindNe:
		fmt.Printf("	cmp rax, rdi\n")
		fmt.Printf("	setne al\n")
		fmt.Printf("	movzb rax, al\n")
	}

	fmt.Printf("	push rax\n")
}

func genAddr(n *node) {
	switch n.kind {
	case nodeKindVar:
		fmt.Printf("	lea rax, [rbp%d]\n", n.variable.offset)
		fmt.Printf("	push rax\n")
	case nodeKindDeref:
		gen(n.lhs)
	default:
		_, _ = fmt.Fprintln(os.Stderr, "Not a variable")
		os.Exit(1)
	}
}

func load() {
	fmt.Printf("	pop rax\n")
	fmt.Printf("	mov rax, [rax]\n")
	fmt.Printf("	push rax\n")
}

func store() {
	fmt.Printf("	pop rdi\n")
	fmt.Printf("	pop rax\n")
	fmt.Printf("	mov [rax], rdi\n")
	fmt.Printf("	push rdi\n")
}
