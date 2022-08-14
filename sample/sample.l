
{%
var NumLines = 0
var NumChars = 0

type Type = int
const (
    Newline Type = iota + 1
    Char
)

type Token struct {text string}

type ReturnType = Token
%}

%%
"\n"      { NumLines++; NumChars++ }
"."       { NumChars++ }
"[a-z]+"  {
    NumChars += len(yytext)
    return Token{text: yytext}, nil
}
%%
