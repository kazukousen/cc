package main

import (
	"fmt"
	"os"
	"strings"
)

var userIn string
var in string
var currentFilename string

type token struct {
	kind tokenKind
	val  string
	num  int
	str  string
}

type tokenKind int

const (
	tokenKindReserved tokenKind = iota
	tokenKindNumberLiteral
	tokenKindStringLiteral
	tokenKindIdent
	tokenKindType
)

func tokenizeFile(filename string) {

	b, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	userIn = string(b)
	in = string(b)
	currentFilename = filename

	tokenize()
}

func tokenize() {
	for len(in) > 0 {

		if in[0] == '"' {
			tokens = append(tokens, toString())
			continue
		}

		if in[0] == ' ' || in[0] == '\n' {
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
			if len(in) > 1 && (in[0:2] == "<=" || in[0:2] == ">=" || in[0:2] == "==" || in[0:2] == "!=" || in[0:2] == "->") {
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

// foo.c:10:5: x + y = 1;
//               ^ error message here
func errorAt(msg string) {

	pos := len(userIn) - len(in)

	start := pos
	for i := pos; i >= 0 && userIn[i] != '\n'; i-- {
		start = i
	}

	end := pos
	for i := 0; i < len(in) && in[i] != '\n'; i++ {
		end++
	}

	lineNo := 1
	for _, c := range userIn[:pos] {
		if c == '\n' {
			lineNo++
		}
	}

	pos = pos - start
	indent, _ := fmt.Fprintf(os.Stderr, "%s:%d:%d: ", currentFilename, lineNo, pos+1)
	_, _ = fmt.Fprintf(os.Stderr, "%s\n", userIn[start:end])
	_, _ = fmt.Fprintln(os.Stderr, strings.Repeat(" ", indent+pos-1), "^")
	_, _ = fmt.Fprintln(os.Stderr, strings.Repeat(" ", indent+pos-1), msg)
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
