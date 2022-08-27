package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/goccy/go-graphviz"
	"github.com/goropikari/golex/automaton"
	"github.com/goropikari/golex/collection"
	"github.com/goropikari/golex/compile"
)

func handMaid() {
	fa0 := automaton.NewNFA(
		collection.NewSet[automaton.State]().Insert(automaton.NewState("I0")).Insert(automaton.NewState("F0")),
		// collection.NewSet[rune]().Insert('a'),
		map[collection.Tuple[automaton.State, rune]]collection.Set[automaton.State]{
			collection.NewTuple(automaton.NewState("I0"), 'a'): collection.NewSet[automaton.State]().Insert(automaton.NewState("F0")),
		},
		collection.NewSet[automaton.State]().Insert(automaton.NewState("I0")),
		collection.NewSet[automaton.State]().Insert(automaton.NewState("F0")),
	)

	fa1 := automaton.NewNFA(
		collection.NewSet[automaton.State]().Insert(automaton.NewState("I1")).Insert(automaton.NewState("F1")),
		// collection.NewSet[rune]().Insert('b'),
		map[collection.Tuple[automaton.State, rune]]collection.Set[automaton.State]{
			collection.NewTuple(automaton.NewState("I1"), 'b'): collection.NewSet[automaton.State]().Insert(automaton.NewState("F1")),
		},
		collection.NewSet[automaton.State]().Insert(automaton.NewState("I1")),
		collection.NewSet[automaton.State]().Insert(automaton.NewState("F1")),
	)

	fa2 := automaton.NewNFA(
		collection.NewSet[automaton.State]().Insert(automaton.NewState("I2")).Insert(automaton.NewState("F2")),
		// collection.NewSet[rune]().Insert('b'),
		map[collection.Tuple[automaton.State, rune]]collection.Set[automaton.State]{
			collection.NewTuple(automaton.NewState("I2"), 'b'): collection.NewSet[automaton.State]().Insert(automaton.NewState("F2")),
		},
		collection.NewSet[automaton.State]().Insert(automaton.NewState("I2")),
		collection.NewSet[automaton.State]().Insert(automaton.NewState("F2")),
	)

	fa3 := automaton.NewNFA(
		collection.NewSet[automaton.State]().Insert(automaton.NewState("I3")).Insert(automaton.NewState("F3")),
		// collection.NewSet[rune]().Insert('a'),
		map[collection.Tuple[automaton.State, rune]]collection.Set[automaton.State]{
			collection.NewTuple(automaton.NewState("I3"), 'a'): collection.NewSet[automaton.State]().Insert(automaton.NewState("F3")),
		},
		collection.NewSet[automaton.State]().Insert(automaton.NewState("I3")),
		collection.NewSet[automaton.State]().Insert(automaton.NewState("F3")),
	)

	// s, _ := fa0.Concat(fa1).ToDot() // ab
	// s, _ := fa0.Star().Sum(fa1).ToDot() // (a*|b)
	// s, _ := fa0.Sum(fa1).Concat(fa2).ToDot() // (a|b)b
	// s, _ := fa0.Star().Sum(fa1).Concat(fa2).ToDot() // (a*|b)b
	// s, _ := (fa0.Sum(fa1).Concat(fa2)).Star().ToDot()           // ((a|b)b)*
	s, _ := fa0.Star().Sum(fa1).Concat(fa2).Concat(fa3).ToDot() // (a*|b)ba
	fmt.Println(s)
}

func Dot() {
	g := graphviz.New()
	graph, _ := g.Graph()
	defer func() {
		if err := graph.Close(); err != nil {
			log.Fatal(err)
		}
		g.Close()
	}()
	graph.SetRankDir("LR") // 図を横長にする

	n, _ := graph.CreateNode("S0")
	e1, _ := graph.CreateEdge("id", n, n)
	e1.SetLabel("a\nb")
	// graph.CreateEdge("bcb", n, n)
	// e2.SetLabel("b")
	// e2.SetID("b")
	// e3, _ := graph.CreateEdge("c", n, n)
	// e3.SetLabel("c")

	var buf bytes.Buffer
	g.Render(graph, "dot", &buf)

	fmt.Println(buf.String())
}

func convertNFA(regex string) {
	// lex := compile.NewLexer("(a*|b)cde*|fghh*")
	lex := compile.NewLexer(regex)
	tokens := lex.Scan()
	parser := compile.NewParser(tokens)
	ast, _ := parser.Parse()
	gen := compile.NewCodeGenerator()
	ast.Accept(gen)

	s, _ := gen.GetNFA().ToDot()
	fmt.Println(s)
}

func convertDFA(regex string) {
	// lex := compile.NewLexer("(a*|b)cde*|fghh*")
	lex := compile.NewLexer(regex)
	tokens := lex.Scan()
	parser := compile.NewParser(tokens)
	ast, _ := parser.Parse()
	gen := compile.NewCodeGenerator()
	ast.Accept(gen)

	// s, _ := gen.GetNFA().ToDFA().Minimize().Totalize().ToDot()
	// s, _ := gen.GetNFA().ToDFA().Minimize().ToDot()
	// s, _ := gen.GetNFA().ToDFA().Totalize().ToDot()
	s, _ := gen.GetNFA().ToDFA().ToDot()
	// s, _ := gen.GetNFA().ToDFA().ToDot()
	fmt.Println(s)
}

func parse(regex string) automaton.NFA {
	lex := compile.NewLexer(regex)
	tokens := lex.Scan()
	parser := compile.NewParser(tokens)
	ast, _ := parser.Parse()
	gen := compile.NewCodeGenerator()
	ast.Accept(gen)

	return gen.GetNFA()
}

func lexerNFA(regexs []string) automaton.NFA {
	nfas := make([]*automaton.NFA, 0)
	for i, regex := range regexs {
		nfa := parse(regex)
		(&nfa).SetTokenID(automaton.TokenID(i + 1))
		nfas = append(nfas, &nfa)
	}

	nfa := *nfas[0]
	for _, n := range nfas[1:] {
		nfa = nfa.Sum(*n)
	}

	return nfa
}

func main() {
	// regex := "a"
	// convertNFA(regex)
	// convertDFA(regex)

	nfa := lexerNFA([]string{"a", "abb", "a*bb*", ".*"})
	// letter := "(a|b|c|d|e|f|g|h|i|j|k|l|m|n|o|p|q|r|s|t|u|v|w|x|y|z)"
	// digit := "(0|1|2|3|4|5|6|7|8|9)"
	// digits := digit + digit + "*"
	// id := fmt.Sprintf("%v(%v|%v)*", letter, letter, digit)
	// nfa := lexerNFA([]string{
	// 	digits,
	// 	"if|then|begin|end|func",
	// 	id,
	// 	"\\+|\\-|\\*|/",
	// 	"( |\n|\t|\r)",
	// 	"\\.",
	// 	".",
	// })
	dfa := nfa.ToDFA().LexerMinimize()
	s, _ := dfa.RemoveBH().ToDot()
	fmt.Println(s)
}
