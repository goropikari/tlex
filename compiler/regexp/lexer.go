package regexp

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
	MinusTokenType
	LParenTokenType
	RParenTokenType
	LSqBracketTokenType
	RSqBracketTokenType
	LCurryBracketTokenType
	RCurryBracketTokenType
	BarTokenType
	NegationTokenType
	DigitTokenType
	CommaTokenType
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

func (lex *Lexer) read() (rune, error) {
	b, err := lex.peek()
	if err != nil {
		return 0, err
	}
	lex.advance()

	return b, nil
}

func (lex *Lexer) advance() {
	lex.pos++
}

func (lex *Lexer) Scan() []Token {
	r, err := lex.peek()
	if errors.Is(err, io.EOF) {
		return lex.tokens
	}

	if r == '"' {
		return lex.scanStringLiteral()
	}

	for {
		var typ TokenType
		r, err := lex.read()
		if errors.Is(err, io.EOF) {
			return lex.tokens
		}
		switch r {
		case '\\':
			r2, err := lex.read()
			if errors.Is(err, io.EOF) {
				panic(ErrInvalidRegex)
			}
			switch r2 {
			case 'a':
				r = '\a'
			case 'b':
				r = '\b'
			case 'f':
				r = '\f'
			case 'n':
				r = '\n'
			case 'r':
				r = '\r'
			case 't':
				r = '\t'
			case 'v':
				r = '\v'
			case '.', '+', '-', '*', '/', '(', ')', '[', ']', '{', '}', '\\':
				r = r2
			default:
				panic(ErrInvalidRegex)
			}
			typ = SymbolTokenType
		case '*':
			typ = StarTokenType
		case '-':
			typ = MinusTokenType
		case '(':
			typ = LParenTokenType
		case ')':
			typ = RParenTokenType
		case '[':
			lex.tokens = append(lex.tokens, NewToken(LSqBracketTokenType, r))
			r2, err := lex.read()
			if errors.Is(err, io.EOF) {
				panic(ErrInvalidRegex)
			}
			r = r2
			switch r {
			case '^':
				typ = NegationTokenType
			default:
				typ = SymbolTokenType
			}
		case ']':
			typ = RSqBracketTokenType
		case '{':
			lex.tokens = append(lex.tokens, NewToken(LCurryBracketTokenType, r))
			lower := lex.scanDigit()
			lex.tokens = append(lex.tokens, NewToken(DigitTokenType, rune(lower)))
			r, err = lex.read()
			if err != nil {
				panic(err)
			}
			if r == ',' {
				lex.tokens = append(lex.tokens, NewToken(CommaTokenType, r))
				upper := lex.scanDigit()
				lex.tokens = append(lex.tokens, NewToken(DigitTokenType, rune(upper)))
				r, err = lex.read()
				if err != nil {
					panic(err)
				}
				if r == '}' {
					typ = RCurryBracketTokenType
				} else {
					panic(ErrInvalidRegex)
				}
			} else if r == '}' {
				typ = RCurryBracketTokenType
			} else {
				panic(ErrInvalidRegex)
			}
		case '|':
			typ = BarTokenType
		case '.':
			typ = DotTokenType
		default:
			typ = SymbolTokenType
		}
		lex.tokens = append(lex.tokens, NewToken(typ, r))
	}
}

func (lex *Lexer) scanDigit() int {
	num := 0
	for {
		r, err := lex.peek()
		if err != nil {
			panic(err)
		}
		if '0' <= r && r <= '9' {
			num = num*10 + int(r-'0')
		} else {
			return num
		}
		lex.read()
	}
}

func (lex *Lexer) scanStringLiteral() []Token {
	_, err := lex.read()
	if err != nil {
		panic(err)
	}
	var prev rune
	tokens := make([]Token, 0)
	for {
		r, err := lex.read()
		if err != nil {
			panic(err)
		}
		switch prev {
		case '\\':
			switch r {
			case 'a':
				tokens = append(tokens, NewToken(SymbolTokenType, '\a'))
				r = '\a'
			case 'b':
				tokens = append(tokens, NewToken(SymbolTokenType, '\b'))
				r = '\b'
			case 'f':
				tokens = append(tokens, NewToken(SymbolTokenType, '\f'))
				r = '\f'
			case 'n':
				tokens = append(tokens, NewToken(SymbolTokenType, '\n'))
				r = '\n'
			case 'r':
				tokens = append(tokens, NewToken(SymbolTokenType, '\r'))
				r = '\r'
			case 't':
				tokens = append(tokens, NewToken(SymbolTokenType, '\t'))
				r = '\t'
			case 'v':
				tokens = append(tokens, NewToken(SymbolTokenType, '\v'))
				r = '\v'
			case '0':
				tokens = append(tokens, NewToken(SymbolTokenType, 0))
				r = 0
			case '\\':
				tokens = append(tokens, NewToken(SymbolTokenType, '\\'))
				r = -1
			case '"':
				tokens = append(tokens, NewToken(SymbolTokenType, '"'))
				r = -1
			}
		default:
			switch r {
			case '"':
				return tokens
			case '\\':
				break
			default:
				tokens = append(tokens, NewToken(SymbolTokenType, r))
			}
		}

		prev = r
	}
}
