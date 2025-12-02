package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/MistyPigeon/lingo/pkg/lexer"
)

type Parser struct {
	tokens    []lexer.Token
	pos       int
	current   lexer.Token
	peekToken lexer.Token
}

func New(tokens []lexer.Token) *Parser {
	p := &Parser{
		tokens: tokens,
		pos:    0,
	}
	p.advance()
	p.advance()
	return p
}

func (p *Parser) Parse() (*Program, error) {
	program := &Program{Items: []ASTNode{}}

	for ! p.is(lexer.TOKEN_EOF) {
		item, err := p.parseTopLevel()
		if err != nil {
			return nil, err
		}
		if item != nil {
			program.Items = append(program.Items, item)
		}
	}

	return program, nil
}

func (p *Parser) parseTopLevel() (ASTNode, error) {
	switch p.current.Type {
	case lexer.TOKEN_PACKAGE:
		return p.parsePackage()
	case lexer. TOKEN_IMPORT:
		return p.parseImport()
	case lexer.TOKEN_FUNC:
		return p.parseFunc()
	case lexer. TOKEN_TYPE:
		return p.parseType()
	case lexer.TOKEN_VAR:
		return p.parseVar()
	case lexer. TOKEN_CONST:
		return p.parseConst()
	default:
		return nil, fmt. Errorf("unexpected token: %v", p.current. Type)
	}
}

func (p *Parser) parsePackage() (*PackageDecl, error) {
	if ! p.match(lexer.TOKEN_PACKAGE) {
		return nil, fmt.Errorf("expected package")
	}

	name := p.current.Value
	p. advance()

	return &PackageDecl{Name: name}, nil
}

func (p *Parser) parseImport() (*ImportDecl, error) {
	if !p.match(lexer.TOKEN_IMPORT) {
		return nil, fmt. Errorf("expected import")
	}

	var path, alias string

	if p.is(lexer.TOKEN_LPAREN) {
		p.advance()
		path = strings. Trim(p.current.Value, `"`)
		p.advance()
		p.expect(lexer.TOKEN_RPAREN)
	} else {
		if p.is(lexer.TOKEN_IDENT) {
			alias = p.current.Value
			p.advance()
		}
		path = strings.Trim(p.current.Value, `"`)
		p.advance()
	}

	return &ImportDecl{Path: path, Alias: alias}, nil
}

func (p *Parser) parseFunc() (*FuncDecl, error) {
	if !p.match(lexer.TOKEN_FUNC) {
		return nil, fmt.Errorf("expected func")
	}

	var receiver *Param

	if p.is(lexer.TOKEN_LPAREN) {
		p.advance()
		receiver = &Param{
			Name: p.current.Value,
			Type: "",
		}
		p.advance()
		if p.is(lexer.TOKEN_MUL) {
			receiver.Type = "*"
			p.advance()
		}
		receiver.Type += p.current.Value
		p. advance()
		p.expect(lexer.TOKEN_RPAREN)
	}

	name := p.current.Value
	p.advance()

	p. expect(lexer.TOKEN_LPAREN)
	params, err := p.parseParamList()
	if err != nil {
		return nil, err
	}
	p.expect(lexer.TOKEN_RPAREN)

	returns := []string{}
	if p.is(lexer.TOKEN_LPAREN) || p.is(lexer.TOKEN_IDENT) || p.is(lexer.TOKEN_MUL) {
		if p.is(lexer.TOKEN_LPAREN) {
			p. advance()
			for ! p.is(lexer.TOKEN_RPAREN) {
				ret := ""
				if p.is(lexer.TOKEN_MUL) {
					ret = "*"
					p.advance()
				}
				if p.is(lexer.TOKEN_LBRACKET) {
					p.advance()
					ret += "[]"
					p.advance()
				}
				ret += p.current.Value
				p.advance()
				returns = append(returns, ret)
				if p.is(lexer. TOKEN_COMMA) {
					p.advance()
				}
			}
			p.expect(lexer.TOKEN_RPAREN)
		} else {
			ret := ""
			if p.is(lexer.TOKEN_MUL) {
				ret = "*"
				p.advance()
			}
			if p.is(lexer.TOKEN_LBRACKET) {
				p.advance()
				ret += "[]"
				p.advance()
			}
			ret += p.current.Value
			p. advance()
			returns = append(returns, ret)
		}
	}

	p.expect(lexer.TOKEN_LBRACE)
	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}
	p.expect(lexer.TOKEN_RBRACE)

	return &FuncDecl{
		Name:     name,
		Receiver: receiver,
		Params:   params,
		Returns:  returns,
		Body:     body,
	}, nil
}

