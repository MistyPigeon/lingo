package lexer

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type TokenType string

const (
	// Literals
	TOKEN_INT    TokenType = "INT"
	TOKEN_FLOAT  TokenType = "FLOAT"
	TOKEN_STRING TokenType = "STRING"
	TOKEN_BOOL   TokenType = "BOOL"
	TOKEN_NULL   TokenType = "NULL"

	// Keywords
	TOKEN_FUNC    TokenType = "FUNC"
	TOKEN_VAR     TokenType = "VAR"
	TOKEN_CONST   TokenType = "CONST"
	TOKEN_TYPE    TokenType = "TYPE"
	TOKEN_STRUCT  TokenType = "STRUCT"
	TOKEN_RETURN  TokenType = "RETURN"
	TOKEN_IF      TokenType = "IF"
	TOKEN_ELSE    TokenType = "ELSE"
	TOKEN_FOR     TokenType = "FOR"
	TOKEN_PACKAGE TokenType = "PACKAGE"
	TOKEN_IMPORT  TokenType = "IMPORT"
	TOKEN_NULL_CHECK TokenType = "?"
	TOKEN_INTERFACE TokenType = "INTERFACE"
	TOKEN_CHAN    TokenType = "CHAN"
	TOKEN_GO      TokenType = "GO"
	TOKEN_SELECT  TokenType = "SELECT"
	TOKEN_CASE    TokenType = "CASE"
	TOKEN_DEFAULT TokenType = "DEFAULT"
	TOKEN_DEFER   TokenType = "DEFER"
	TOKEN_PANIC   TokenType = "PANIC"
	TOKEN_RECOVER TokenType = "RECOVER"

	// Identifiers
	TOKEN_IDENT TokenType = "IDENT"

	// Operators
	TOKEN_ASSIGN   TokenType = "="
	TOKEN_PLUS     TokenType = "+"
	TOKEN_MINUS    TokenType = "-"
	TOKEN_MUL      TokenType = "*"
	TOKEN_DIV      TokenType = "/"
	TOKEN_MOD      TokenType = "%"
	TOKEN_EQ       TokenType = "=="
	TOKEN_NEQ      TokenType = "!="
	TOKEN_LT       TokenType = "<"
	TOKEN_LTE      TokenType = "<="
	TOKEN_GT       TokenType = ">"
	TOKEN_GTE      TokenType = ">="
	TOKEN_LAND     TokenType = "&&"
	TOKEN_LOR      TokenType = "||"
	TOKEN_LNOT     TokenType = "!"
	TOKEN_AND      TokenType = "&"
	TOKEN_OR       TokenType = "|"
	TOKEN_XOR      TokenType = "^"
	TOKEN_LSHIFT   TokenType = "<<"
	TOKEN_RSHIFT   TokenType = ">>"
	TOKEN_WALRUS   TokenType = ":="
	TOKEN_DOT      TokenType = "."
	TOKEN_COMMA    TokenType = ","
	TOKEN_COLON    TokenType = ":"
	TOKEN_SEMICOLON TokenType = ";"
	TOKEN_ARROW    TokenType = "->"
	TOKEN_LPAREN   TokenType = "("
	TOKEN_RPAREN   TokenType = ")"
	TOKEN_LBRACE   TokenType = "{"
	TOKEN_RBRACE   TokenType = "}"
	TOKEN_LBRACKET TokenType = "["
	TOKEN_RBRACKET TokenType = "]"
	TOKEN_QUESTION TokenType = "?"

	// Special
	TOKEN_EOF   TokenType = "EOF"
	TOKEN_NEWLINE TokenType = "NEWLINE"
)

type Token struct {
	Type  TokenType
	Value string
	Line  int
	Col   int
}

type Lexer struct {
	input  string
	pos    int
	line   int
	col    int
	tokens []Token
}

func New(input string) *Lexer {
	return &Lexer{
		input:  input,
		pos:    0,
		line:   1,
		col:    1,
		tokens: []Token{},
	}
}

