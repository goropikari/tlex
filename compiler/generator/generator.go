package generator

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
	// parse lexer configuration
	parser := NewParser(r)
	def, rules, userCode := parser.Parse()

	// compile regex and generate DFA
	regexs := make([]string, 0)
	actions := make([]string, 0)
	for _, v := range rules {
		regexs = append(regexs, v[0])
		actions = append(actions, v[1])
	}
	dfa := lexerDFA(regexs)
	oldstIDToNewStID := make(map[automata.StateID]automata.StateID)
	id := automata.StateID(1) // state id = 0 is reserved for dead state.
	oldstIDToNewStID[dfa.GetInitState()] = automata.StateID(id)
	id++
	for _, st := range dfa.GetStates() {
		if st == dfa.GetInitState() {
			continue
		}
		oldstIDToNewStID[st] = id
		id++
	}
	idToRegexID := make([]automata.RegexID, id)
	newStIDToOldStID := make([]automata.StateID, id)
	for oldid, newid := range oldstIDToNewStID {
		idToRegexID[newid] = dfa.GetRegexID(oldid)
		newStIDToOldStID[newid] = oldid
	}

	// generate lexer file
	embeddedTmpl := def
	stateIDToRegexIDTmpl := genStIdToRegexID(idToRegexID)
	finStatesTmpl := genFinStates(newStIDToOldStID, dfa.GetFinStates())
	transitionTableTmpl := genTransitionTable(oldstIDToNewStID, newStIDToOldStID, dfa.GetTransitionTable())
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
	// t.Execute(os.Stdout, lexCfg)

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

func genFinStates(newStIDToOldStID []automata.StateID, finStates *collection.Set[automata.StateID]) string {
	var buf bytes.Buffer
	for i, st := range newStIDToOldStID {
		if finStates.Contains(st) {
			buf.WriteString(fmt.Sprintf("%v: {},\n", i))
		}
	}

	return buf.String()
}

func genTransitionTable(oldIDToNewID map[automata.StateID]automata.StateID, newIDToOldID []automata.StateID, delta *automata.DFATransition) string {
	var buf bytes.Buffer
	for fromID := automata.StateID(1); fromID <= automata.StateID(len(oldIDToNewID)); fromID++ {
		mp, ok := delta.GetMap(newIDToOldID[fromID])
		if !ok {
			continue
		}
		buf.WriteString(fmt.Sprintf("%v: {\n", fromID))
		for intv, oldtoID := range mp {
			toID := oldIDToNewID[oldtoID]
			buf.WriteString(fmt.Sprintf("yyinterval{l: %v, r: %v}: %v,\n", intv.L, intv.R, toID))
		}
		buf.WriteString("},\n")
	}

	return buf.String()
}

func lexerNFA(regexs []string) *automata.NFA {
	nfas := make([]*automata.NFA, 0)
	for i, regex := range regexs {
		nfa := parse(regex)
		nfa.SetRegexID(automata.RegexID(i + 1))
		nfas = append(nfas, nfa)
	}

	nfa := nfas[0]
	for _, n := range nfas[1:] {
		nfa = nfa.Sum(n)
	}

	return nfa
}

func lexerDFA(regexs []string) *automata.DFA {
	nfa := lexerNFA(regexs)

	return nfa.ToImdNFA().ToDFA().LexerMinimize()
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

func parse(regex string) *automata.NFA {
	lex := regexp.NewLexer(regex)
	tokens := lex.Scan()
	parser := regexp.NewParser(tokens)
	ast, _ := parser.Parse()
	gen := regexp.NewCodeGenerator()
	ast.Accept(gen)

	return gen.GetNFA()
}
