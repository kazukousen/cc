package main

import (
	"fmt"
	"os"
)

var label = 0
var funcName string

var argRegisters64 = []string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}
var argRegisters8 = []string{"dil", "sil", "dl", "cl", "r8b", "r9b"}

func codegen(prog *program) {
	fmt.Printf(".intel_syntax noprefix\n")
	emitData()
	emitText(prog.funcs)
}

func emitData() {
	for _, gv := range globals {
		fmt.Printf("	.data\n")
		fmt.Printf("	.globl %s\n", gv.name)
		fmt.Printf("%s:\n", gv.name)
		if gv.initData != nil {
			for _, b := range gv.initData {
				fmt.Printf("	.byte %d\n", b)
			}
		} else {
			fmt.Printf("	.zero %d\n", gv.ty.size)
		}
	}
}

func emitText(funcs []*function) {
	for _, f := range funcs {
		funcName = f.name
		fmt.Printf(`	.globl %[1]s
	.text
%[1]s:
	push rbp
	mov rbp, rsp
	sub rsp, %[2]d
`, funcName, f.stackSize)

		for i, p := range f.params {
			if p.size == 1 {
				fmt.Printf("	mov [rbp-%d], %s\n", f.getLocal(p.name).offset, argRegisters8[i])
			} else {
				fmt.Printf("	mov [rbp-%d], %s\n", f.getLocal(p.name).offset, argRegisters64[i])
			}
		}

		gen(f.body)

		fmt.Printf(`.Lreturn.%s:
	mov rsp, rbp
	pop rbp
	ret
`, funcName)
	}
}

func gen(n interface{}) {
	switch n := n.(type) {
	case *returnStmtNode:
		gen(n.child)
		fmt.Printf("	pop rax\n")
		fmt.Printf("	jmp .Lreturn.%s\n", funcName)
		return
	case *intLit:
		fmt.Printf("	push %d\n", n.val)
		return
	case *obj:
		genAddr(n)
		load(n.getType())
		return
	case *memberNode:
		genAddr(n)
		load(n.getType())
		return
	case *assignNode:
		genAddr(n.lhs)
		gen(n.rhs)
		store(n.getType())
		return
	case *ifStmtNode:
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
	case *forStmtNode:
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
		return
	case *blockStmtNode:
		for _, s := range n.code {
			gen(s)
		}
		return
	case *exprStmtNode:
		gen(n.child)
		return
	case *funcCallNode:

		for _, arg := range n.args {
			gen(arg)
		}

		for i := len(n.args) - 1; i >= 0; i-- {
			fmt.Printf("	pop %s\n", argRegisters64[i])
		}

		fmt.Printf("	call %s\n", n.name)
		fmt.Printf("	push rax\n")
		return
	case *addrNode:
		genAddr(n.child)
		return
	case *derefNode:
		gen(n.child)
		load(n.getType())
		return
	}

	b := n.(*binaryNode)

	if b.lhs != nil {
		gen(b.lhs)
	}
	if b.rhs != nil {
		gen(b.rhs)
	}

	fmt.Printf("	pop rdi\n")
	fmt.Printf("	pop rax\n")

	switch b.op {
	case "+":
		fmt.Printf("	add rax, rdi\n")
	case "-":
		fmt.Printf("	sub rax, rdi\n")
	case "*":
		fmt.Printf("	imul rax, rdi\n")
	case "/":
		fmt.Printf("	cqo\n")
		fmt.Printf("	idiv rdi\n")
	case "<":
		fmt.Printf("	cmp rax, rdi\n")
		fmt.Printf("	setl al\n")
		fmt.Printf("	movzb rax, al\n")
	case "<=":
		fmt.Printf("	cmp rax, rdi\n")
		fmt.Printf("	setle al\n")
		fmt.Printf("	movzb rax, al\n")
	case "==":
		fmt.Printf("	cmp rax, rdi\n")
		fmt.Printf("	sete al\n")
		fmt.Printf("	movzb rax, al\n")
	case "!=":
		fmt.Printf("	cmp rax, rdi\n")
		fmt.Printf("	setne al\n")
		fmt.Printf("	movzb rax, al\n")
	}

	fmt.Printf("	push rax\n")
}

func genAddr(n expression) {
	switch n := n.(type) {
	case *obj:
		if n.isLocal {
			fmt.Printf("	lea rax, [rbp-%d]\n", n.offset)
			fmt.Printf("	push rax\n")
		} else {
			fmt.Printf("	push offset %s\n", n.name)
		}
	case *derefNode:
		gen(n.child)
	case *memberNode:
		genAddr(n.child)
		fmt.Printf("	pop rax\n")
		fmt.Printf("	add rax, %d\n", n.member.offset)
		fmt.Printf("	push rax\n")
	default:
		_, _ = fmt.Fprintln(os.Stderr, "Not an identifier")
		os.Exit(1)
	}
}

func load(ty *typ) {
	if ty.kind == typeKindArray {
		return
	}
	if ty.size == 1 {
		fmt.Printf("	pop rax\n")
		fmt.Printf("	movzx rax, byte ptr [rax]\n")
		fmt.Printf("	push rax\n")
	} else {
		fmt.Printf("	pop rax\n")
		fmt.Printf("	mov rax, [rax]\n")
		fmt.Printf("	push rax\n")
	}
}

func store(ty *typ) {
	fmt.Printf("	pop rdi\n")
	fmt.Printf("	pop rax\n")
	if ty.size == 1 {
		fmt.Printf("	mov [rax], dil\n")
	} else {
		fmt.Printf("	mov [rax], rdi\n")
	}
	fmt.Printf("	push rdi\n") // e.g. a=b=3
}
