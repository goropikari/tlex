package regexp

import "github.com/goropikari/golex/automata"

func Compile(regexp string) automata.DFA {
	lex := NewLexer(regexp)
	tokens := lex.Scan()
	parser := NewParser(tokens)
	ast, _ := parser.Parse()
	gen := NewCodeGenerator()
	ast.Accept(gen)
	dfa := gen.GetNFA().ToDFA().LexerMinimize()

	return dfa
}
