package lex

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"text/template"

	"github.com/goropikari/tlex/automata"
	"github.com/goropikari/tlex/collection"
	"github.com/goropikari/tlex/compiler/regexp"
	"golang.org/x/tools/imports"
)

type LexerTemplate struct {
	PackageName          string
	EmbeddedTmpl         string
	StateIDToRegexIDTmpl string
	FinStatesTmpl        string
	TransitionTableTmpl  string
	RegexActionsTmpl     string
	UserCodeTmpl         string
}

func Generate(r *bufio.Reader, pkgName string, outfile string) {
	parser := NewParser(r)
	def, rules, userCode := parser.Parse()

	regexs := make([]string, 0)
	actions := make([]string, 0)
	for _, v := range rules {
		regexs = append(regexs, v[0])
		actions = append(actions, v[1])
	}

	dfa := lexerDFA(regexs)
	stToID := make(map[automata.State]int)
	id := 1
	stToID[dfa.GetInitState()] = id
	id++
	for _, st := range dfa.GetStates() {
		if st == dfa.GetInitState() {
			continue
		}
		stToID[st] = id
		id++
	}
	idToSt := make([]automata.State, id)
	idToRegexID := make([]automata.RegexID, id)
	for st, id := range stToID {
		idToSt[id] = st
		idToRegexID[id] = dfa.GetRegexID(st)
	}

	embeddedTmpl := def
	stateIDToRegexIDTmpl := genStIdToRegexID(idToRegexID)
	finStatesTmpl := genFinStates(idToSt, dfa.GetFinStates())
	transitionTableTmpl := genTransitionTable(stToID, idToSt, dfa.GetTransitionTable())
	regexActionsTmpl := genRegexActions(actions)
	userCodeTmpl := userCode

	lexCfg := LexerTemplate{
		PackageName:          pkgName,
		EmbeddedTmpl:         embeddedTmpl,
		StateIDToRegexIDTmpl: stateIDToRegexIDTmpl,
		FinStatesTmpl:        finStatesTmpl,
		TransitionTableTmpl:  transitionTableTmpl,
		RegexActionsTmpl:     regexActionsTmpl,
		UserCodeTmpl:         userCodeTmpl,
	}
	s := tmpl
	t := template.Must(template.New("lexer").Parse(s))

	var buf bytes.Buffer
	if err := t.Execute(&buf, lexCfg); err != nil {
		panic(err)
	}

	// outfile := "tlex.yy.go"
	f, err := os.OpenFile(outfile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	data, err := imports.Process(outfile, buf.Bytes(), nil)
	if err != nil {
		panic(err)
	}
	io.Copy(f, bytes.NewReader(data))
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func genStIdToRegexID(idToRegexID []automata.RegexID) string {
	var buf bytes.Buffer
	for _, rid := range idToRegexID[1:] {
		buf.WriteString(fmt.Sprintf("%v,\n", rid))
	}

	return buf.String()
}

func genFinStates(idToSt []automata.State, finStates *collection.Set[automata.State]) string {
	var buf bytes.Buffer
	for i, st := range idToSt {
		if finStates.Contains(st) {
			buf.WriteString(fmt.Sprintf("%v: {},\n", i))
		}
	}

	return buf.String()
}

func genTransitionTable(stToID map[automata.State]int, idToSt []automata.State, delta *automata.DFATransition) string {
	tbl := make(map[int]map[byte]int)
	var buf bytes.Buffer
	iter := delta.Iterator()
	for iter.HasNext() {
		pair, to := iter.Next()
		from := pair.First
		b := pair.Second
		if _, ok := tbl[stToID[from]]; !ok {
			tbl[stToID[from]] = make(map[byte]int)
		}
		tbl[stToID[from]][b] = stToID[to]
	}

	for fromID := 1; fromID <= len(stToID); fromID++ {
		if _, ok := tbl[fromID]; ok {
			buf.WriteString(fmt.Sprintf("%v: {\n", fromID))
			for _, b := range automata.SupportedChars {
				if toID, ok2 := tbl[fromID][b]; ok2 {
					buf.WriteString(fmt.Sprintf("%v: %v,\n", b, toID))
				}
			}
			buf.WriteString("},\n")
		}
	}

	return buf.String()
}

func lexerNFA(regexs []string) automata.NFA {
	nfas := make([]*automata.NFA, 0)
	for i, regex := range regexs {
		nfa := parse(regex)
		(&nfa).SetRegexID(automata.RegexID(i + 1))
		nfas = append(nfas, &nfa)
	}

	nfa := *nfas[0]
	for _, n := range nfas[1:] {
		nfa = nfa.SumWithRegexID(*n)
	}

	return nfa
}

func lexerDFA(regexs []string) automata.DFA {
	nfa := lexerNFA(regexs)

	return nfa.ToImNFA().ToDFA().LexerMinimize().RemoveBH()
}

func genRegexActions(actions []string) string {

	var buf bytes.Buffer
	for i, v := range actions {
		buf.WriteString(fmt.Sprintf("case %v:\n", i+1))
		buf.WriteString(v + "\n")
		buf.WriteString("goto yystart\n")
	}

	return buf.String()
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
