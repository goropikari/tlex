package main

import (
	"errors"
	"testing"
)

func TestLexer(t *testing.T) {
	program := `
func foo000() int {
    x := 1 * 10 + 123 - 1000 / 5432

    return x
}
`

	lex := New(program)
	ns := make([]int, 0)
	strs := make([]string, 0)

	expected := []struct {
		typ  int
		text string
	}{
		{Keyword, "func"},
		{Identifier, "foo000"},
		{LParen, "("},
		{RParen, ")"},
		{Type, "int"},
		{LBracket, "{"},
		{Identifier, "x"},
		{Operator, ":="},
		{Digit, "1"},
		{Operator, "*"},
		{Digit, "10"},
		{Operator, "+"},
		{Digit, "123"},
		{Operator, "-"},
		{Digit, "1000"},
		{Operator, "/"},
		{Digit, "5432"},
		{Keyword, "return"},
		{Identifier, "x"},
		{RBracket, "}"},
	}
	for {
		n, err := lex.Next()
		if err != nil {
			if errors.Is(err, EOF) {
				break
			}
			return
		}
		ns = append(ns, n)
		strs = append(strs, lex.YYText)
	}

	if len(expected) != len(ns) {
		t.Error("type is different")
	}
	for i, v := range expected {
		if v.typ != ns[i] {
			t.Error("type is different")
		}
	}
	if len(expected) != len(strs) {
		t.Error("token is different")
	}
	for i, v := range expected {
		if v.text != strs[i] {
			t.Error("token is different")
		}
	}
}
