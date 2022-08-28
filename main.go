package main

import (
	"bufio"
	"os"

	"github.com/goropikari/golex/automata"
	"github.com/goropikari/golex/compile/golex"
	"github.com/goropikari/golex/compile/regexp"
)

func main() {
	f, _ := os.OpenFile(os.Args[1], os.O_RDONLY, 0755)
	// f, _ := os.OpenFile("sample/sample.l", os.O_RDONLY, 0755)
	r := bufio.NewReader(f)
	golex.Generate(r)

	// s, _ := lexerNFA([]string{"a", "abb", "a*bb*"}).ToDFA().LexerMinimize().RemoveBH().ToDot()
	// fmt.Println(s)
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
