```
$ go install github.com/goropikari/golex@latest
$ golex sample.l lex/gen.go
$ printf 'foo "a \nbar\nbaz\n' | go run main.go
output:
        {foo 4}
4
double quote
5
6
        {bar 4}
        {baz 4}
15 3
```
`a` doesn't match `[a-z]+` because `.` is higher precedence than `[a-z]+` for single character.