func (p *Parser) parseParamList() ([]*Param, error) {
	params := []*Param{}

	for !p.is(lexer. TOKEN_RPAREN) && ! p.is(lexer.TOKEN_EOF) {
		name := p.current.Value
		p.advance()

		p.expect(lexer.TOKEN_COLON)

		varType := ""
		if p. is(lexer.TOKEN_MUL) {
			varType = "*"
			p.advance()
		}
		if p.is(lexer.TOKEN_LBRACKET) {
			varType += "["
			p.advance()
			varType += "]"
			p.advance()
		}
		varType += p.current.Value
		p.advance()

		params = append(params, &Param{Name: name, Type: varType})

		if p.is(lexer.TOKEN_COMMA) {
			p.advance()
		}
	}

	return params, nil
}

func (p *Parser) parseType() (*TypeDecl, error) {
	if !p.match(lexer. TOKEN_TYPE) {
		return nil, fmt.Errorf("expected type")
	}

	name := p.current.Value
	p.advance()

	varType := ""
	isNullable := false

	if p.is(lexer.TOKEN_QUESTION) {
		isNullable = true
		p.advance()
	}

	if p.is(lexer.TOKEN_MUL) {
		varType = "*"
		p.advance()
	}
	if p.is(lexer.TOKEN_LBRACKET) {
		varType += "["
		p.advance()
		varType += "]"
		p. advance()
	}

	varType += p.current.Value
	p.advance()

	return &TypeDecl{Name: name, Type: varType, IsNullable: isNullable}, nil
}

func (p *Parser) parseVar() (*VarDecl, error) {
	if !p.match(lexer.TOKEN_VAR) {
		return nil, fmt. Errorf("expected var")
	}

	name := p.current.Value
	p.advance()

	var varType string
	isNullable := false

	if p.is(lexer.TOKEN_COLON) {
		p.advance()
		if p.is(lexer.TOKEN_QUESTION) {
			isNullable = true
			p.advance()
		}
		varType = p.parseTypeAnnotation()
	}

	var value ASTNode
	if p.is(lexer.TOKEN_ASSIGN) {
		p.advance()
		var err error
		value, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
	}

	return &VarDecl{
		Name:       name,
		Type:       varType,
		Value:      value,
		IsNullable: isNullable,
	}, nil
}

func (p *Parser) parseTypeAnnotation() string {
	varType := ""
	if p.is(lexer.TOKEN_MUL) {
		varType = "*"
		p.advance()
	}
	if p.is(lexer.TOKEN_LBRACKET) {
		varType += "["
		p.advance()
		varType += "]"
		p.advance()
	}
	if p.is(lexer.TOKEN_LBRACE) {
		p.advance()
		keyType := p.current.Value
		p.advance()
		p.expect(lexer.TOKEN_RBRACE)
		varType += "map[" + keyType + "]"
		if p.is(lexer.TOKEN_MUL) {
			varType += "*"
			p. advance()
		}
		varType += p.current.Value
		p. advance()
	} else {
		varType += p.current.Value
		p.advance()
	}
	return varType
}

func (p *Parser) parseConst() (*ConstDecl, error) {
	if !p.match(lexer.TOKEN_CONST) {
		return nil, fmt. Errorf("expected const")
	}

	name := p.current.Value
	p.advance()

	var varType string
	if p.is(lexer.TOKEN_COLON) {
		p.advance()
		varType = p.parseTypeAnnotation()
	}

	p.expect(lexer.TOKEN_ASSIGN)
	value, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	return &ConstDecl{Name: name, Type: varType, Value: value}, nil
}

func (p *Parser) parseBlock() ([]ASTNode, error) {
	statements := []ASTNode{}

	for !p.is(lexer.TOKEN_RBRACE) && !p.is(lexer. TOKEN_EOF) {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			statements = append(statements, stmt)
		}
	}

	return statements, nil
}

