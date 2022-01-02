package main

import (
	"fmt"
	"os"
	"strings"
)

type tokenKind int

const (
	tokenKindReserved tokenKind = iota
	tokenKindNumberLiteral
	tokenKindStringLiteral
	tokenKindIdent
	tokenKindType
)

type token struct {
	kind tokenKind
	val  string
	num  int
	str  string
}

func isAlpha() bool {
	return (in[0] >= 'a' && in[0] <= 'z') || (in[0] >= 'A' && in[0] <= 'Z') || in[0] == '_'
}

func isDigit() bool {
	return in[0] >= '0' && in[0] <= '9'
}

func identifierToken(val string) *token {
	for _, w := range []string{"return", "if", "else", "while", "for", "sizeof"} {
		if val == w {
			return &token{kind: tokenKindReserved, val: val}
		}
	}
	for _, w := range []string{"int", "char", "struct"} {
		if val == w {
			return &token{kind: tokenKindType, val: val}
		}
	}
	return &token{kind: tokenKindIdent, val: val}
}

func tokenize() {
	for len(in) > 0 {

		if in[0] == '"' {
			tokens = append(tokens, toString())
			continue
		}

		if in[0] == ' ' {
			in = in[1:]
			continue
		}

		if in[0] >= 'a' && in[0] <= 'z' {
			name := in[0:1]
			in = in[1:]
			for len(in) > 0 && (isAlpha() || isDigit()) {
				name += in[0:1]
				in = in[1:]
			}
			tokens = append(tokens, identifierToken(name))
			continue
		}

		if strings.Contains("+-*/()<>=!;{},&[].", string(in[0])) {
			if len(in) > 1 && (in[0:2] == "<=" || in[0:2] == ">=" || in[0:2] == "==" || in[0:2] == "!=") {
				tokens = append(tokens, &token{kind: tokenKindReserved, val: in[0:2]})
				in = in[2:]
			} else {
				tokens = append(tokens, &token{kind: tokenKindReserved, val: string(in[0])})
				in = in[1:]
			}
			continue
		}

		if isDigit() {
			n := toInt()
			tokens = append(tokens, &token{kind: tokenKindNumberLiteral, num: n})
			continue
		}

		errorAt("unexpected character: " + string(in[0]))
	}
}

func errorAt(msg string) {
	n := len(userIn) - len(in)
	_, _ = fmt.Fprintln(os.Stderr, userIn)
	_, _ = fmt.Fprintln(os.Stderr, strings.Repeat(" ", n), "^")
	_, _ = fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func toString() *token {
	start := len(userIn) - len(in)
	in = in[1:]
	for len(in) > 0 && in[0] != '"' {
		in = in[1:]
	}
	end := len(userIn) - len(in)
	str := userIn[start+1 : end]
	str += "\000"

	in = in[1:]

	return &token{kind: tokenKindStringLiteral, str: str}
}

func toInt() int {
	var ret int
	for len(in) > 0 && in[0] >= '0' && in[0] <= '9' {
		ret = ret*10 + int(in[0]-'0')
		in = in[1:]
	}
	return ret
}
