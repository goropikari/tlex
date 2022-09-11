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
			regex: `a(b|c*)de\t\n[a-z][^A-Z]\+\-\*\/`,
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
		{
			name:  "lexer string literal",
			regex: `"hoge\"piyo[xyz]ab*"`,
			expected: []regexp.Token{
				regexp.NewToken(regexp.SymbolTokenType, 'h'),
				regexp.NewToken(regexp.SymbolTokenType, 'o'),
				regexp.NewToken(regexp.SymbolTokenType, 'g'),
				regexp.NewToken(regexp.SymbolTokenType, 'e'),
				regexp.NewToken(regexp.SymbolTokenType, '"'),
				regexp.NewToken(regexp.SymbolTokenType, 'p'),
				regexp.NewToken(regexp.SymbolTokenType, 'i'),
				regexp.NewToken(regexp.SymbolTokenType, 'y'),
				regexp.NewToken(regexp.SymbolTokenType, 'o'),
				regexp.NewToken(regexp.SymbolTokenType, '['),
				regexp.NewToken(regexp.SymbolTokenType, 'x'),
				regexp.NewToken(regexp.SymbolTokenType, 'y'),
				regexp.NewToken(regexp.SymbolTokenType, 'z'),
				regexp.NewToken(regexp.SymbolTokenType, ']'),
				regexp.NewToken(regexp.SymbolTokenType, 'a'),
				regexp.NewToken(regexp.SymbolTokenType, 'b'),
				regexp.NewToken(regexp.SymbolTokenType, '*'),
			},
		},
		{
			name:  "alternation",
			regex: `if|for|while`,
			expected: []regexp.Token{
				regexp.NewToken(regexp.SymbolTokenType, 'i'),
				regexp.NewToken(regexp.SymbolTokenType, 'f'),
				regexp.NewToken(regexp.BarTokenType, '|'),
				regexp.NewToken(regexp.SymbolTokenType, 'f'),
				regexp.NewToken(regexp.SymbolTokenType, 'o'),
				regexp.NewToken(regexp.SymbolTokenType, 'r'),
				regexp.NewToken(regexp.BarTokenType, '|'),
				regexp.NewToken(regexp.SymbolTokenType, 'w'),
				regexp.NewToken(regexp.SymbolTokenType, 'h'),
				regexp.NewToken(regexp.SymbolTokenType, 'i'),
				regexp.NewToken(regexp.SymbolTokenType, 'l'),
				regexp.NewToken(regexp.SymbolTokenType, 'e'),
			},
		},
		{
			name:  "character class",
			regex: `a[ \t\n\r]*b`,
			expected: []regexp.Token{
				regexp.NewToken(regexp.SymbolTokenType, 'a'),
				regexp.NewToken(regexp.LSqBracketTokenType, '['),
				regexp.NewToken(regexp.SymbolTokenType, ' '),
				regexp.NewToken(regexp.SymbolTokenType, '\t'),
				regexp.NewToken(regexp.SymbolTokenType, '\n'),
				regexp.NewToken(regexp.SymbolTokenType, '\r'),
				regexp.NewToken(regexp.RSqBracketTokenType, ']'),
				regexp.NewToken(regexp.StarTokenType, '*'),
				regexp.NewToken(regexp.SymbolTokenType, 'b'),
			},
		},
		{
			name:  "curry bracket",
			regex: `a{123}`,
			expected: []regexp.Token{
				regexp.NewToken(regexp.SymbolTokenType, 'a'),
				regexp.NewToken(regexp.LCurryBracketTokenType, '{'),
				regexp.NewToken(regexp.DigitTokenType, rune(123)),
				regexp.NewToken(regexp.RCurryBracketTokenType, '}'),
			},
		},
		{
			name:  "curry bracket two arguments",
			regex: `a{2,10}`,
			expected: []regexp.Token{
				regexp.NewToken(regexp.SymbolTokenType, 'a'),
				regexp.NewToken(regexp.LCurryBracketTokenType, '{'),
				regexp.NewToken(regexp.DigitTokenType, rune(2)),
				regexp.NewToken(regexp.CommaTokenType, ','),
				regexp.NewToken(regexp.DigitTokenType, rune(10)),
				regexp.NewToken(regexp.RCurryBracketTokenType, '}'),
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
