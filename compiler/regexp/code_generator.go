package regexp

import (
	"unicode"

	"github.com/goropikari/tlex/automata"
	"github.com/goropikari/tlex/collection"
)

type CodeGenerator struct {
	nfa *automata.NFA
}

func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{}
}

func (gen *CodeGenerator) GetNFA() *automata.NFA {
	return gen.nfa
}

func (gen *CodeGenerator) VisitSumExpr(expr SumExpr) {
	expr.lhs.Accept(gen)
	lhs := gen.nfa
	expr.rhs.Accept(gen)
	rhs := gen.nfa

	gen.nfa = lhs.Sum(rhs)
}

func (gen *CodeGenerator) VisitConcatExpr(expr ConcatExpr) {
	expr.lhs.Accept(gen)
	lhs := gen.nfa
	expr.rhs.Accept(gen)
	rhs := gen.nfa

	gen.nfa = lhs.Concat(rhs)
}

func (gen *CodeGenerator) VisitStarExpr(expr StarExpr) {
	expr.expr.Accept(gen)
	gen.nfa = gen.nfa.Star()
}

func (gen *CodeGenerator) VisitSymbolExpr(expr SymbolExpr) {
	from := automata.NewStateID()
	to := automata.NewStateID()
	gen.nfa = automata.NewNFA(
		collection.NewSet[automata.StateID]().Insert(from).Insert(to),
		automata.NewEpsilonTransition(),
		automata.NewNFATransition().Set(from, automata.NewInterval(int(expr.sym), int(expr.sym)), to),
		collection.NewSet[automata.StateID]().Insert(from),
		collection.NewSet[automata.StateID]().Insert(to),
	)
}

func (gen *CodeGenerator) VisitRangeExpr(expr RangeExpr) {
	from := automata.NewStateID()
	to := automata.NewStateID()
	trans := automata.NewNFATransition()

	intvs := expr.intervals()
	for _, intv := range intvs {
		trans.Set(from, intv, to)
	}

	gen.nfa = automata.NewNFA(
		collection.NewSet[automata.StateID]().Insert(from).Insert(to),
		automata.NewEpsilonTransition(),
		trans,
		collection.NewSet[automata.StateID]().Insert(from),
		collection.NewSet[automata.StateID]().Insert(to),
	)
}

var dotRanges = []automata.Interval{
	automata.NewInterval(0, 9),
	automata.NewInterval(11, int(unicode.MaxRune)),
}

func (gen *CodeGenerator) VisitDotExpr(expr DotExpr) {
	from := automata.NewStateID()
	to := automata.NewStateID()
	initStates := collection.NewSet[automata.StateID]().Insert(from)
	trans := automata.NewNFATransition()
	states := collection.NewSet[automata.StateID]().Insert(from).Insert(to)
	finStates := collection.NewSet[automata.StateID]()

	for _, intv := range dotRanges {
		states = states.Insert(to)
		finStates = finStates.Insert(to)
		trans.Set(from, intv, to)
	}

	gen.nfa = automata.NewNFA(
		states,
		automata.NewEpsilonTransition(),
		trans,
		initStates,
		finStates,
	)
}
