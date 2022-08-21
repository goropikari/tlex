package main

import (
	"fmt"

	"github.com/goropikari/golex/collection"
	"github.com/goropikari/golex/compile"
)

func handMaid() {
	fa0 := compile.NewNFA(
		collection.NewSet[compile.State]().Insert(compile.NewState("I0")).Insert(compile.NewState("F0")),
		// collection.NewSet[rune]().Insert('a'),
		map[compile.Tuple[compile.State, rune]]collection.Set[compile.State]{
			compile.NewTuple(compile.NewState("I0"), 'a'): collection.NewSet[compile.State]().Insert(compile.NewState("F0")),
		},
		collection.NewSet[compile.State]().Insert(compile.NewState("I0")),
		collection.NewSet[compile.State]().Insert(compile.NewState("F0")),
	)

	fa1 := compile.NewNFA(
		collection.NewSet[compile.State]().Insert(compile.NewState("I1")).Insert(compile.NewState("F1")),
		// collection.NewSet[rune]().Insert('b'),
		map[compile.Tuple[compile.State, rune]]collection.Set[compile.State]{
			compile.NewTuple(compile.NewState("I1"), 'b'): collection.NewSet[compile.State]().Insert(compile.NewState("F1")),
		},
		collection.NewSet[compile.State]().Insert(compile.NewState("I1")),
		collection.NewSet[compile.State]().Insert(compile.NewState("F1")),
	)

	fa2 := compile.NewNFA(
		collection.NewSet[compile.State]().Insert(compile.NewState("I2")).Insert(compile.NewState("F2")),
		// collection.NewSet[rune]().Insert('b'),
		map[compile.Tuple[compile.State, rune]]collection.Set[compile.State]{
			compile.NewTuple(compile.NewState("I2"), 'b'): collection.NewSet[compile.State]().Insert(compile.NewState("F2")),
		},
		collection.NewSet[compile.State]().Insert(compile.NewState("I2")),
		collection.NewSet[compile.State]().Insert(compile.NewState("F2")),
	)

	fa3 := compile.NewNFA(
		collection.NewSet[compile.State]().Insert(compile.NewState("I3")).Insert(compile.NewState("F3")),
		// collection.NewSet[rune]().Insert('a'),
		map[compile.Tuple[compile.State, rune]]collection.Set[compile.State]{
			compile.NewTuple(compile.NewState("I3"), 'a'): collection.NewSet[compile.State]().Insert(compile.NewState("F3")),
		},
		collection.NewSet[compile.State]().Insert(compile.NewState("I3")),
		collection.NewSet[compile.State]().Insert(compile.NewState("F3")),
	)

	// s, _ := fa0.Concat(fa1).ToDot() // ab
	// s, _ := fa0.Star().Sum(fa1).ToDot() // (a*|b)
	// s, _ := fa0.Sum(fa1).Concat(fa2).ToDot() // (a|b)b
	// s, _ := fa0.Star().Sum(fa1).Concat(fa2).ToDot() // (a*|b)b
	// s, _ := (fa0.Sum(fa1).Concat(fa2)).Star().ToDot()           // ((a|b)b)*
	s, _ := fa0.Star().Sum(fa1).Concat(fa2).Concat(fa3).ToDot() // (a*|b)ba
	fmt.Println(s)
}

func convertNFA() {
	lex := compile.NewLexer("(a*|b)cde*|fghh*")
	tokens := lex.Scan()
	parser := compile.NewParser(tokens)
	ast, _ := parser.Parse()
	gen := compile.NewCodeGenerator()
	ast.Accept(gen)

	s, _ := gen.GetNFA().ToDot()
	fmt.Println(s)
}

func main() {
	convertNFA()
}
