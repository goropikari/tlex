%{

import (
    "fmt"
    "log"
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
)

%}

%%
"if|for|while|func|return" { return Keyword, nil }
"int|float64" { return Type, nil }
"[a-zA-Z][a-zA-Z0-9]*" { return Identifier, nil }
"[1-9][0-9]*" { return Digit, nil }
"[ \t\n\r]*" { }
"\\(" { return LParen, nil }
"\\)" { return RParen, nil }
"{" { return LBracket, nil }
"}" { return RBracket, nil }
"\\+|\\-|\\*|/|:=|==|!=" { return Operator, nil }
"." {}
%%

// This part is optional
func main() {
    program := `
func foo123bar() int {
    x := 1 * 10 + 123 - 1000 / 5432
    y := float64(x)

    return x + y
}
`
    fmt.Println(program)
    fmt.Println("-----------------")

    lex := New(program)
    for {
        n, err := lex.Next()
        if err != nil {
            return
        }
        switch n {
        case Keyword:
            fmt.Println("Keyword")
        case Type:
            fmt.Println("Type")
        case Identifier:
            fmt.Println("Identifier")
        case Digit:
            fmt.Println("Digit")
        case Whitespace:
            fmt.Println("Whitespace")
        case LParen:
            fmt.Println("LParen")
        case RParen:
            fmt.Println("RParen")
        case LBracket:
            fmt.Println("LBracket")
        case RBracket:
            fmt.Println("RBracket")
        case Operator:
            fmt.Println("Operator")
        }
        fmt.Printf("\t %#v\n",lex.YYText)
    }
}
