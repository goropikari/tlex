# golex: Toy implementation of Lexer Generator

golex is lexical analyzer generator such as lex.
This is toy implementation for study, so don't use for production.
golex supports only ASCII string, doesn't do unicode.


```bash
$ go install github.com/goropikari/golex@v0.3.0

$ golex -h
Usage of ./golex:
  -o string
        generated file path (default "golex.yy.go")
  -pkg string
        generated go file package name (default "main")
  -src string
        input lexer configuration file

# golex [-src srcfile] [-pkg output_pkg_name] [-o outfile]
$ golex -src sample.l -pkg main -o lexer.go

$ go run lexer.go
3 ab
2 abb
1 a
2022/08/28 22:49:32 EOF
exit status 1
```
