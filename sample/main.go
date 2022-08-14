package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"sample/lex"
)

func main() {
	b, _ := io.ReadAll(os.Stdin)
	data := string(b)
	lexer := lex.NewLexer(data)

	fmt.Printf("output:\n")
	for {
		tok, err := lexer.NextToken()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
		}

		fmt.Printf("\t%v\n", tok)
	}
	fmt.Println(lex.NumChars, lex.NumLines)
}
