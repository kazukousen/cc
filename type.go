package main

type typeKind int

const (
	typeKindInt typeKind = iota
	typeKindBool
	typeKindChar
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
	return ty.kind == typeKindInt || ty.kind == typeKindChar
}

func (ty *typ) hasBase() bool {
	return ty.base != nil
}

func newLiteralType(s string) *typ {
	typeKindMap := map[string]typeKind{
		"int":  typeKindInt,
		"bool": typeKindBool,
		"char": typeKindChar,
	}
	typeKindSize := map[string]int{
		"int":  8,
		"bool": 1,
		"char": 1,
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

func assignLVarOffsets() int {
	offset := 0
	for i := len(locals) - 1; i >= 0; i-- {
		lv := locals[i]
		offset += lv.ty.size
		lv.offset = offset
	}
	return offset
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
		if ty.hasBase() {
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
		n.setType(n.lhs.getType())
		return
	}
}
