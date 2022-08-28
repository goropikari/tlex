package automata_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/goccy/go-graphviz"
	"github.com/goropikari/golex/automata"
	"github.com/goropikari/golex/compile/regexp"
	"github.com/stretchr/testify/require"
)

func TestDFA_Accept(t *testing.T) {
	t.Parallel()

	letter := "(a|b|c|d|e|f|g|h|i|j|k|l|m|n|o|p|q|r|s|t|u|v|w|x|y|z)"
	digit := "(0|1|2|3|4|5|6|7|8|9)"
	digits := digit + digit + "*"
	id := fmt.Sprintf("%v(%v|%v)*", letter, letter, digit)

	regexs := []string{
		digits,                   // regexID: 1
		"if|then|begin|end|func", // regexID: 2
		id,                       // regexID: 3
		"\\+|\\-|\\*|/",          // regexID: 4
		"( |\n|\t|\r)",           // regexID: 5
		"\\.",                    // regexID: 6
		".",                      // regexID: 7
	}

	tests := []struct {
		name   string
		regexs []string
		given  string
		// expected
		accept  bool
		regexID automata.TokenID
	}{
		{
			name:    "digits",
			regexs:  regexs,
			given:   "123",
			accept:  true,
			regexID: 1,
		},
		{
			name:    "keyword: if",
			regexs:  regexs,
			given:   "if",
			accept:  true,
			regexID: 2,
		},
		{
			name:    "keyword: then",
			regexs:  regexs,
			given:   "then",
			accept:  true,
			regexID: 2,
		},
		{
			name:    "keyword: begin",
			regexs:  regexs,
			given:   "begin",
			accept:  true,
			regexID: 2,
		},
		{
			name:    "keyword: end",
			regexs:  regexs,
			given:   "end",
			accept:  true,
			regexID: 2,
		},
		{
			name:    "keyword: func",
			regexs:  regexs,
			given:   "func",
			accept:  true,
			regexID: 2,
		},
		{
			name:    "identifier",
			regexs:  regexs,
			given:   "ifhoge",
			accept:  true,
			regexID: 3,
		},
		{
			name:    "identifier: hoge",
			regexs:  regexs,
			given:   "hoge",
			accept:  true,
			regexID: 3,
		},
		{
			name:    "operator: +",
			regexs:  regexs,
			given:   "+",
			accept:  true,
			regexID: 4,
		},
		{
			name:    "operator: -",
			regexs:  regexs,
			given:   "-",
			accept:  true,
			regexID: 4,
		},
		{
			name:    "operator: *",
			regexs:  regexs,
			given:   "*",
			accept:  true,
			regexID: 4,
		},
		{
			name:    "operator: /",
			regexs:  regexs,
			given:   "/",
			accept:  true,
			regexID: 4,
		},
		{
			name:    "whitespace: space",
			regexs:  regexs,
			given:   " ",
			accept:  true,
			regexID: 5,
		},
		{
			name:    "dot: .",
			regexs:  regexs,
			given:   ".",
			accept:  true,
			regexID: 6,
		},
		{
			name:    "other: %",
			regexs:  regexs,
			given:   "%",
			accept:  true,
			regexID: 7,
		},
		{
			name:    "dot2: ..",
			regexs:  regexs,
			given:   "..",
			accept:  false,
			regexID: 0,
		},
		{
			name:    "identifier: start with digit",
			regexs:  regexs,
			given:   "0hoge",
			accept:  false,
			regexID: 0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dfa := lexerNFA(tt.regexs).ToDFA().LexerMinimize().RemoveBH()

			regexID, accept := dfa.Accept(tt.given)

			require.Equal(t, tt.accept, accept)
			require.Equal(t, tt.regexID, regexID)
		})
	}
}

func TestDot(t *testing.T) {
	// generate dot file
	// go test ./automata/ -run TestDot

	s, _ := lexerNFA([]string{"a", "abb", "a*bb*"}).ToDFA().LexerMinimize().RemoveBH().ToDot()
	err := os.WriteFile("ex.dot", []byte(s), 0666)
	if err != nil {
		log.Fatal(err)
	}

	graph, err := graphviz.ParseBytes([]byte(s))
	if err != nil {
		log.Fatal(err)
	}
	g := graphviz.New()
	if err := g.RenderFilename(graph, graphviz.PNG, "ex.png"); err != nil {
		log.Fatal(err)
	}
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

func parse(regex string) automata.NFA {
	lex := regexp.NewLexer(regex)
	tokens := lex.Scan()
	parser := regexp.NewParser(tokens)
	ast, _ := parser.Parse()
	gen := regexp.NewCodeGenerator()
	ast.Accept(gen)

	return gen.GetNFA()
}
