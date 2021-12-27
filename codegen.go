package main

import (
	"fmt"
	"os"
)

var label = 0

func gen(n *node) {
	switch n.kind {
	case nodeKindReturn:
		gen(n.lhs)
		fmt.Printf("	pop rax\n")
		fmt.Printf("	mov rsp, rbp\n")
		fmt.Printf("	pop rbp\n")
		fmt.Printf("	ret\n")
		return
	case nodeKindNum:
		fmt.Printf("	push %d\n", n.num)
		return
	case nodeKindLocal:
		genLocal(n)
		fmt.Printf("	pop rax\n")
		fmt.Printf("	mov rax, [rax]\n")
		fmt.Printf("	push rax\n")
		return
	case nodeKindAssign:
		genLocal(n.lhs)
		gen(n.rhs)
		fmt.Printf("	pop rdi\n")
		fmt.Printf("	pop rax\n")
		fmt.Printf("	mov [rax], rdi\n")
		fmt.Printf("	push rdi\n")
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
	case nodeKindWhile:
		fmt.Printf(".Lbegin%d:\n", label)
		gen(n.cond)
		fmt.Printf("	pop rax\n")
		fmt.Printf("	cmp rax, 0\n")
		fmt.Printf("	je .Lend%d\n", label)
		gen(n.then)
		fmt.Printf("	jmp .Lbegin%d\n", label)
		fmt.Printf(".Lend%d:\n", label)
	case nodeKindBlock:
		for _, s := range n.code {
			gen(s)
		}
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

func genLocal(n *node) {
	if n.kind != nodeKindLocal {
		_, _ = fmt.Fprintln(os.Stderr, "left side of assignment is not variable")
		os.Exit(1)
	}
	fmt.Printf("	mov rax, rbp\n")
	fmt.Printf("	sub rax, %d\n", n.offset)
	fmt.Printf("	push rax\n")
}
