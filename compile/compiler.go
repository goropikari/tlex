package compile

func Compile(regexp string) NFA {
	lex := NewLexer(regexp)
	tokens := lex.Scan()
	parser := NewParser(tokens)
	ast, _ := parser.Parse()
	gen := NewCodeGenerator()
	ast.Accept(gen)

	return gen.GetNFA()
}
