package regexp_test

import (
	"testing"

	"github.com/goropikari/tlex/compiler/regexp"
	"github.com/stretchr/testify/require"
)

func TestLexer_Scan(t *testing.T) {
	tests := []struct {
		name     string
		regex    string
		expected []regexp.Token
	}{
		{
			name:  "lexer test",
			regex: "a(b|c*)de\t\n[a-z][^A-Z]\\+\\-\\*/",
			expected: []regexp.Token{
				regexp.NewToken(regexp.SymbolTokenType, 'a'),
				regexp.NewToken(regexp.LParenTokenType, '('),
				regexp.NewToken(regexp.SymbolTokenType, 'b'),
				regexp.NewToken(regexp.BarTokenType, '|'),
				regexp.NewToken(regexp.SymbolTokenType, 'c'),
				regexp.NewToken(regexp.StarTokenType, '*'),
				regexp.NewToken(regexp.RParenTokenType, ')'),
				regexp.NewToken(regexp.SymbolTokenType, 'd'),
				regexp.NewToken(regexp.SymbolTokenType, 'e'),
				regexp.NewToken(regexp.SymbolTokenType, '\t'),
				regexp.NewToken(regexp.SymbolTokenType, '\n'),
				regexp.NewToken(regexp.LSqBracketTokenType, '['),
				regexp.NewToken(regexp.SymbolTokenType, 'a'),
				regexp.NewToken(regexp.MinusTokenType, '-'),
				regexp.NewToken(regexp.SymbolTokenType, 'z'),
				regexp.NewToken(regexp.RSqBracketTokenType, ']'),
				regexp.NewToken(regexp.LSqBracketTokenType, '['),
				regexp.NewToken(regexp.NegationTokenType, '^'),
				regexp.NewToken(regexp.SymbolTokenType, 'A'),
				regexp.NewToken(regexp.MinusTokenType, '-'),
				regexp.NewToken(regexp.SymbolTokenType, 'Z'),
				regexp.NewToken(regexp.RSqBracketTokenType, ']'),
				regexp.NewToken(regexp.SymbolTokenType, '+'),
				regexp.NewToken(regexp.SymbolTokenType, '-'),
				regexp.NewToken(regexp.SymbolTokenType, '*'),
				regexp.NewToken(regexp.SymbolTokenType, '/'),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := regexp.NewLexer(tt.regex)
			toks := lexer.Scan()

			require.Equal(t, tt.expected, toks)
		})
	}

}
