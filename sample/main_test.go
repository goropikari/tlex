package main

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

func TestLexer(t *testing.T) {
	program := `
func foo000あいう() int {
    x := 1 * 10 + 123 - 1000 / 5432

    return x
}
`

	lex := New(bytes.NewReader([]byte(program)))
	ns := make([]int, 0)
	strs := make([]string, 0)

	expected := []struct {
		typ  int
		text string
	}{
		{Keyword, "func"},
		{Identifier, "foo000"},
		{Hiragana, "あいう"},
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
			if errors.Is(err, io.EOF) {
				break
			}
			t.Fatal(err)
		}
		ns = append(ns, n)
		strs = append(strs, lex.YYText)
	}

	if len(expected) != len(ns) {
		t.Error("the number of recognize token is different from expected length.")
	}
	for i, v := range expected {
		if v.typ != ns[i] {
			t.Errorf("type is different: expected %v but %v", v.typ, ns[i])
		}
		if v.text != strs[i] {
			t.Errorf("token is different: expected %v but %v", v.text, strs[i])
		}
	}
	// for i, v := range expected {
	// 	if v.text != strs[i] {
	// 		t.Errorf("token is different: expected %v but %v", v.text, strs[i])
	// 	}
	// }
}
