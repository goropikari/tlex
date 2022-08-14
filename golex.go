package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"text/template"
)

func main() {
	// data, _ := os.ReadFile("/home/arch/workspace/github/golex/test.l")
	// Read config file
	data, _ := os.ReadFile(os.Args[1])
	readdata := string(data)

	const decRegex = `((?s)(\n{%\n(.*)\n%}|^{%\n(.*)\n%}))`
	re := regexp.MustCompile(decRegex)
	ds := re.FindStringSubmatch(readdata)
	if len(ds) == 0 {
		panic("You must define type alias of ReturnType")
	}
	var decBody string
	if ds[len(ds)-1] == "" {
		decBody = ds[len(ds)-2]
	} else {
		decBody = ds[len(ds)-1]
	}

	const regPattStr = `(?s)\n%%\n(.*)\n%%`
	re2 := regexp.MustCompile(regPattStr)
	regPatts := re2.FindStringSubmatch(readdata)
	if len(regPatts) == 0 {
		panic("regex pattern is required")
	}

	regLexer := NewRegexLexer(regPatts[1])
	tokens := regLexer.lex()
	regexPatterns := makeRegexPatterns(tokens)

	matchers := makeMachers(tokens)

	cfg := LexerTemplate{
		Declare:       decBody,
		RegexPatterns: regexPatterns,
		Matchers:      matchers,
	}
	s := tmpl
	t := template.Must(template.New("lexer").Parse(s))

	// Generate lexer file
	f, err := os.OpenFile(os.Args[2], os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(f, cfg)
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func makeRegexPatterns(tokens []Regex) string {
	n := len(tokens)
	patts := []string{}
	for i := 0; i < n; i++ {
		patts = append(patts, fmt.Sprintf("regPattern%v", i))
	}

	ret := make([]byte, 0)
	for i, tok := range tokens {
		ret = append(ret, []byte(fmt.Sprintf("regPattern%v = \"^%v\"\n", i, tok.pattern))...)
	}

	return fmt.Sprintf("regPatterns = %v\n%v", strings.Join(patts, " + \"|\" + "), string(ret))
}

func makeMachers(tokens []Regex) string {
	ms := make([]byte, 0)
	for i, tok := range tokens {
		s := fmt.Sprintf("if matched, _ := regexp.MatchString(regPattern%v+\"$\", yytext); matched {\n%v\n}\n", i, tok.body)
		ms = append(ms, []byte(s)...)
	}

	return string(ms)
}

type RegexLexer struct {
	body string
	pos  int
	len  int
}

type Regex struct {
	pattern string
	body    string
}

func NewRegexLexer(regBody string) *RegexLexer {
	return &RegexLexer{
		body: regBody,
		pos:  0,
		len:  len(regBody),
	}
}

func (lex *RegexLexer) lex() []Regex {
	regexs := make([]Regex, 0)
	for {
		reg, err := lex.nextToken()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return regexs
			}
		}
		body, err := lex.nextToken()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return regexs
			}
		}
		regexs = append(regexs, Regex{pattern: reg, body: body})
	}
}

func (lex *RegexLexer) skipWhitespace() {
	if lex.pos >= lex.len {
		return
	}
	for {
		switch lex.peek() {
		case ' ', '\n', '\t', '\r':
			lex.next()
			continue
		default:
			return
		}
	}
}

func (lex *RegexLexer) nextToken() (string, error) {
	lex.skipWhitespace()
	if lex.pos >= lex.len {
		return "", io.EOF
	}
	switch lex.peek() {
	case '"':
		return lex.readRegexPattern(), nil
	case '{':
		return lex.readBody(), nil
	}

	panic(errors.New("lexer error"))
}

func (lex *RegexLexer) readRegexPattern() string {
	s := make([]byte, 0)
	lex.next()
	for lex.peek() != '"' {
		s = append(s, lex.peek())
		lex.next()
	}
	lex.next()

	return string(s)
}

func (lex *RegexLexer) readBody() string {
	lex.skipWhitespace()
	s := []byte{'{'}
	lex.next()
	stack := []struct{}{{}}
	for {
		if lex.peek() == '}' {
			if len(stack) == 1 {
				break
			}
			stack = stack[0 : len(stack)-1]
		}
		if lex.peek() == '{' {
			stack = append(stack, struct{}{})
		}
		s = append(s, lex.peek())
		lex.next()
	}
	s = append(s, '}')
	lex.next()

	return string(s)
}

func (lex *RegexLexer) peek() byte {
	return lex.body[lex.pos]
}

func (lex *RegexLexer) next() {
	lex.pos++
}

type LexerTemplate struct {
	Declare       string
	RegexPatterns string
	Matchers      string
}

// Declare
// RegexPatterns
// Matchers
const tmpl = `package lex

import (
	"errors"
	"io"
	"regexp"
)

var yytext string

{{ .Declare }}

const (
{{ .RegexPatterns }}
)

type Lexer struct {
	data string
	pos  int
}

func NewLexer(data string) *Lexer {
	return &Lexer{
		data: data,
		pos:  0,
	}
}

func (lex *Lexer) NextToken() (ReturnType, error) {
	for {
		prevPos := lex.pos
		if lex.pos >= len(lex.data) {
			return ReturnType{}, io.EOF
		}
		// 自動生成コードはじめ
		if !lex.match(regPatterns) {
			break
		}
		{{ .Matchers }}
		// おわり

		if prevPos == lex.pos {
			return ReturnType{}, errors.New("infinite loop")
		}
	}

	return ReturnType{}, errors.New("no match")
}

func (lex *Lexer) match(expr string) bool {
	re := regexp.MustCompile(expr)
	re.Longest()
	yytext = re.FindString(lex.data[lex.pos:])
	lex.pos += len(yytext)

	return yytext != ""
}`