func (p *Parser) parseStatement() (ASTNode, error) {
	switch p.current.Type {
	case lexer.TOKEN_VAR:
		return p.parseVar()
	case lexer. TOKEN_CONST:
		return p.parseConst()
	case lexer.TOKEN_RETURN:
		return p.parseReturn()
	case lexer.TOKEN_IF:
		return p. parseIf()
	case lexer. TOKEN_FOR:
		return p.parseFor()
	case lexer.TOKEN_DEFER:
		return p.parseDefer()
	case lexer.TOKEN_GO:
		return p.parseGo()
	case lexer. TOKEN_SELECT:
		return p.parseSelect()
	case lexer.TOKEN_PANIC:
		return p.parsePanic()
	case lexer. TOKEN_IDENT:
		return p.parseAssignmentOrCall()
	default:
		return nil, fmt. Errorf("unexpected statement: %v", p.current. Type)
	}
}

func (p *Parser) parseReturn() (*ReturnStmt, error) {
	if !p.match(lexer. TOKEN_RETURN) {
		return nil, fmt.Errorf("expected return")
	}

	values := []ASTNode{}

	if !p.is(lexer.TOKEN_RBRACE) && !p.is(lexer.TOKEN_EOF) {
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		values = append(values, expr)

		for p.is(lexer.TOKEN_COMMA) {
			p.advance()
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			values = append(values, expr)
		}
	}

	return &ReturnStmt{Values: values}, nil
}

func (p *Parser) parseIf() (*IfStmt, error) {
	if !p.match(lexer. TOKEN_IF) {
		return nil, fmt.Errorf("expected if")
	}

	cond, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	p.expect(lexer.TOKEN_LBRACE)
	thenBlock, err := p.parseBlock()
	if err != nil {
		return nil, err
	}
	p.expect(lexer.TOKEN_RBRACE)

	elseBlock := []ASTNode{}
	if p.is(lexer.TOKEN_ELSE) {
		p.advance()
		if p.is(lexer.TOKEN_IF) {
			elseIf, err := p.parseIf()
			if err != nil {
				return nil, err
			}
			elseBlock = append(elseBlock, elseIf)
		} else {
			p.expect(lexer.TOKEN_LBRACE)
			elseBlock, err = p.parseBlock()
			if err != nil {
				return nil, err
			}
			p.expect(lexer.TOKEN_RBRACE)
		}
	}

	return &IfStmt{Condition: cond, Then: thenBlock, Else: elseBlock}, nil
}

func (p *Parser) parseFor() (ASTNode, error) {
	if !p.match(lexer.TOKEN_FOR) {
		return nil, fmt.Errorf("expected for")
	}

	p.expect(lexer.TOKEN_LBRACE)
	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}
	p.expect(lexer.TOKEN_RBRACE)

	return &ForStmt{Body: body}, nil
}

func (p *Parser) parseDefer() (*DeferStmt, error) {
	if !p.match(lexer.TOKEN_DEFER) {
		return nil, fmt. Errorf("expected defer")
	}

	if !p.is(lexer.TOKEN_IDENT) {
		return nil, fmt.Errorf("expected function call after defer")
	}

	funcName := p.current.Value
	p.advance()

	p.expect(lexer.TOKEN_LPAREN)
	args, err := p.parseArgList()
	if err != nil {
		return nil, err
	}
	p.expect(lexer.TOKEN_RPAREN)

	return &DeferStmt{Call: &CallExpr{Func: funcName, Args: args}}, nil
}

func (p *Parser) parseGo() (*GoStmt, error) {
	if !p.match(lexer.TOKEN_GO) {
		return nil, fmt.Errorf("expected go")
	}

	if !p.is(lexer.TOKEN_IDENT) {
		return nil, fmt. Errorf("expected function call after go")
	}

	funcName := p.current.Value
	p.advance()

	p.expect(lexer.TOKEN_LPAREN)
	args, err := p.parseArgList()
	if err != nil {
		return nil, err
	}
	p.expect(lexer.TOKEN_RPAREN)

	return &GoStmt{Call: &CallExpr{Func: funcName, Args: args}}, nil
}

