package main

import (
	"errors"
	"log"
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
		{Identifier, "int"},
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
		strs = append(strs, YYtext)
	}

	if len(expected) != len(ns) {
		log.Fatal(errors.New("type is different"))
	}
	for i, v := range expected {
		if v.typ != ns[i] {
			log.Fatal(errors.New("type is different"))
		}
	}
	if len(expected) != len(strs) {
		log.Fatal(errors.New("token is different"))
	}
	for i, v := range expected {
		if v.text != strs[i] {
			log.Fatal(errors.New("token is different"))
		}
	}
}
