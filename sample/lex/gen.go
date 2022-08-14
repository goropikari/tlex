package lex

import (
	"errors"
	"io"
	"regexp"
)

var yytext string

var NumLines = 0
var NumChars = 0

type Type = int
const (
    Newline Type = iota + 1
    Char
)

type Token struct {text string}

type ReturnType = Token

const (
regPatterns = regPattern0 + "|" + regPattern1 + "|" + regPattern2
regPattern0 = "^\n"
regPattern1 = "^."
regPattern2 = "^[a-z]+"

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
		if matched, _ := regexp.MatchString(regPattern0+"$", yytext); matched {
{ NumLines++; NumChars++ }
}
if matched, _ := regexp.MatchString(regPattern1+"$", yytext); matched {
{ NumChars++ }
}
if matched, _ := regexp.MatchString(regPattern2+"$", yytext); matched {
{
    NumChars += len(yytext)
    return Token{text: yytext}, nil
}
}

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
}