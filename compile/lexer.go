package compile

import (
	"errors"
	"io"
)

var (
	ErrInvalidRegex = errors.New("invalid regular expression")
)

type TokenType int

const (
	SymbolTokenType TokenType = iota + 1
	DotTokenType
	StarTokenType
	LParenTokenType
	RParenTokenType
	BarTokenType
)

type Token struct {
	typ TokenType
	val rune
}

func NewToken(typ TokenType, val rune) Token {
	return Token{typ: typ, val: val}
}

func (tok Token) GetType() TokenType {
	return tok.typ
}

func (tok Token) GetRune() rune {
	return tok.val
}

type Lexer struct {
	regexp []rune
	tokens []Token
	pos    int
	length int
}

func NewLexer(regexp string) *Lexer {
	return &Lexer{regexp: []rune(regexp), pos: 0, length: len(regexp)}
}

func (lex *Lexer) peek() (rune, error) {
	if lex.pos >= len(lex.regexp) {
		return 0, io.EOF
	}
	return lex.regexp[lex.pos], nil
}

// func (lex *Lexer) next() (rune, error) {
// 	if lex.pos+1 >= lex.length {
// 		return 0, io.EOF
// 	}
// 	return lex.regexp[lex.pos+1], nil
// }

func (lex *Lexer) read() (rune, error) {
	ru, err := lex.peek()
	if err != nil {
		return 0, err
	}
	lex.advance()

	return ru, nil
}

func (lex *Lexer) advance() {
	lex.pos++
}

func (lex *Lexer) Scan() []Token {
	for {
		ru, err := lex.read()
		if errors.Is(err, io.EOF) {
			return lex.tokens
		}

		var typ TokenType
		switch ru {
		case '\\':
			ru2, err := lex.read()
			if errors.Is(err, io.EOF) {
				panic(ErrInvalidRegex)
			}
			switch ru2 {
			case '.', '+', '-', '*', '(', ')', '[', ']':
				ru = ru2
			default:
				panic(ErrInvalidRegex)
			}
			typ = SymbolTokenType
		case '*':
			typ = StarTokenType
		case '(':
			typ = LParenTokenType
		case ')':
			typ = RParenTokenType
		case '|':
			typ = BarTokenType
		case '.':
			typ = DotTokenType
		default:
			typ = SymbolTokenType
		}
		lex.tokens = append(lex.tokens, NewToken(typ, ru))
	}
}
