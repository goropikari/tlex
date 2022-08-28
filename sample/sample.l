
%{
import (
    "fmt"
    "log"
)

type Type = int
const (
    State1 Type = iota + 1
    State2
    State3
    Other
)

%}

%%
"a"      {  return State1, nil }
"abb"      {  return State2, nil }
"a*bb*"       {  return State3, nil }
"ab" { return Other, nil }
%%

func main() {
    lex := New("ababba")
    for {
        n, err := lex.Next()
        if err != nil {
            log.Fatal(err)
            return
        }
        switch n {
        case State1:
            fmt.Println(State1, YYtext)
        case State2:
            fmt.Println(State2, YYtext)
        case State3:
            fmt.Println(State3, YYtext)
        default:
            fmt.Println(n, YYtext)
        }
    }
}