func (l *Lexer) Tokenize() []Token {
	for l.pos < len(l. input) {
		l.skipWhitespaceAndComments()

		if l.pos >= len(l.input) {
			break
		}

		ch := l.current()

		if unicode.IsLetter(rune(ch)) || ch == '_' {
			l.readIdentifierOrKeyword()
		} else if unicode.IsDigit(rune(ch)) {
			l.readNumber()
		} else if ch == '"' {
			l.readString()
		} else if ch == '\'' {
			l.readChar()
		} else {
			l.readOperator()
		}
	}

	l.tokens = append(l.tokens, Token{Type: TOKEN_EOF, Value: "", Line: l.line, Col: l.col})
	return l. tokens
}

func (l *Lexer) current() byte {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}

func (l *Lexer) peek(n int) byte {
	pos := l.pos + n
	if pos >= len(l.input) {
		return 0
	}
	return l.input[pos]
}

func (l *Lexer) advance() {
	if l.pos < len(l.input) && l.input[l.pos] == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col++
	}
	l.pos++
}

func (l *Lexer) addToken(typ TokenType, value string) {
	l.tokens = append(l.tokens, Token{
		Type:  typ,
		Value: value,
		Line:  l.line,
		Col:   l.col,
	})
}

func (l *Lexer) skipWhitespaceAndComments() {
	for l.pos < len(l.input) {
		ch := l.current()

		if ch == ' ' || ch == '\t' || ch == '\r' {
			l.advance()
		} else if ch == '\n' {
			l.advance()
		} else if ch == '/' && l.peek(1) == '/' {
			for l.pos < len(l.input) && l.current() != '\n' {
				l. advance()
			}
		} else if ch == '/' && l.peek(1) == '*' {
			l.advance()
			l.advance()
			for l.pos < len(l.input) {
				if l.current() == '*' && l.peek(1) == '/' {
					l.advance()
					l.advance()
					break
				}
				l.advance()
			}
		} else {
			break
		}
	}
}

func (l *Lexer) readIdentifierOrKeyword() {
	start := l.pos
	startCol := l.col

	for l.pos < len(l. input) && (unicode.IsLetter(rune(l.current())) || unicode.IsDigit(rune(l.current())) || l.current() == '_') {
		l.advance()
	}

	value := l.input[start:l.pos]

	// Check if it's a keyword
	typ := TOKEN_IDENT
	switch value {
	case "func":
		typ = TOKEN_FUNC
	case "var":
		typ = TOKEN_VAR
	case "const":
		typ = TOKEN_CONST
	case "type":
		typ = TOKEN_TYPE
	case "struct":
		typ = TOKEN_STRUCT
	case "return":
		typ = TOKEN_RETURN
	case "if":
		typ = TOKEN_IF
	case "else":
		typ = TOKEN_ELSE
	case "for":
		typ = TOKEN_FOR
	case "package":
		typ = TOKEN_PACKAGE
	case "import":
		typ = TOKEN_IMPORT
	case "interface":
		typ = TOKEN_INTERFACE
	case "chan":
		typ = TOKEN_CHAN
	case "go":
		typ = TOKEN_GO
	case "select":
		typ = TOKEN_SELECT
	case "case":
		typ = TOKEN_CASE
	case "default":
		typ = TOKEN_DEFAULT
	case "defer":
		typ = TOKEN_DEFER
	case "panic":
		typ = TOKEN_PANIC
	case "recover":
		typ = TOKEN_RECOVER
	case "true", "false":
		typ = TOKEN_BOOL
	case "null":
		typ = TOKEN_NULL
	}

	l.tokens = append(l.tokens, Token{
		Type:  typ,
		Value: value,
		Line:  l.line,
		Col:   startCol,
	})
}

