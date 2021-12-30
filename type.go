package main

type typeKind int

const (
	typeKindInt typeKind = iota
	typeKindPtr
)

type typ struct {
	kind typeKind
	base *typ
}

func newLiteralType(s string) *typ {
	typeKindMap := map[string]typeKind{
		"int": typeKindInt,
	}
	return &typ{
		kind: typeKindMap[s],
	}
}

func pointerTo(ty *typ) *typ {
	return &typ{
		kind: typeKindPtr,
		base: ty,
	}
}
