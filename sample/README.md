```bash
$ go install github.com/goropikari/golex@v0.2.0
$ golex sample.l
$ go run golex.yy.go

2 abb
3 ab
1 a
1 a
1 a
```


`golex.yy.go`
```go
// "a"      {  return State1, nil }
// "abb"      {  return State2, nil }
// "a*bb*"       {  return State3, nil }

func main() {
    lex := New("abbabaaa")
    for {
        n, err := lex.Next()
        if err != nil {
            return
        }
        switch n {
        case State1:
            fmt.Println(State1, YYtext)
        case State2:
            fmt.Println(State2, YYtext)
        case State3:
            fmt.Println(State3, YYtext)
        }
    }
}
```