func (l *Lexer) readNumber() {
	start := l.pos
	startCol := l.col

	for l.pos < len(l.input) && unicode.IsDigit(rune(l.current())) {
		l.advance()
	}

	if l.current() == '.' && unicode.IsDigit(rune(l.peek(1))) {
		l.advance()
		for l.pos < len(l. input) && unicode.IsDigit(rune(l.current())) {
			l.advance()
		}
		l.tokens = append(l.tokens, Token{
			Type:  TOKEN_FLOAT,
			Value: l.input[start:l.pos],
			Line:  l. line,
			Col:   startCol,
		})
	} else {
		l.tokens = append(l.tokens, Token{
			Type:  TOKEN_INT,
			Value: l.input[start:l.pos],
			Line:  l.line,
			Col:   startCol,
		})
	}
}

func (l *Lexer) readString() {
	startCol := l.col
	l.advance() // Skip opening quote
	start := l.pos

	for l.pos < len(l.input) && l.current() != '"' {
		if l.current() == '\\' {
			l.advance()
		}
		l.advance()
	}

	value := l.input[start:l.pos]
	l.advance() // Skip closing quote

	l.tokens = append(l.tokens, Token{
		Type:  TOKEN_STRING,
		Value: value,
		Line:  l.line,
		Col:   startCol,
	})
}

func (l *Lexer) readChar() {
	startCol := l.col
	l.advance() // Skip opening quote
	start := l.pos

	for l.pos < len(l.input) && l.current() != '\'' {
		if l.current() == '\\' {
			l.advance()
		}
		l.advance()
	}

	value := l. input[start:l.pos]
	l.advance() // Skip closing quote

	l.tokens = append(l.tokens, Token{
		Type:  TOKEN_STRING,
		Value: value,
		Line:  l.line,
		Col:   startCol,
	})
}

