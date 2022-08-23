package compile_test

import (
	"testing"

	"github.com/goropikari/golex/compile"
	"github.com/stretchr/testify/require"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []compile.Token
		expected string
	}{
		{
			name: "parser test",
			tokens: []compile.Token{ // a(b|c*)deあいう|fg*hi|.*
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
				compile.NewToken(compile.BarTokenType, '|'),
				compile.NewToken(compile.SymbolTokenType, 'f'),
				compile.NewToken(compile.SymbolTokenType, 'g'),
				compile.NewToken(compile.StarTokenType, '*'),
				compile.NewToken(compile.SymbolTokenType, 'h'),
				compile.NewToken(compile.SymbolTokenType, 'i'),
				compile.NewToken(compile.BarTokenType, '|'),
				compile.NewToken(compile.DotTokenType, '.'),
				compile.NewToken(compile.StarTokenType, '*'),
			},
			expected: `
SumExpr
	ConcatExpr
		SymbolExpr
			a
		ConcatExpr
			SumExpr
				SymbolExpr
					b
				StarExpr
					SymbolExpr
						c
			ConcatExpr
				SymbolExpr
					d
				ConcatExpr
					SymbolExpr
						e
					ConcatExpr
						SymbolExpr
							あ
						ConcatExpr
							SymbolExpr
								い
							SymbolExpr
								う
	SumExpr
		ConcatExpr
			SymbolExpr
				f
			ConcatExpr
				StarExpr
					SymbolExpr
						g
				ConcatExpr
					SymbolExpr
						h
					SymbolExpr
						i
		StarExpr
			DotExpr
				.
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := compile.NewParser(tt.tokens)
			expr, err := parser.Parse()

			printer := compile.NewASTPrinter()
			expr.Accept(printer)
			// printer.Print()

			require.NoError(t, err)
			require.Equal(t, tt.expected, "\n"+printer.String())
		})
	}
}

func TestParser_Lexer_Parse(t *testing.T) {
	tests := []struct {
		name     string
		given    string
		expected string
	}{
		{
			name:  "lexer & parser test",
			given: "a(b|c*)deあいう|fg*hi|.*|\t",
			expected: `
SumExpr
	ConcatExpr
		SymbolExpr
			a
		ConcatExpr
			SumExpr
				SymbolExpr
					b
				StarExpr
					SymbolExpr
						c
			ConcatExpr
				SymbolExpr
					d
				ConcatExpr
					SymbolExpr
						e
					ConcatExpr
						SymbolExpr
							あ
						ConcatExpr
							SymbolExpr
								い
							SymbolExpr
								う
	SumExpr
		ConcatExpr
			SymbolExpr
				f
			ConcatExpr
				StarExpr
					SymbolExpr
						g
				ConcatExpr
					SymbolExpr
						h
					SymbolExpr
						i
		StarExpr
			DotExpr
				.
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := compile.NewLexer(tt.given)
			tokens := lexer.Scan()
			parser := compile.NewParser(tokens)
			expr, err := parser.Parse()

			printer := compile.NewASTPrinter()
			expr.Accept(printer)
			// printer.Print()

			require.NoError(t, err)
			require.Equal(t, tt.expected, "\n"+printer.String())
		})
	}
}