func (p *Parser) parseSelect() (*SelectStmt, error) {
	if !p.match(lexer.TOKEN_SELECT) {
		return nil, fmt.Errorf("expected select")
	}

	p.expect(lexer.TOKEN_LBRACE)

	cases := []*SelectCase{}
	for ! p.is(lexer.TOKEN_RBRACE) && !p.is(lexer.TOKEN_EOF) {
		if !p.match(lexer.TOKEN_CASE) {
			break
		}

		var chanOp *ChanOp
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		chanOp = &ChanOp{Op: "<-", Expr: expr}

		p.expect(lexer.TOKEN_COLON)
		body, err := p.parseBlock()
		if err != nil {
			return nil, err
		}

		cases = append(cases, &SelectCase{ChanOp: chanOp, Body: body})
	}

	p.expect(lexer. TOKEN_RBRACE)

	return &SelectStmt{Cases: cases}, nil
}

func (p *Parser) parsePanic() (*PanicStmt, error) {
	if !p. match(lexer.TOKEN_PANIC) {
		return nil, fmt.Errorf("expected panic")
	}

	p.expect(lexer.TOKEN_LPAREN)
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	p.expect(lexer.TOKEN_RPAREN)

	return &PanicStmt{Expr: expr}, nil
}

func (p *Parser) parseAssignmentOrCall() (ASTNode, error) {
	name := p.current.Value
	p.advance()

	if p. is(lexer.TOKEN_ASSIGN) {
		p.advance()
		value, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		return &AssignStmt{Name: name, Value: value}, nil
	} else if p.is(lexer.TOKEN_WALRUS) {
		p.advance()
		value, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		return &ShortAssignStmt{Name: name, Value: value}, nil
	} else if p.is(lexer.TOKEN_LPAREN) {
		p.advance()
		args, err := p.parseArgList()
		if err != nil {
			return nil, err
		}
		p.expect(lexer.TOKEN_RPAREN)
		return &CallExpr{Func: name, Args: args}, nil
	} else if p.is(lexer.TOKEN_DOT) {
		p.advance()
		method := p.current.Value
		p.advance()
		p.expect(lexer.TOKEN_LPAREN)
		args, err := p.parseArgList()
		if err != nil {
			return nil, err
		}
		p.expect(lexer.TOKEN_RPAREN)
		return &MethodCall{Receiver: name, Method: method, Args: args}, nil
	}

	return &Identifier{Name: name}, nil
}

func (p *Parser) parseArgList() ([]ASTNode, error) {
	args := []ASTNode{}

	for !p.is(lexer. TOKEN_RPAREN) && ! p.is(lexer.TOKEN_EOF) {
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		args = append(args, expr)

		if p. is(lexer.TOKEN_COMMA) {
			p.advance()
		}
	}

	return args, nil
}

func (p *Parser) parseExpr() (ASTNode, error) {
	return p.parseLogicalOr()
}

func (p *Parser) parseLogicalOr() (ASTNode, error) {
	left, err := p.parseLogicalAnd()
	if err != nil {
		return nil, err
	}

	for p.is(lexer.TOKEN_LOR) {
		op := p.current.Value
		p.advance()
		right, err := p.parseLogicalAnd()
		if err != nil {
			return nil, err
		}
		left = &BinaryOp{Left: left, Op: op, Right: right}
	}

	return left, nil
}

func (p *Parser) parseLogicalAnd() (ASTNode, error) {
	left, err := p.parseEquality()
	if err != nil {
		return nil, err
	}

	for p.is(lexer.TOKEN_LAND) {
		op := p. current.Value
		p.advance()
		right, err := p.parseEquality()
		if err != nil {
			return nil, err
		}
		left = &BinaryOp{Left: left, Op: op, Right: right}
	}

	return left, nil
}

func (p *Parser) parseEquality() (ASTNode, error) {
	left, err := p.parseComparison()
	if err != nil {
		return nil, err
	}

	for p.is(lexer. TOKEN_EQ) || p.is(lexer.TOKEN_NEQ) {
		op := p.current. Value
		p.advance()
		right, err := p.parseComparison()
		if err != nil {
			return nil, err
		}
		left = &BinaryOp{Left: left, Op: op, Right: right}
	}

	return left, nil
}

