package main

import "fmt"

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
	case nodeKindLt:
		fmt.Printf("	cmp rax,rdi\n")
		fmt.Printf("	setl al\n")
		fmt.Printf("	movzb rax,al\n")
	case nodeKindLe:
		fmt.Printf("	cmp rax,rdi\n")
		fmt.Printf("	setle al\n")
		fmt.Printf("	movzb rax,al\n")
	case nodeKindEq:
		fmt.Printf("	cmp rax,rdi\n")
		fmt.Printf("	sete al\n")
		fmt.Printf("	movzb rax,al\n")
	case nodeKindNe:
		fmt.Printf("	cmp rax,rdi\n")
		fmt.Printf("	setne al\n")
		fmt.Printf("	movzb rax,al\n")
	}

	fmt.Printf("	push rax\n")
}
