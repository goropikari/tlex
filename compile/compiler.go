package compile

import "github.com/goropikari/golex/automata"

func Compile(regexp string) automata.NFA {
	lex := NewLexer(regexp)
	tokens := lex.Scan()
	parser := NewParser(tokens)
	ast, _ := parser.Parse()
	gen := NewCodeGenerator()
	ast.Accept(gen)

	return gen.GetNFA()
}
