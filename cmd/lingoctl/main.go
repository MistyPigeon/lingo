package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/MistyPigeon/lingo/pkg/lexer"
	"github.com/MistyPigeon/lingo/pkg/parser"
)

func main() {
	var (
		command = flag.String("cmd", "", "Command: lex, parse")
		file    = flag.String("file", "", "Input file")
	)

	flag.Parse()

	if *command == "" || *file == "" {
		fmt.Fprintf(os.Stderr, "Usage: lingoctl -cmd <lex|parse> -file <file. lingo>\n")
		os.Exit(1)
	}

	source, err := os.ReadFile(*file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	switch *command {
	case "lex":
		lex := lexer.New(string(source))
		tokens := lex.Tokenize()
		for _, tok := range tokens {
			fmt. Printf("%v: %q\n", tok.Type, tok.Value)
		}

	case "parse":
		lex := lexer.New(string(source))
		tokens := lex.Tokenize()
		p := parser.New(tokens)
		ast, err := p.Parse()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
			os. Exit(1)
		}
		fmt.Printf("AST parsed successfully.  Items: %d\n", len(ast.Items))

	default:
		fmt. Fprintf(os.Stderr, "Unknown command: %s\n", *command)
		os.Exit(1)
	}
}
