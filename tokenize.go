package main

import (
	"fmt"
	"os"
	"strings"
)

type tokenKind int

const (
	tokenKindReserved tokenKind = iota
	tokenKindNumber
	tokenKindIdent
)

type token struct {
	kind tokenKind
	val  string
	num  int
}

func tokenize() {
	for len(in) > 0 {

		if in[0] == ' ' {
			in = in[1:]
			continue
		}

		if in[0] >= 'a' && in[0] <= 'z' {
			tokens = append(tokens, &token{kind: tokenKindIdent, val: string(in[0])})
			in = in[1:]
			continue
		}

		if strings.Contains("+-*/()<>=!;", string(in[0])) {
			if len(in) > 1 && (in[0:2] == "<=" || in[0:2] == ">=" || in[0:2] == "==" || in[0:2] == "!=") {
				tokens = append(tokens, &token{kind: tokenKindReserved, val: in[0:2]})
				in = in[2:]
			} else {
				tokens = append(tokens, &token{kind: tokenKindReserved, val: string(in[0])})
				in = in[1:]
			}
			continue
		}

		if in[0] >= '0' && in[0] <= '9' {
			n := toInt()
			tokens = append(tokens, &token{kind: tokenKindNumber, num: n})
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

func toInt() int {
	var ret int
	for len(in) > 0 && in[0] >= '0' && in[0] <= '9' {
		ret = ret*10 + int(in[0]-'0')
		in = in[1:]
	}
	return ret
}
