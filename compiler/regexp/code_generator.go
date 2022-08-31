package regexp

import (
	"github.com/goropikari/golex/automata"
	"github.com/goropikari/golex/collection"
	"github.com/goropikari/golex/utils/guid"
)

type CodeGenerator struct {
	nfa automata.NFA
}

func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{}
}

func (gen *CodeGenerator) GetNFA() automata.NFA {
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
	from := automata.NewState(automata.StateID(guid.New()))
	to := automata.NewState(automata.StateID(guid.New()))

	gen.nfa = automata.NewNFA(
		collection.NewSet[automata.State]().Insert(from).Insert(to),
		automata.NFATransition{
			collection.NewTuple(from, expr.sym): collection.NewSet[automata.State]().Insert(to),
		},
		collection.NewSet[automata.State]().Insert(from),
		collection.NewSet[automata.State]().Insert(to),
	)
}

func (gen *CodeGenerator) VisitDotExpr(expr DotExpr) {
	from := automata.NewState(automata.StateID(guid.New()))
	trans := make(automata.NFATransition)
	states := collection.NewSet[automata.State]().Insert(from)
	finStates := collection.NewSet[automata.State]()

	for _, ru := range automata.SupportedChars {
		to := automata.NewState(automata.StateID(guid.New()))
		states = states.Insert(to)
		finStates = finStates.Insert(to)
		trans[collection.NewTuple(from, ru)] = collection.NewSet[automata.State]().Insert(to)
	}

	gen.nfa = automata.NewNFA(
		states,
		trans,
		collection.NewSet[automata.State]().Insert(from),
		finStates,
	)
}
