package compile

import (
	"github.com/google/uuid"
	"github.com/goropikari/golex/automaton"
	"github.com/goropikari/golex/collection"
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
	from := automaton.NewState(uuid.New().String())
	to := automaton.NewState(uuid.New().String())

	gen.nfa = automaton.NewNFA(
		collection.NewSet[automaton.State]().Insert(from).Insert(to),
		automaton.Transition{
			// collection.NewTuple[automaton.State, rune](from, expr.sym): collection.NewSet[automaton.State]().Insert(to),
			collection.NewTuple(from, expr.sym): collection.NewSet[automaton.State]().Insert(to),
		},
		collection.NewSet[automaton.State]().Insert(from),
		collection.NewSet[automaton.State]().Insert(to),
	)
}