func (l *Lexer) readOperator() {
	startCol := l.col
	ch := l.current()

	// Two-character operators
	if l.pos+1 < len(l.input) {
		twoChar := string([]byte{ch, l.peek(1)})
		switch twoChar {
		case "==":
			l.advance()
			l.advance()
			l.tokens = append(l. tokens, Token{Type: TOKEN_EQ, Value: "==", Line: l.line, Col: startCol})
			return
		case "!=":
			l.advance()
			l.advance()
			l.tokens = append(l.tokens, Token{Type: TOKEN_NEQ, Value: "!=", Line: l.line, Col: startCol})
			return
		case "<=":
			l.advance()
			l.advance()
			l.tokens = append(l.tokens, Token{Type: TOKEN_LTE, Value: "<=", Line: l.line, Col: startCol})
			return
		case ">=":
			l.advance()
			l.advance()
			l.tokens = append(l.tokens, Token{Type: TOKEN_GTE, Value: ">=", Line: l.line, Col: startCol})
			return
		case "&&":
			l.advance()
			l.advance()
			l.tokens = append(l.tokens, Token{Type: TOKEN_LAND, Value: "&&", Line: l.line, Col: startCol})
			return
		case "||":
			l.advance()
			l.advance()
			l.tokens = append(l.tokens, Token{Type: TOKEN_LOR, Value: "||", Line: l.line, Col: startCol})
			return
		case ":=":
			l.advance()
			l.advance()
			l.tokens = append(l.tokens, Token{Type: TOKEN_WALRUS, Value: ":=", Line: l.line, Col: startCol})
			return
		case "->":
			l.advance()
			l.advance()
			l.tokens = append(l.tokens, Token{Type: TOKEN_ARROW, Value: "->", Line: l.line, Col: startCol})
			return
		case "<<":
			l.advance()
			l.advance()
			l.tokens = append(l.tokens, Token{Type: TOKEN_LSHIFT, Value: "<<", Line: l.line, Col: startCol})
			return
		case ">>":
			l. advance()
			l.advance()
			l.tokens = append(l.tokens, Token{Type: TOKEN_RSHIFT, Value: ">>", Line: l.line, Col: startCol})
			return
		}
	}

	// Single-character operators
	switch ch {
	case '+':
		l. advance()
		l.tokens = append(l.tokens, Token{Type: TOKEN_PLUS, Value: "+", Line: l.line, Col: startCol})
	case '-':
		l.advance()
		l.tokens = append(l.tokens, Token{Type: TOKEN_MINUS, Value: "-", Line: l.line, Col: startCol})
	case '*':
		l.advance()
		l.tokens = append(l.tokens, Token{Type: TOKEN_MUL, Value: "*", Line: l.line, Col: startCol})
	case '/':
		l.advance()
		l.tokens = append(l.tokens, Token{Type: TOKEN_DIV, Value: "/", Line: l.line, Col: startCol})
	case '%':
		l.advance()
		l. tokens = append(l.tokens, Token{Type: TOKEN_MOD, Value: "%", Line: l.line, Col: startCol})
	case '=':
		l.advance()
		l. tokens = append(l.tokens, Token{Type: TOKEN_ASSIGN, Value: "=", Line: l.line, Col: startCol})
	case '<':
		l.advance()
		l. tokens = append(l.tokens, Token{Type: TOKEN_LT, Value: "<", Line: l.line, Col: startCol})
	case '>':
		l.advance()
		l. tokens = append(l.tokens, Token{Type: TOKEN_GT, Value: ">", Line: l.line, Col: startCol})
	case '!':
		l.advance()
		l. tokens = append(l.tokens, Token{Type: TOKEN_LNOT, Value: "!", Line: l.line, Col: startCol})
	case '&':
		l.advance()
		l.tokens = append(l.tokens, Token{Type: TOKEN_AND, Value: "&", Line: l.line, Col: startCol})
	case '|':
		l.advance()
		l. tokens = append(l.tokens, Token{Type: TOKEN_OR, Value: "|", Line: l.line, Col: startCol})
	case '^':
		l.advance()
		l. tokens = append(l.tokens, Token{Type: TOKEN_XOR, Value: "^", Line: l.line, Col: startCol})
	case '. ':
		l.advance()
		l.tokens = append(l.tokens, Token{Type: TOKEN_DOT, Value: ".", Line: l.line, Col: startCol})
	case ',':
		l.advance()
		l.tokens = append(l.tokens, Token{Type: TOKEN_COMMA, Value: ",", Line: l.line, Col: startCol})
	case ':':
		l.advance()
		l.tokens = append(l.tokens, Token{Type: TOKEN_COLON, Value: ":", Line: l.line, Col: startCol})
	case ';':
		l.advance()
		l.tokens = append(l.tokens, Token{Type: TOKEN_SEMICOLON, Value: ";", Line: l.line, Col: startCol})
	case '(':
		l.advance()
		l.tokens = append(l.tokens, Token{Type: TOKEN_LPAREN, Value: "(", Line: l.line, Col: startCol})
	case ')':
		l.advance()
		l.tokens = append(l.tokens, Token{Type: TOKEN_RPAREN, Value: ")", Line: l.line, Col: startCol})
	case '{':
		l.advance()
		l.tokens = append(l.tokens, Token{Type: TOKEN_LBRACE, Value: "{", Line: l.line, Col: startCol})
	case '}':
		l.advance()
		l.tokens = append(l.tokens, Token{Type: TOKEN_RBRACE, Value: "}", Line: l.line, Col: startCol})
	case '[':
		l.advance()
		l.tokens = append(l.tokens, Token{Type: TOKEN_LBRACKET, Value: "[", Line: l.line, Col: startCol})
	case ']':
		l.advance()
		l. tokens = append(l.tokens, Token{Type: TOKEN_RBRACKET, Value: "]", Line: l.line, Col: startCol})
	case '? ':
		l.advance()
		l.tokens = append(l.tokens, Token{Type: TOKEN_QUESTION, Value: "? ", Line: l.line, Col: startCol})
	default:
		l.advance()
	}
}
