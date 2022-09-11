package generator_test

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"

	"github.com/goropikari/tlex/compiler/generator"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name  string
		given string
	}{
		{
			name: "parser test",
			given: `
%{

import (
	"fmt"
	"log"
	"bytes"
)

// generated lexer returned types are (int, error).
const (
	Keyword int = iota + 1
	Type
	Identifier
	Digit
	Whitespace
	LParen
	RParen
	LBracket
	RBracket
	Operator
	Hiragana
)

%}

%%
if|for|while|func|return { return Keyword, nil }
int|float64 { return Type, nil }
[a-zA-Z][a-zA-Z0-9]* { return Identifier, nil }
[1-9][0-9]* { return Digit, nil }
[ \t\n\r]* { }
"(" { return LParen, nil }
")" { return RParen, nil }
"{" { return LBracket, nil }
"}" { return RBracket, nil }
"+" { return Operator, nil }
"-" { return Operator, nil }
"*" { return Operator, nil }
"/" { return Operator, nil }
":=" { return Operator, nil }
"==" { return Operator, nil }
"!=" { return Operator, nil }
[ぁ-ゔ]* { return Hiragana, nil }
. {}
%%

user defined code
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewBufferString(tt.given)
			p := generator.NewParser(bufio.NewReader(r))
			def, rules, userCode := p.Parse()
			fmt.Println(def, rules, userCode)
		})
	}
}
