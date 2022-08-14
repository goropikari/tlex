
{%
import "fmt"

var NumLines = 0
var NumChars = 0

type Type = int
const (
    Newline Type = iota + 1
    Char
    DoubleQuote
    Identifier
)

type Token struct {
    text string
    typ  Type
}

type ReturnType = Token
%}

%%

"\n"      { NumLines++; NumChars++; }
"\""      { fmt.Println("double quote") }
"."       { NumChars++; fmt.Println(NumChars) }
"[a-z]+"  {
    NumChars += len(yytext)
    return Token{text: yytext, typ: Identifier}, nil
}
%%
