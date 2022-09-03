%{

import "bytes"

var nc = 0
var nw = 0
var nl = 0

%}

%%
"[a-zA-Z0-9()]*" {
    nc += len(YYText)
    nw++
}
"[ \t\r]*" { nc++ }
"\n" {
    nl++
    nc++
}
%%

// This part is optional
func main() {
    program := `hello world
hello tlex
`
    fmt.Print(program)
    fmt.Println("-----------------")

    lex := New(bytes.NewReader([]byte(program)))
    for {
        _, err := lex.Next()
        if err != nil {
            break
        }
    }
    fmt.Printf("number of lines: %d\n", nl)
    fmt.Printf("number of words: %d\n", nw)
    fmt.Printf("number of chars: %d\n", nc)
}
