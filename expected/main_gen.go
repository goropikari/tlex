package lex

import (
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
)

// declare start
var numLines = 0
var numChars = 0

var yytext = ""

type Type = int

const (
	Newline Type = iota + 1
	Char
)

type Token struct {
	text string
}

type ReturnType = Token

// declare finish

// 生成される regex pattern (頭に ^ をつけただけ)
const (
	regPattern1 = "^\n"
	regPattern2 = "^."
	regPattern3 = "^[a-z]+"
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
		if !lex.match(strings.Join([]string{regPattern1, regPattern2, regPattern3}, "|")) {
			break
		}
		if matched, _ := regexp.MatchString(regPattern1+"$", yytext); matched {
			{
				numLines++
			}
		}
		if matched, _ := regexp.MatchString(regPattern2+"$", yytext); matched {
			{
				numChars++
			}
		}
		if matched, _ := regexp.MatchString(regPattern3+"$", yytext); matched {
			{
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
	// fmt.Println(expr, yytext)

	return yytext != ""
}

func main() {
	lex := NewLexer("hoge\npiyo")
	for {
		token, err := lex.NextToken()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}

			log.Fatal(err)
		}
		fmt.Printf("token: %v\n", token)
	}

	// fmt.Println(numLines)
}
