package regexp_test

import (
	"testing"

	"github.com/goropikari/tlex/compiler/regexp"
	"github.com/stretchr/testify/require"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []regexp.Token
		expected string
	}{
		{
			name: "parser test",
			tokens: []regexp.Token{ // a(b|c*)deあいう|fg*hi|.*
				regexp.NewToken(regexp.SymbolTokenType, 'a'),
				regexp.NewToken(regexp.LParenTokenType, '('),
				regexp.NewToken(regexp.SymbolTokenType, 'b'),
				regexp.NewToken(regexp.BarTokenType, '|'),
				regexp.NewToken(regexp.SymbolTokenType, 'c'),
				regexp.NewToken(regexp.StarTokenType, '*'),
				regexp.NewToken(regexp.RParenTokenType, ')'),
				regexp.NewToken(regexp.SymbolTokenType, 'd'),
				regexp.NewToken(regexp.SymbolTokenType, 'e'),
				regexp.NewToken(regexp.SymbolTokenType, 'あ'),
				regexp.NewToken(regexp.SymbolTokenType, 'い'),
				regexp.NewToken(regexp.SymbolTokenType, 'う'),
				regexp.NewToken(regexp.BarTokenType, '|'),
				regexp.NewToken(regexp.SymbolTokenType, 'f'),
				regexp.NewToken(regexp.SymbolTokenType, 'g'),
				regexp.NewToken(regexp.StarTokenType, '*'),
				regexp.NewToken(regexp.SymbolTokenType, 'h'),
				regexp.NewToken(regexp.SymbolTokenType, 'i'),
				regexp.NewToken(regexp.BarTokenType, '|'),
				regexp.NewToken(regexp.DotTokenType, '.'),
				regexp.NewToken(regexp.StarTokenType, '*'),
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
			parser := regexp.NewParser(tt.tokens)
			expr, err := parser.Parse()

			printer := regexp.NewASTPrinter()
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
			given: "a(b|c*)deあいう|fg*hi|.*",
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
		{
			name:  "lexer & parser test 2",
			given: "(.*)*",
			expected: `
StarExpr
	StarExpr
		DotExpr
			.
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := regexp.NewLexer(tt.given)
			tokens := lexer.Scan()
			parser := regexp.NewParser(tokens)
			expr, err := parser.Parse()

			printer := regexp.NewASTPrinter()
			expr.Accept(printer)
			// printer.Print()

			require.NoError(t, err)
			require.Equal(t, tt.expected, "\n"+printer.String())
		})
	}
}
