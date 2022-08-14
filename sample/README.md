```
$ go install github.com/goropikari/golex@latest
$ golex sample.l lex/gen.go
$ echo -n "foo bar baz" | go run main.go
output:
        {foo}
        {bar}
        {baz}
11 0
```
