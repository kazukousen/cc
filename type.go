package main

import "fmt"

type typeKind int

const (
	typeKindInt typeKind = iota
	typeKindBool
	typeKindArray
	typeKindPtr
)

type typ struct {
	kind   typeKind
	base   *typ
	name   string
	params []*typ
	size   int
	length int
}

func (ty *typ) isInteger() bool {
	return ty.kind == typeKindInt
}

func (ty *typ) hasBase() bool {
	return ty.base != nil
}

func newLiteralType(s string) *typ {
	typeKindMap := map[string]typeKind{
		"int":  typeKindInt,
		"bool": typeKindBool,
	}
	typeKindSize := map[string]int{
		"int":  8,
		"bool": 1,
	}
	return &typ{
		kind: typeKindMap[s],
		size: typeKindSize[s],
	}
}

func pointerTo(base *typ) *typ {
	return &typ{
		kind: typeKindPtr,
		base: base,
		name: base.name,
		size: 8,
	}
}

func arrayOf(base *typ, length int) *typ {
	return &typ{
		kind:   typeKindArray,
		base:   base,
		name:   base.name,
		size:   base.size * length,
		length: length,
	}
}

func calcStackSize(vs []*obj) int {
	var ret int
	for _, v := range vs {
		ret += v.ty.size
	}
	return ret
}

func addType(n interface{}) {
	if n == nil {
		return
	}
	if n, ok := n.(interface{ getType() *typ }); ok && n.getType() != nil {
		return
	}

	switch n := n.(type) {
	case *intLit:
		n.setType(newLiteralType("int"))
		return
	case *addrNode:
		addType(n.child)
		ct := n.child.getType()
		if ct.kind == typeKindArray {
			ct = n.child.getType().base
		}
		ty := pointerTo(ct)
		n.setType(ty)
		return
	case *derefNode:
		addType(n.child)
		ty := n.child.getType()
		if !ty.hasBase() {
			panic(fmt.Sprintf("invalid pointer dereference: %v", ty))
		}
		n.setType(ty.base)
		return
	case *binaryNode:
		addType(n.lhs)
		addType(n.rhs)
		switch n.op {
		case "+", "-", "*", "/":
			n.setType(n.lhs.getType())
		case "==", "!=", "<", "<=":
			n.setType(newLiteralType("bool"))
		}
		return
	case *obj:
		return
	case *funcCallNode:
		for _, c := range n.args {
			addType(c)
		}
		return
	case *exprStmtNode:
		addType(n.child)
		return
	case *returnStmtNode:
		addType(n.child)
		return
	case *blockStmtNode:
		for _, c := range n.code {
			addType(c)
		}
		return
	case *ifStmtNode:
		addType(n.cond)
		addType(n.then)
		addType(n.els)
		return
	case *forStmtNode:
		addType(n.ini)
		addType(n.cond)
		addType(n.step)
		addType(n.then)
		return
	case *assignNode:
		addType(n.lhs)
		addType(n.rhs)
		n.lhs.setType(n.rhs.getType())
		n.setType(n.lhs.getType())
		return
	}
}
