# golex: Toy implementation of Lexer Generator

golex is lexical analyzer generator such as lex.
Generated lexical analyzer is based on DFA.

This supports only ASCII string, doesn't do unicode.


```bash
$ go install github.com/goropikari/golex@v0.2.0
# golex [-src srcfile] [-pkg pkgName] [-o outfile]
$ golex -src sample.l -pkg main -o lexer.go

$ go run lexer.go
3 ab
2 abb
1 a
2022/08/28 22:49:32 EOF
exit status 1
```