func (p *Parser) parseComparison() (ASTNode, error) {
	left, err := p.parseBitwiseOr()
	if err != nil {
		return nil, err
	}

	for p.is(lexer.TOKEN_LT) || p.is(lexer.TOKEN_LTE) || p.is(lexer.TOKEN_GT) || p.is(lexer.TOKEN_GTE) {
		op := p.current.Value
		p.advance()
		right, err := p.parseBitwiseOr()
		if err != nil {
			return nil, err
		}
		left = &BinaryOp{Left: left, Op: op, Right: right}
	}

	return left, nil
}

func (p *Parser) parseBitwiseOr() (ASTNode, error) {
	left, err := p.parseBitwiseXor()
	if err != nil {
		return nil, err
	}

	for p.is(lexer.TOKEN_OR) && !p.peekIs(lexer.TOKEN_OR) {
		op := p.current.Value
		p. advance()
		right, err := p.parseBitwiseXor()
		if err != nil {
			return nil, err
		}
		left = &BinaryOp{Left: left, Op: op, Right: right}
	}

	return left, nil
}

func (p *Parser) parseBitwiseXor() (ASTNode, error) {
	left, err := p.parseBitwiseAnd()
	if err != nil {
		return nil, err
	}

	for p.is(lexer.TOKEN_XOR) {
		op := p. current.Value
		p.advance()
		right, err := p.parseBitwiseAnd()
		if err != nil {
			return nil, err
		}
		left = &BinaryOp{Left: left, Op: op, Right: right}
	}

	return left, nil
}

func (p *Parser) parseBitwiseAnd() (ASTNode, error) {
	left, err := p.parseShift()
	if err != nil {
		return nil, err
	}

	for p.is(lexer.TOKEN_AND) && !p.peekIs(lexer.TOKEN_AND) {
		op := p.current.Value
		p.advance()
		right, err := p.parseShift()
		if err != nil {
			return nil, err
		}
		left = &BinaryOp{Left: left, Op: op, Right: right}
	}

	return left, nil
}

func (p *Parser) parseShift() (ASTNode, error) {
	left, err := p.parseAdditive()
	if err != nil {
		return nil, err
	}

	for p.is(lexer.TOKEN_LSHIFT) || p.is(lexer.TOKEN_RSHIFT) {
		op := p. current.Value
		p.advance()
		right, err := p.parseAdditive()
		if err != nil {
			return nil, err
		}
		left = &BinaryOp{Left: left, Op: op, Right: right}
	}

	return left, nil
}

func (p *Parser) parseAdditive() (ASTNode, error) {
	left, err := p.parseMultiplicative()
	if err != nil {
		return nil, err
	}

	for p.is(lexer.TOKEN_PLUS) || p.is(lexer.TOKEN_MINUS) {
		op := p.current.Value
		p.advance()
		right, err := p.parseMultiplicative()
		if err != nil {
			return nil, err
		}
		left = &BinaryOp{Left: left, Op: op, Right: right}
	}

	return left, nil
}

func (p *Parser) parseMultiplicative() (ASTNode, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	for p. is(lexer.TOKEN_MUL) || p.is(lexer.TOKEN_DIV) || p.is(lexer. TOKEN_MOD) {
		op := p.current.Value
		p.advance()
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		left = &BinaryOp{Left: left, Op: op, Right: right}
	}

	return left, nil
}

func (p *Parser) parseUnary() (ASTNode, error) {
	if p.is(lexer.TOKEN_LNOT) || p.is(lexer.TOKEN_MINUS) || p.is(lexer.TOKEN_PLUS) || p.is(lexer.TOKEN_AND) || p.is(lexer.TOKEN_MUL) {
		op := p.current. Value
		p.advance()
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return &UnaryOp{Op: op, Right: right}, nil
	}

	return p.parsePostfix()
}

func (p *Parser) parsePostfix() (ASTNode, error) {
	left, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	for {
		if p.is(lexer.TOKEN_LBRACKET) {
			p.advance()
			index, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			p.expect(lexer.TOKEN_RBRACKET)
			left = &IndexExpr{Expr: left, Index: index}
		} else if p.is(lexer.TOKEN_DOT) {
			p.advance()
			field := p.current.Value
			p.advance()
			left = &BinaryOp{Left: left, Op: ".", Right: &Identifier{Name: field}}
		} else if p.is(lexer.TOKEN_QUESTION) {
			p.advance()
			if p.is(lexer.TOKEN_COLON) {
				p.advance()
				def, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				left = &NullCheckExpr{Expr: left, DefaultExpr: def}
			} else {
				left = &NullableExpr{Expr: left}
			}
		} else {
			break
		}
	}

	return left, nil
}

