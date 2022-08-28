package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/goccy/go-graphviz"
	"github.com/goropikari/golex/automata"
	"github.com/goropikari/golex/collection"
	"github.com/goropikari/golex/compile/golex"
	"github.com/goropikari/golex/compile/regexp"
	lexer "github.com/goropikari/golex/expected"
)

func handMaid() {
	fa0 := automata.NewNFA(
		collection.NewSet[automata.State]().Insert(automata.NewState("I0")).Insert(automata.NewState("F0")),
		// collection.NewSet[rune]().Insert('a'),
		map[collection.Tuple[automata.State, rune]]collection.Set[automata.State]{
			collection.NewTuple(automata.NewState("I0"), 'a'): collection.NewSet[automata.State]().Insert(automata.NewState("F0")),
		},
		collection.NewSet[automata.State]().Insert(automata.NewState("I0")),
		collection.NewSet[automata.State]().Insert(automata.NewState("F0")),
	)

	fa1 := automata.NewNFA(
		collection.NewSet[automata.State]().Insert(automata.NewState("I1")).Insert(automata.NewState("F1")),
		// collection.NewSet[rune]().Insert('b'),
		map[collection.Tuple[automata.State, rune]]collection.Set[automata.State]{
			collection.NewTuple(automata.NewState("I1"), 'b'): collection.NewSet[automata.State]().Insert(automata.NewState("F1")),
		},
		collection.NewSet[automata.State]().Insert(automata.NewState("I1")),
		collection.NewSet[automata.State]().Insert(automata.NewState("F1")),
	)

	fa2 := automata.NewNFA(
		collection.NewSet[automata.State]().Insert(automata.NewState("I2")).Insert(automata.NewState("F2")),
		// collection.NewSet[rune]().Insert('b'),
		map[collection.Tuple[automata.State, rune]]collection.Set[automata.State]{
			collection.NewTuple(automata.NewState("I2"), 'b'): collection.NewSet[automata.State]().Insert(automata.NewState("F2")),
		},
		collection.NewSet[automata.State]().Insert(automata.NewState("I2")),
		collection.NewSet[automata.State]().Insert(automata.NewState("F2")),
	)

	fa3 := automata.NewNFA(
		collection.NewSet[automata.State]().Insert(automata.NewState("I3")).Insert(automata.NewState("F3")),
		// collection.NewSet[rune]().Insert('a'),
		map[collection.Tuple[automata.State, rune]]collection.Set[automata.State]{
			collection.NewTuple(automata.NewState("I3"), 'a'): collection.NewSet[automata.State]().Insert(automata.NewState("F3")),
		},
		collection.NewSet[automata.State]().Insert(automata.NewState("I3")),
		collection.NewSet[automata.State]().Insert(automata.NewState("F3")),
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
	// lex := regexp.NewLexer("(a*|b)cde*|fghh*")
	lex := regexp.NewLexer(regex)
	tokens := lex.Scan()
	parser := regexp.NewParser(tokens)
	ast, _ := parser.Parse()
	gen := regexp.NewCodeGenerator()
	ast.Accept(gen)

	s, _ := gen.GetNFA().ToDot()
	fmt.Println(s)
}

func convertDFA(regex string) {
	// lex := regexp.NewLexer("(a*|b)cde*|fghh*")
	lex := regexp.NewLexer(regex)
	tokens := lex.Scan()
	parser := regexp.NewParser(tokens)
	ast, _ := parser.Parse()
	gen := regexp.NewCodeGenerator()
	ast.Accept(gen)

	// s, _ := gen.GetNFA().ToDFA().Minimize().Totalize().ToDot()
	// s, _ := gen.GetNFA().ToDFA().Minimize().ToDot()
	// s, _ := gen.GetNFA().ToDFA().Totalize().ToDot()
	s, _ := gen.GetNFA().ToDFA().ToDot()
	// s, _ := gen.GetNFA().ToDFA().ToDot()
	fmt.Println(s)
}

func parse(regex string) automata.NFA {
	lex := regexp.NewLexer(regex)
	tokens := lex.Scan()
	parser := regexp.NewParser(tokens)
	ast, _ := parser.Parse()
	gen := regexp.NewCodeGenerator()
	ast.Accept(gen)

	return gen.GetNFA()
}

func lexerNFA(regexs []string) automata.NFA {
	nfas := make([]*automata.NFA, 0)
	for i, regex := range regexs {
		nfa := parse(regex)
		(&nfa).SetTokenID(automata.TokenID(i + 1))
		nfas = append(nfas, &nfa)
	}

	nfa := *nfas[0]
	for _, n := range nfas[1:] {
		nfa = nfa.Sum(*n)
	}

	return nfa
}

func testSample() {
	lex := lexer.New("aaa")
	for {
		tok, err := lex.Next()
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Printf("\t%v %v\n", tok, lexer.YYtext)
	}
}

func gen() {
	golex.Gen()
}

func main() {
	// // regex := "a"
	// // convertNFA(regex)
	// // convertDFA(regex)
	//
	// nfa := lexerNFA([]string{"a", "abb", "a*bb*", ".*"})
	// // letter := "(a|b|c|d|e|f|g|h|i|j|k|l|m|n|o|p|q|r|s|t|u|v|w|x|y|z)"
	// // digit := "(0|1|2|3|4|5|6|7|8|9)"
	// // digits := digit + digit + "*"
	// // id := fmt.Sprintf("%v(%v|%v)*", letter, letter, digit)
	// // nfa := lexerNFA([]string{
	// // 	digits,
	// // 	"if|then|begin|end|func",
	// // 	id,
	// // 	"\\+|\\-|\\*|/",
	// // 	"( |\n|\t|\r)",
	// // 	"\\.",
	// // 	".",
	// // })
	// dfa := nfa.ToDFA().LexerMinimize()
	// s, _ := dfa.RemoveBH().ToDot()
	// fmt.Println(s)

	gen()
}
