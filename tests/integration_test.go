package lingo

import (
	"os"
	"testing"

	"github.com/MistyPigeon/lingo/pkg/codegen"
	"github.com/MistyPigeon/lingo/pkg/lexer"
	"github.com/MistyPigeon/lingo/pkg/parser"
	"github.com/MistyPigeon/lingo/pkg/typechecker"
)

func TestBasicLexing(t *testing. T) {
	source := `package main
func main() {
	var x: int = 42
}`

	lex := lexer.New(source)
	tokens := lex. Tokenize()

	if len(tokens) == 0 {
		t. Fatal("No tokens generated")
	}

	if tokens[0]. Type != lexer.TOKEN_PACKAGE {
		t. Errorf("Expected PACKAGE token, got %v", tokens[0].Type)
	}
}

func TestBasicParsing(t *testing.T) {
	source := `package main
func main() {
	var x: int = 42
}`

	lex := lexer. New(source)
	tokens := lex.Tokenize()

	p := parser.New(tokens)
	ast, err := p.Parse()

	if err != nil
