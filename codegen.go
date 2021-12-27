package main

import (
	"fmt"
	"os"
)

func gen(n *node) {
	switch n.kind {
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
