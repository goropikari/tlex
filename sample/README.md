
Lexical analize following function defition.
```go
func foo000() int {
    x := 1 * 10 + 123 - 1000 / 5432

    return x
}
```


```bash
$ go install github.com/goropikari/golex@v0.2.0
$ golex sample.l
$ go run golex.yy.go
Keyword
         "func"
Identifier
         "foo000"
LParen
         "("
RParen
         ")"
Identifier
         "int"
LBracket
         "{"
Identifier
         "x"
Operator
         ":="
Digit
         "1"
Operator
         "*"
Digit
         "10"
Operator
         "+"
Digit
         "123"
Operator
         "-"
Digit
         "1000"
Operator
         "/"
Digit
         "5432"
Keyword
         "return"
Identifier
         "x"
RBracket
         "}"
2022/08/29 01:39:20 EOF
exit status 1
```


`sample.l`
```go
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
"if|for|while|func" {  return Keyword, nil }
"[a-zA-Z][a-zA-Z0-9]*" {  return Identifier, nil }
"[1-9][0-9]*" {  return Digit, nil }
"[ \t\n\r]*" { return Whitespace, nil }
"\\(" { return LParen, nil }
"\\)" { return RParen, nil }
"{" { return LBracket, nil }
"}" { return RBracket, nil }
"[\\+|\\-|\\*|/|:=|==|!=]" { return Operator, nil }
"." {}
%%

func main() {
    lex := New(`
func foo000() {
    x := 1 * 10 + 123 - 1000 / 5432
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
```
