# golex: Toy implementation of Lexical analyzer generator

golex is lexical analyzer generator such as Lex.
This is toy implementation for my study, so don't use for production.
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
$ golex -src sample.l -pkg main -o main.go
$ go run main.go

func foo123bar() int {
    x := 1 * 10 + 123 - 1000 / 5432
    y := float64(x)

    return x + y
}

-----------------
Keyword
	 "func"
Identifier
	 "foo123bar"
LParen
	 "("
RParen
	 ")"
Type
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
Identifier
	 "y"
Operator
	 ":="
Type
	 "float64"
LParen
	 "("
Identifier
	 "x"
RParen
	 ")"
Keyword
	 "return"
Identifier
	 "x"
Operator
	 "+"
Identifier
	 "y"
RBracket
	 "}"
```
