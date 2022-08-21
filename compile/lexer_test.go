package compile_test

import (
	"testing"

	"github.com/goropikari/golex/compile"
	"github.com/stretchr/testify/require"
)

func TestLexer_Scan(t *testing.T) {
	tests := []struct {
		name     string
		regex    string
		expected []compile.Token
	}{
		{
			name:  "lexer test",
			regex: "a(b|c*)deあいう",
			expected: []compile.Token{
				compile.NewToken(compile.SymbolTokenType, 'a'),
				compile.NewToken(compile.LParenTokenType, '('),
				compile.NewToken(compile.SymbolTokenType, 'b'),
				compile.NewToken(compile.BarTokenType, '|'),
				compile.NewToken(compile.SymbolTokenType, 'c'),
				compile.NewToken(compile.StarTokenType, '*'),
				compile.NewToken(compile.RParenTokenType, ')'),
				compile.NewToken(compile.SymbolTokenType, 'd'),
				compile.NewToken(compile.SymbolTokenType, 'e'),
				compile.NewToken(compile.SymbolTokenType, 'あ'),
				compile.NewToken(compile.SymbolTokenType, 'い'),
				compile.NewToken(compile.SymbolTokenType, 'う'),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := compile.NewLexer(tt.regex)
			toks := lexer.Scan()

			require.Equal(t, tt.expected, toks)
		})
	}

}
