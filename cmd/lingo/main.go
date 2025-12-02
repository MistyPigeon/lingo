package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/MistyPigeon/lingo/pkg/codegen"
	"github.com/MistyPigeon/lingo/pkg/lexer"
	"github.com/MistyPigeon/lingo/pkg/parser"
	"github.com/MistyPigeon/lingo/pkg/typechecker"
)

func main() {
	var (
		inputFile  = flag.String("file", "", "Input . lingo file to compile")
		outputFile = flag. String("out", "", "Output . go file (default: input filename with .go extension)")
		checkOnly  = flag.Bool("check", false, "Only perform type checking without generating code")
		verbose    = flag.Bool("v", false, "Verbose output")
	)

	flag.Parse()

	if *inputFile == "" {
		fmt.Fprintf(os.Stderr, "Usage: lingo -file <input.lingo> [-out <output.go>] [-check] [-v]\n")
		os.Exit(1)
	}

	// Read input file
	source, err := os.ReadFile(*inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Determine output file
	outFile := *outputFile
	if outFile == "" {
		outFile = filepath.Join(
			filepath.Dir(*inputFile),
			filepath.Base(*inputFile[:len(*inputFile)-len(filepath. Ext(*inputFile))])+".go",
		)
	}

	// Lexing
	lex := lexer.New(string(source))
	tokens := lex.Tokenize()

	if *verbose {
		fmt.Println("=== TOKENS ===")
		for _, tok := range tokens {
			fmt. Printf("%v: %q\n", tok.Type, tok.Value)
		}
		fmt.Println()
	}

	// Parsing
	p := parser.New(tokens)
	ast, err := p.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
		os.Exit(1)
	}

	if *verbose {
		fmt. Println("=== AST ===")
		printAST(ast, 0)
		fmt.Println()
	}

	// Type checking
	tc := typechecker.New()
	err = tc.Check(ast)
	if err != nil {
		fmt.Fprintf(os. Stderr, "Type error: %v\n", err)
		os.Exit(1)
	}

	if *checkOnly {
		fmt.Println("Type checking passed!")
		return
	}

	// Code generation
	gen := codegen.New()
	goCode, err := gen.Generate(ast)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Codegen error: %v\n", err)
		os.Exit(1)
	}

	// Write output file
	err = os.WriteFile(outFile, []byte(goCode), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully compiled %s -> %s\n", *inputFile, outFile)
}

func printAST(node interface{}, depth int) {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}

	switch n := node.(type) {
	case *parser.Program:
		fmt.Printf("%sProgram\n", indent)
		for _, item := range n.Items {
			printAST(item, depth+1)
		}
	case *parser.FuncDecl:
		fmt.Printf("%sFunction: %s\n", indent, n.Name)
		if n.Receiver != nil {
			fmt.Printf("%s  Receiver: %s\n", indent, n.Receiver)
		}
		fmt.Printf("%s  Params:\n", indent)
		for _, p := range n.Params {
			fmt.Printf("%s    %s: %s\n", indent, p.Name, p.Type)
		}
		fmt.Printf("%s  Returns:\n", indent)
		for _, r := range n.Returns {
			fmt.Printf("%s    %s\n", indent, r)
		}
		fmt.Printf("%s  Body:\n", indent)
		for _, stmt := range n.Body {
			printAST(stmt, depth+2)
		}
	case *parser.VarDecl:
		fmt. Printf("%sVar: %s: %s\n", indent, n.Name, n.Type)
	case *parser. StructDecl:
		fmt.Printf("%sStruct: %s\n", indent, n.Name)
		for _, f := range n.Fields {
			fmt.Printf("%s  %s: %s\n", indent, f.Name, f.Type)
		}
	}
}
