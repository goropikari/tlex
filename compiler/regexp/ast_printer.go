package regexp

import (
	"fmt"
	"strings"
)

type ASTPrinter struct {
	depth int
	str   string
}

func NewASTPrinter() *ASTPrinter {
	return &ASTPrinter{depth: 0}
}

func (p *ASTPrinter) Print() {
	fmt.Print(p.str)
}

func (p *ASTPrinter) String() string {
	return p.str
}

func (p *ASTPrinter) VisitSumExpr(expr SumExpr) {
	p.str += p.header("SumExpr")
	p.depth++
	expr.lhs.Accept(p)
	expr.rhs.Accept(p)
	p.depth--
}

func (p *ASTPrinter) VisitConcatExpr(expr ConcatExpr) {
	p.str += p.header("ConcatExpr")
	p.depth++
	expr.lhs.Accept(p)
	expr.rhs.Accept(p)
	p.depth--
}

func (p *ASTPrinter) VisitStarExpr(expr StarExpr) {
	p.str += p.header("StarExpr")
	p.depth++
	expr.expr.Accept(p)
	p.depth--
}

func (p *ASTPrinter) VisitSymbolExpr(expr SymbolExpr) {
	p.str += p.header("SymbolExpr")
	s := fmt.Sprintf("%v%v\n", repTab(p.depth+1), string(expr.sym))
	p.str += s
}

func (p *ASTPrinter) VisitDotExpr(expr DotExpr) {
	p.str += p.header("DotExpr")
	s := fmt.Sprintf("%v%v\n", repTab(p.depth+1), ".")
	p.str += s
}

func (p *ASTPrinter) VisitRangeExpr(expr RangeExpr) {
	p.str += p.header("RangeExpr")
	strs := make([]string, 0)
	strs = append(strs, fmt.Sprintf("%v%v", repTab(p.depth+1), expr.neg))
	for _, intv := range expr.intvs {
		strs = append(strs, fmt.Sprintf("%v[%v-%v]", repTab(p.depth+1), intv.l, intv.r))
	}
	p.str += strings.Join(strs, "\n") + "\n"
}

func (p *ASTPrinter) header(name string) string {
	s := repTab(p.depth)
	s += name + "\n"
	return s
}

func repTab(n int) string {
	return strings.Repeat("\t", n)
}
