package main

type typeKind int

const (
	typeKindInt typeKind = iota
	typeKindPtr
)

type typ struct {
	kind typeKind
	base *typ
	size int
}

func newLiteralType(s string) *typ {
	typeKindMap := map[string]typeKind{
		"int": typeKindInt,
	}
	typeKindSize := map[string]int{
		"int": 8,
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
