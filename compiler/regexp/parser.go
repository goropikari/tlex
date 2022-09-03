package regexp

import (
	"errors"
	"io"

	"github.com/goropikari/tlex/automata"
	"github.com/goropikari/tlex/collection"
)

var (
	ErrParse = errors.New("parse error")
)

type Parser struct {
	tokens []Token
	pos    int
	length int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens, pos: 0, length: len(tokens)}
}

type NodeVisitor interface {
	VisitSumExpr(SumExpr)
	VisitConcatExpr(ConcatExpr)
	VisitStarExpr(StarExpr)
	VisitSymbolExpr(SymbolExpr)
	VisitDotExpr(DotExpr)
}

type RegexExpr interface {
	Accept(v NodeVisitor)
}

func (p *Parser) Parse() (RegexExpr, error) {
	return p.sum()
}

func (p *Parser) read() (Token, error) {
	if p.pos >= p.length {
		return Token{}, io.EOF
	}
	p.pos++
	return p.tokens[p.pos-1], nil
}

func (p *Parser) peek() (Token, error) {
	if p.pos >= p.length {
		return Token{}, io.EOF
	}
	return p.tokens[p.pos], nil
}

func (p *Parser) next() (Token, error) {
	if p.pos+1 >= p.length {
		return Token{}, io.EOF
	}
	return p.tokens[p.pos+1], nil
}

func (p *Parser) sum() (RegexExpr, error) {
	lhs, err := p.concat()
	if err != nil {
		return nil, err
	}

	op, err := p.peek()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return lhs, nil
		}
		return nil, err
	}

	if op.GetType() == BarTokenType {
		p.read()
		rhs, err := p.sum()
		if err != nil {
			return nil, err
		}
		return SumExpr{lhs: lhs, rhs: rhs}, nil
	}

	return lhs, nil
}

func (p *Parser) set() (RegexExpr, error) {
	neg := false
	runes := make([]rune, 0)
	var prev rune

	for {
		tok, err := p.peek()
		if err != nil {
			return nil, err
		}
		switch tok.GetType() {
		case RSqBracketTokenType:
			if prev == '-' {
				return nil, ErrParse
			}
			goto Out
		case NegationTokenType:
			prev = tok.GetRune()
			neg = true
		case MinusTokenType:
			prev = tok.GetRune()
		default:
			ru := tok.GetRune()
			if prev == '-' {
				from := runes[len(runes)-1]
				if from > ru {
					return nil, ErrParse
				}
				for t := from + 1; t < ru; t++ {
					runes = append(runes, t)
				}
			}
			runes = append(runes, ru)
			prev = ru
		}
		_, _ = p.read()
	}
Out:
	var expr RegexExpr
	if !neg {
		expr = NewSymbolExpr(runes[0])
		if len(runes) == 1 {
			return expr, nil
		}

		for i := 1; i < len(runes); i++ {
			rhs := NewSymbolExpr(runes[i])
			expr = NewSumExpr(expr, rhs)
		}
		return expr, nil
	}

	ruSet := collection.NewSet[rune]()
	for _, ru := range runes {
		ruSet.Insert(ru)
	}
	for _, ru := range automata.SupportedChars {
		if !ruSet.Contains(ru) {
			if expr == nil {
				expr = NewSymbolExpr(ru)
			} else {
				expr = NewSumExpr(expr, NewSymbolExpr(ru))
			}
		}
	}

	return expr, nil
}

func (p *Parser) concat() (RegexExpr, error) {
	lhs, err := p.star()
	if err != nil {
		return nil, err
	}

	b, err := p.peek()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return lhs, nil
		}
		return nil, err
	}

	switch b.GetType() {
	case SymbolTokenType, DotTokenType, LParenTokenType, LSqBracketTokenType:
		rhs, err := p.concat()
		if err != nil {
			return nil, err
		}
		return NewConcatExpr(lhs, rhs), nil
	default:
		return lhs, nil
	}
}

func (p *Parser) star() (RegexExpr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}
	st, err := p.peek()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return expr, nil
		}
		return nil, err
	}
	if st.GetType() == StarTokenType {
		p.read()
		return NewStarExpr(expr), nil
	}
	return expr, nil
}

func (p *Parser) primary() (RegexExpr, error) {
	s, err := p.read()
	if err != nil {
		return nil, err
	}

	switch s.GetType() {
	case SymbolTokenType:
		return NewSymbolExpr(s.GetRune()), nil
	case DotTokenType:
		return NewDotExpr(), nil
	case LParenTokenType:
		sum, err := p.sum()
		if err != nil {
			return nil, err
		}
		r, err := p.read()
		if err != nil {
			return nil, err
		}
		if r.GetType() == RParenTokenType {
			return sum, nil
		}
	case LSqBracketTokenType:
		set, err := p.set()
		if err != nil {
			return nil, err
		}
		r, err := p.read()
		if err != nil {
			return nil, err
		}
		if r.GetType() == RSqBracketTokenType {
			return set, nil
		}
	}

	return nil, ErrParse
}

type SumExpr struct {
	lhs RegexExpr
	rhs RegexExpr
}

func NewSumExpr(lhs, rhs RegexExpr) SumExpr {
	return SumExpr{lhs: lhs, rhs: rhs}
}

func (expr SumExpr) Accept(v NodeVisitor) {
	v.VisitSumExpr(expr)
}

type ConcatExpr struct {
	lhs RegexExpr
	rhs RegexExpr
}

func NewConcatExpr(lhs, rhs RegexExpr) ConcatExpr {
	return ConcatExpr{lhs: lhs, rhs: rhs}
}

func (expr ConcatExpr) Accept(v NodeVisitor) {
	v.VisitConcatExpr(expr)
}

type StarExpr struct {
	expr RegexExpr
}

func NewStarExpr(expr RegexExpr) RegexExpr {
	return StarExpr{expr: expr}
}

func (expr StarExpr) Accept(v NodeVisitor) {
	v.VisitStarExpr(expr)
}

type SymbolExpr struct {
	sym rune
}

func NewSymbolExpr(sym rune) SymbolExpr {
	return SymbolExpr{sym: sym}
}

func (expr SymbolExpr) Accept(v NodeVisitor) {
	v.VisitSymbolExpr(expr)
}

type DotExpr struct {
}

func NewDotExpr() DotExpr {
	return DotExpr{}
}

func (expr DotExpr) Accept(v NodeVisitor) {
	v.VisitDotExpr(expr)
}
