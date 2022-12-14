package generator_test

import (
	"fmt"
	"testing"

	"github.com/goropikari/tlex/automata"
	"github.com/goropikari/tlex/compiler/generator"
	"github.com/stretchr/testify/require"
)

func TestDFA_Accept(t *testing.T) {
	t.Parallel()

	letter := "(a|b|c|d|e|f|g|h|i|j|k|l|m|n|o|p|q|r|s|t|u|v|w|x|y|z)"
	digit := "(0|1|2|3|4|5|6|7|8|9)"
	digits := digit + digit + "*"
	id := fmt.Sprintf("%v(%v|%v)*", letter, letter, digit)

	regexs := []string{
		digits,                       // regexID: 1
		"if|then|begin|end|func|あいう", // regexID: 2
		id,                           // regexID: 3
		"\\+|\\-|\\*|/",              // regexID: 4
		"( |\n|\t|\r)",               // regexID: 5
		"\\.",                        // regexID: 6
		".",                          // regexID: 7
	}

	tests := []struct {
		name   string
		regexs []string
		given  string
		// expected
		accept  bool
		regexID automata.RegexID
	}{
		{
			name:    "digits",
			regexs:  regexs,
			given:   "123",
			accept:  true,
			regexID: 1,
		},
		{
			name:    "keyword: if",
			regexs:  regexs,
			given:   "if",
			accept:  true,
			regexID: 2,
		},
		{
			name:    "keyword: then",
			regexs:  regexs,
			given:   "then",
			accept:  true,
			regexID: 2,
		},
		{
			name:    "keyword: begin",
			regexs:  regexs,
			given:   "begin",
			accept:  true,
			regexID: 2,
		},
		{
			name:    "keyword: end",
			regexs:  regexs,
			given:   "end",
			accept:  true,
			regexID: 2,
		},
		{
			name:    "keyword: func",
			regexs:  regexs,
			given:   "func",
			accept:  true,
			regexID: 2,
		},
		{
			name:    "keyword: unicode",
			regexs:  regexs,
			given:   "あいう",
			accept:  true,
			regexID: 2,
		},
		{
			name:    "identifier",
			regexs:  regexs,
			given:   "ifhoge",
			accept:  true,
			regexID: 3,
		},
		{
			name:    "identifier: hoge",
			regexs:  regexs,
			given:   "hoge",
			accept:  true,
			regexID: 3,
		},
		{
			name:    "operator: +",
			regexs:  regexs,
			given:   "+",
			accept:  true,
			regexID: 4,
		},
		{
			name:    "operator: -",
			regexs:  regexs,
			given:   "-",
			accept:  true,
			regexID: 4,
		},
		{
			name:    "operator: *",
			regexs:  regexs,
			given:   "*",
			accept:  true,
			regexID: 4,
		},
		{
			name:    "operator: /",
			regexs:  regexs,
			given:   "/",
			accept:  true,
			regexID: 4,
		},
		{
			name:    "whitespace: space",
			regexs:  regexs,
			given:   " ",
			accept:  true,
			regexID: 5,
		},
		{
			name:    "dot: .",
			regexs:  regexs,
			given:   ".",
			accept:  true,
			regexID: 6,
		},
		{
			name:    "other: %",
			regexs:  regexs,
			given:   "%",
			accept:  true,
			regexID: 7,
		},
		{
			name:    "dot2: ..",
			regexs:  regexs,
			given:   "..",
			accept:  false,
			regexID: 0,
		},
		{
			name:    "identifier: start with digit",
			regexs:  regexs,
			given:   "0hoge",
			accept:  false,
			regexID: 0,
		},
		{
			name:    "arbitrary character",
			regexs:  []string{".*"},
			given:   "abc",
			accept:  true,
			regexID: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dfa := generator.LexerNFA(tt.regexs).ToImdNFA().ToDFA().LexerMinimize()

			regexID, accept := dfa.Accept(tt.given)

			require.Equal(t, tt.accept, accept)
			require.Equal(t, tt.regexID, regexID)
		})
	}
}

func TestDot(t *testing.T) {
	// _, _ = generator.LexerNFA([]string{"a", "abb", "a*bb*"}).
	// 	ToImdNFA().
	// 	ToDFA().
	// 	LexerMinimize().
	// 	ToDot()
	generator.LexerNFA(
		[]string{
			"if|for|while|func|return",
			"[a-zA-Z][a-zA-Z0-9]*",
			"[1-9][0-9]*",
			"[ \t\n\r]*",
			"\\(",
			"\\)",
			"\\{",
			"\\}",
			"\\+|\\-|\\*|\\/|:=|==|!=",
			".",
		},
	).
		ToImdNFA().
		ToDFA().
		LexerMinimize().
		ToDot()
}
