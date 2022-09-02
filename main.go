package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/goropikari/golex/compiler/lex"
)

var (
	pkgName string
	srcfile string
	outfile string
)

func main() {
	// go get github.com/pkg/profile
	// go tool pprof -http=":8081" cpu.pprof
	// defer profile.Start(profile.ProfilePath(".")).Stop()
	flag.StringVar(&pkgName, "pkg", "main", "generated go file package name")
	flag.StringVar(&srcfile, "src", "", "input lexer configuration file: sample.l")
	flag.StringVar(&outfile, "o", "golex.yy.go", "generated file path")
	flag.Parse()
	if srcfile == "" {
		fmt.Fprint(os.Stderr, "srcfile is required.\n")
	}

	f, err := os.OpenFile(srcfile, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	r := bufio.NewReader(f)
	lex.Generate(r, pkgName, outfile)
}
