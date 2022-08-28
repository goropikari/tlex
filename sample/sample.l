%{
import (
    "fmt"
    "log"
)

type Type = int
const (
    Keyword Type = iota + 1
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

func main() {
    lex := New(`
func foo000() int {
    x := 1 * 10 + 123 - 1000 / 5432

    return x
}
`)
    for {
        n, err := lex.Next()
        if err != nil {
            log.Fatal(err)
            return
        }
        switch n {
        case Keyword:
            fmt.Println("Keyword")
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
        fmt.Printf("\t %#v\n",YYtext)
    }
}
