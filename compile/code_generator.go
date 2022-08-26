package compile

import (
	"github.com/goropikari/golex/automaton"
	"github.com/goropikari/golex/collection"
	"github.com/goropikari/golex/utils/guid"
)

type CodeGenerator struct {
	nfa automaton.NFA
}

func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{}
}

func (gen *CodeGenerator) GetNFA() automaton.NFA {
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
	from := automaton.NewState(guid.New())
	to := automaton.NewState(guid.New())

	gen.nfa = automaton.NewNFA(
		collection.NewSet[automaton.State]().Insert(from).Insert(to),
		automaton.NFATransition{
			// collection.NewTuple[automaton.State, rune](from, expr.sym): collection.NewSet[automaton.State]().Insert(to),
			collection.NewTuple(from, expr.sym): collection.NewSet[automaton.State]().Insert(to),
		},
		collection.NewSet[automaton.State]().Insert(from),
		collection.NewSet[automaton.State]().Insert(to),
	)
}

func (gen *CodeGenerator) VisitDotExpr(expr DotExpr) {
	from := automaton.NewState(guid.New())
	trans := make(automaton.NFATransition)
	states := collection.NewSet[automaton.State]().Insert(from)
	finStates := collection.NewSet[automaton.State]()

	for _, ru := range automaton.SupportedChars {
		to := automaton.NewState(guid.New())
		states = states.Insert(to)
		finStates = finStates.Insert(to)
		trans[collection.NewTuple(from, ru)] = collection.NewSet[automaton.State]().Insert(to)
	}

	gen.nfa = automaton.NewNFA(
		states,
		trans,
		collection.NewSet[automaton.State]().Insert(from),
		finStates,
	)
}
