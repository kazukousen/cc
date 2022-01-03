package main

type typeKind int

const (
	typeKindInt typeKind = iota
	typeKindBool
	typeKindChar
	typeKindArray
	typeKindStruct
	typeKindPtr
)

type typ struct {
	kind  typeKind
	base  *typ
	name  string
	size  int
	align int

	// func
	params []*typ

	// array
	length int

	// struct
	members []*member
}

func (ty *typ) isInteger() bool {
	return ty.kind == typeKindInt || ty.kind == typeKindChar
}

func (ty *typ) hasBase() bool {
	return ty.base != nil
}

func newType(kind typeKind, size, align int) *typ {
	return &typ{kind: kind, size: size, align: align}
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
	typeKindAlign := map[string]int{
		"int":  8,
		"bool": 1,
		"char": 1,
	}
	return newType(typeKindMap[s], typeKindSize[s], typeKindAlign[s])
}

type member struct {
	ty     *typ
	name   string
	offset int
}

func newStructType(members []*member) *typ {

	align := 1
	offset := 0
	for i := range members {
		m := members[i]
		m.offset = alignTo(offset, m.ty.align)
		offset += m.ty.size

		if align < m.ty.align {
			align = m.ty.align
		}
	}

	ty := newType(typeKindStruct, alignTo(offset, align), align)
	ty.members = members
	return ty
}

func pointerTo(base *typ) *typ {
	ty := newType(typeKindPtr, 8, 8)
	ty.base = base
	ty.name = base.name
	return ty
}

func arrayOf(base *typ, length int) *typ {
	ty := newType(typeKindArray, base.size*length, base.align)
	ty.base = base
	ty.length = length
	ty.name = base.name
	return ty
}

func assignLVarOffsets() int {
	offset := 0
	for i := len(locals) - 1; i >= 0; i-- {
		lv := locals[i]
		offset += lv.ty.size
		offset = alignTo(offset, lv.ty.align)
		lv.offset = offset
	}
	return alignTo(offset, 16)
}

func alignTo(n, align int) int {
	return (n + align - 1) / align * align
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
	case *memberNode:
		addType(n.child)
		n.setType(n.member.ty)
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
