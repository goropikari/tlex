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
	VisitRangeExpr(RangeExpr)
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

type interval struct {
	l int
	r int
}

func newInterval(l, r int) interval {
	return interval{l: l, r: r}
}

func newIntervalRune(r rune) interval {
	return newInterval(int(r), int(r))
}

func (p *Parser) set() (RegexExpr, error) {
	neg := false
	var prev rune
	deq := collection.NewDeque[interval]()

	isFirst := true
	for {
		tok, err := p.peek()
		if err != nil {
			return nil, err
		}
		switch tok.GetType() {
		case RSqBracketTokenType:
			if prev == '-' {
				deq.PushBack(newIntervalRune('-'))
			}
			goto Out
		case NegationTokenType:
			ru := tok.GetRune()
			if isFirst {
				neg = true
			} else {
				deq.PushBack(newIntervalRune(ru))
			}
			prev = ru
		case MinusTokenType:
			if prev == '-' {
				return nil, ErrParse
			}
			prev = tok.GetRune()
		default:
			ru := tok.GetRune()
			if prev == '-' {
				if deq.Size() == 0 {
					return nil, ErrParse
				}
				intv := deq.Back()
				deq.PopBack()
				if intv.l > int(ru) {
					return nil, ErrParse
				}
				intv.r = int(ru)
				deq.PushBack(intv)
			} else {
				deq.PushBack(newIntervalRune(ru))
			}
			prev = ru
		}
		if _, err := p.read(); err != nil {
			return nil, err
		}
		isFirst = false
	}

Out:
	var expr RegexExpr
	intvs := make([]interval, 0)
	for deq.Size() > 0 {
		intv := deq.Front()
		deq.PopFront()
		intvs = append(intvs, intv)
	}
	expr = NewRangeExpr(neg, intvs)

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

type RangeExpr struct {
	neg   bool
	intvs []interval
}

func NewRangeExpr(neg bool, intvs []interval) RangeExpr {
	return RangeExpr{neg: neg, intvs: intvs}
}

func (expr RangeExpr) Accept(v NodeVisitor) {
	v.VisitRangeExpr(expr)
}

func (expr RangeExpr) intervals() []automata.Interval {
	tmpIntvs := make([]automata.Interval, 0)
	for _, intv := range expr.intvs {
		tmpIntvs = append(tmpIntvs, automata.NewInterval(intv.l, intv.r))
	}
	if !expr.neg {
		return tmpIntvs
	}

	deq := collection.NewDeque[automata.Interval]()
	for _, intv := range automata.UnicodeRange {
		deq.PushBack(intv)
	}

	intvs := make([]automata.Interval, 0)
	for deq.Size() > 0 {
		fr := deq.Front()
		deq.PopFront()
		ok := true
		for _, intv := range tmpIntvs {
			if fr.Overlap(intv) {
				ok = false
				ls := fr.Difference(intv)
				for _, t := range ls {
					deq.PushBack(t)
				}
				break
			}
		}
		if ok {
			intvs = append(intvs, fr)
		}
	}

	return intvs
}

type DotExpr struct {
}

func NewDotExpr() DotExpr {
	return DotExpr{}
}

func (expr DotExpr) Accept(v NodeVisitor) {
	v.VisitDotExpr(expr)
}
