package main

type typeKind int

const (
	typeKindInt typeKind = iota
	typeKindBool
	typeKindPtr
)

type typ struct {
	kind typeKind
	base *typ
	size int
}

func (ty *typ) isInteger() bool {
	return ty.kind == typeKindInt
}

func (ty *typ) isPointer() bool {
	return ty.kind == typeKindPtr
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

func pointerTo(ty *typ) *typ {
	return &typ{
		kind: typeKindPtr,
		base: ty,
		size: 8,
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
		ty := pointerTo(n.child.getType())
		n.setType(ty)
		return
	case *derefNode:
		addType(n.child)
		ty := n.child.getType()
		if ty.kind == typeKindPtr {
			n.setType(ty.base)
			return
		}
		n.setType(ty)
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
		return
	}
}
