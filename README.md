<div align="center">
    <img height=200 src="image/logo.png" alt="tlex logo">
</div>

# tlex: Toy LEXical analyzer generator

tlex is lexical analyzer generator such as Lex.
This is toy implementation for my study, so don't use for production.
tlex supports Unicode.


```bash
$ go install github.com/goropikari/tlex@latest

$ tlex -h
Usage of ./tlex:
  -o string
        generated file path (default "tlex.yy.go")
  -pkg string
        generated go file package name (default "main")
  -src string
        input lexer configuration file

# tlex [-src srcfile] [-pkg output_pkg_name] [-o outfile]
$ tlex -src sample.l -pkg main -o main.go
$ go run main.go

func foo123barあいう () int {
    x := 1 * 10 + 123 - 1000 / 5432
    y := float64(x)

    return x + y
}

-----------------
Keyword
	 "func"
Identifier
	 "foo123bar"
Hiragana
	 "あいう"
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
