%{

import "bytes"

var nc = 0
var nw = 0
var nl = 0

%}

%%
[ \t\r]* { nc++ }
[^ \n]* {
    nc += len([]rune(YYText))
    nw++
}
\n {
    nl++
    nc++
}
%%

// This part is optional
func main() {
    program := `hello world
hello tlex
あいう αβγ
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