func (p *Parser) parsePrimary() (ASTNode, error) {
	switch p.current.Type {
	case lexer.TOKEN_INT:
		value := p.current.Value
		p.advance()
		return &LiteralInt{Value: value}, nil

	case lexer.TOKEN_FLOAT:
		value := p. current.Value
		p.advance()
		return &LiteralFloat{Value: value}, nil

	case lexer.TOKEN_STRING:
		value := p.current.Value
		p.advance()
		return &LiteralString{Value: value}, nil

	case lexer.TOKEN_BOOL:
		value := p.current.Value == "true"
		p.advance()
		return &LiteralBool{Value: value}, nil

	case lexer.TOKEN_NULL:
		p.advance()
		return &LiteralNull{}, nil

	case lexer. TOKEN_IDENT:
		name := p. current.Value
		p.advance()
		if p.is(lexer.TOKEN_LPAREN) {
			p.advance()
			args, err := p.parseArgList()
			if err != nil {
				return nil, err
			}
			p.expect(lexer.TOKEN_RPAREN)
			return &CallExpr{Func: name, Args: args}, nil
		}
		return &Identifier{Name: name}, nil

	case lexer.TOKEN_LPAREN:
		p.advance()
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		p.expect(lexer.TOKEN_RPAREN)
		return expr, nil

	case lexer.TOKEN_LBRACKET:
		return p.parseArrayOrSlice()

	case lexer.TOKEN_LBRACE:
		return p.parseMapOrStruct()

	case lexer. TOKEN_RECOVER:
		p.advance()
		return &RecoverExpr{}, nil

	default:
		return nil, fmt.Errorf("unexpected primary: %v", p.current. Type)
	}
}

func (p *Parser) parseArrayOrSlice() (ASTNode, error) {
	p. expect(lexer.TOKEN_LBRACKET)

	if p.is(lexer.TOKEN_RBRACKET) {
		p.advance()
		elemType := p.current.Value
		p.advance()
		return &ArrayLiteral{Type: elemType, Elements: []ASTNode{}}, nil
	}

	elements := []ASTNode{}
	for ! p.is(lexer.TOKEN_RBRACKET) {
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		elements = append(elements, expr)
		if p.is(lexer.TOKEN_COMMA) {
			p.advance()
		}
	}
	p.expect(lexer. TOKEN_RBRACKET)

	return &ArrayLiteral{Type: "", Elements: elements}, nil
}

func (p *Parser) parseMapOrStruct() (ASTNode, error) {
	p. expect(lexer.TOKEN_LBRACE)

	pairs := make(map[string]ASTNode)
	for !p.is(lexer.TOKEN_RBRACE) {
		key := p.current.Value
		p.advance()
		p.expect(lexer.TOKEN_COLON)
		value, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		pairs[key] = value
		if p.is(lexer. TOKEN_COMMA) {
			p.advance()
		}
	}
	p.expect(lexer.TOKEN_RBRACE)

	return &MapLiteral{Pairs: pairs}, nil
}

func (p *Parser) advance() {
	p.pos++
	if p.pos < len(p.tokens) {
		p.current = p.tokens[p.pos]
		if p.pos+1 < len(p.tokens) {
			p.peekToken = p.tokens[p.pos+1]
		}
	}
}

func (p *Parser) is(typ lexer.TokenType) bool {
	return p.current.Type == typ
}

func (p *Parser) peekIs(typ lexer.TokenType) bool {
	return p.peekToken.Type == typ
}

func (p *Parser) match(typ lexer.TokenType) bool {
	if p.is(typ) {
		p.advance()
		return true
	}
	return false
}

func (p *Parser) expect(typ lexer.TokenType) error {
	if !p.is(typ) {
		return fmt.Errorf("expected %v, got %v", typ, p.current.Type)
	}
	p.advance()
	return nil
}
