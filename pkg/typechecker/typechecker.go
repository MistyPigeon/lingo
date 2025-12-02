package typechecker

import (
	"fmt"

	"github.com/MistyPigeon/lingo/pkg/parser"
)

type TypeChecker struct {
	scopes       []map[string]string
	nullableVars map[string]bool
}

func New() *TypeChecker {
	return &TypeChecker{
		scopes:       []map[string]string{make(map[string]string)},
		nullableVars: make(map[string]bool),
	}
}

func (tc *TypeChecker) Check(program *parser.Program) error {
	for _, item := range program.Items {
		switch node := item.(type) {
		case *parser.FuncDecl:
			if err := tc.checkFunc(node); err != nil {
				return err
			}
		case *parser.VarDecl:
			if err := tc.checkVar(node); err != nil {
				return err
			}
		case *parser.ConstDecl:
			if err := tc.checkConst(node); err != nil {
				return err
			}
		case *parser.TypeDecl:
			if err := tc.checkType(node); err != nil {
				return err
			}
		}
	}
	return nil
}

func (tc *TypeChecker) checkFunc(fn *parser.FuncDecl) error {
	tc.pushScope()
	defer tc.popScope()

	for _, param := range fn.Params {
		tc.defineVar(param.Name, param.Type)
	}

	for _, stmt := range fn.Body {
		if err := tc.checkStatement(stmt); err != nil {
			return err
		}
	}

	return nil
}

func (tc *TypeChecker) checkVar(v *parser.VarDecl) error {
	if v.Value != nil {
		exprType, err := tc.inferExprType(v.Value)
		if err != nil {
			return err
		}

		if v.Type != "" && v.Type != exprType && !tc.isCompatible(v.Type, exprType) {
			return fmt. Errorf("type mismatch for var %s: expected %s, got %s", v.Name, v.Type, exprType)
		}

		if v.IsNullable {
			tc.nullableVars[v.Name] = true
		}

		tc.defineVar(v.Name, exprType)
	} else if v.Type != "" {
		if v.IsNullable {
			tc.nullableVars[v.Name] = true
		}
		tc.defineVar(v. Name, v.Type)
	}

	return nil
}

func (tc *TypeChecker) checkConst(c *parser. ConstDecl) error {
	exprType, err := tc.inferExprType(c.Value)
	if err != nil {
		return err
	}

	if c.Type != "" && c.Type != exprType && !tc.isCompatible(c.Type, exprType) {
		return fmt. Errorf("type mismatch for const %s: expected %s, got %s", c.Name, c.Type, exprType)
	}

	tc.defineVar(c.Name, exprType)
	return nil
}

func (tc *TypeChecker) checkType(t *parser.TypeDecl) error {
	if t.IsNullable {
		tc. nullableVars[t.Name] = true
	}
	tc.defineVar(t. Name, t.Type)
	return nil
}

func (tc *TypeChecker) checkStatement(stmt interface{}) error {
	switch s := stmt.(type) {
	case *parser.VarDecl:
		return tc.checkVar(s)
	case *parser.ConstDecl:
		return tc.checkConst(s)
	case *parser.ReturnStmt:
		return tc.checkReturn(s)
	case *parser. IfStmt:
		return tc.checkIf(s)
	case *parser.ForStmt:
		return tc. checkFor(s)
	case *parser.CallExpr:
		_, err := tc.inferExprType(s)
		return err
	case *parser.AssignStmt:
		return tc.checkAssign(s)
	}
	return nil
}

func (tc *TypeChecker) checkReturn(ret *parser.ReturnStmt) error {
	for _, val := range ret.Values {
		_, err := tc.inferExprType(val)
		if err != nil {
			return err
		}
	}
	return nil
}

func (tc *TypeChecker) checkIf(ifStmt *parser.IfStmt) error {
	_, err := tc.inferExprType(ifStmt.Condition)
	if err != nil {
		return err
	}

	for _, stmt := range ifStmt.Then {
		if err := tc. checkStatement(stmt); err != nil {
			return err
		}
	}

	for _, stmt := range ifStmt. Else {
		if err := tc.checkStatement(stmt); err != nil {
			return err
		}
	}

	return nil
}

func (tc *TypeChecker) checkFor(forStmt *parser.ForStmt) error {
	tc.pushScope()
	defer tc.popScope()

	if forStmt.Init != nil {
		if err := tc.checkStatement(forStmt.Init); err != nil {
			return err
		}
	}

	if forStmt.Condition != nil {
		_, err := tc. inferExprType(forStmt.Condition)
		if err != nil {
			return err
		}
	}

	for _, stmt := range forStmt.Body {
		if err := tc.checkStatement(stmt); err != nil {
			return err
		}
	}

	return nil
}

func (tc *TypeChecker) checkAssign(assign *parser.AssignStmt) error {
	varType := tc.lookupVar(assign.Name)
	if varType == "" {
		return fmt.Errorf("undefined variable: %s", assign.Name)
	}

	exprType, err := tc.inferExprType(assign.Value)
	if err != nil {
		return err
	}

	if ! tc.isCompatible(varType, exprType) {
		return fmt.Errorf("cannot assign %s to %s", exprType, varType)
	}

	return nil
}

func (tc *TypeChecker) inferExprType(expr interface{}) (string, error) {
	switch e := expr.(type) {
	case *parser.LiteralInt:
		return "int", nil
	case *parser.LiteralFloat:
		return "float64", nil
	case *parser. LiteralString:
		return "string", nil
	case *parser.LiteralBool:
		return "bool", nil
	case *parser.LiteralNull:
		return "nil", nil
	case *parser. Identifier:
		varType := tc.lookupVar(e.Name)
		if varType == "" {
			return "", fmt.Errorf("undefined variable: %s", e.Name)
		}
		return varType, nil
	case *parser. BinaryOp:
		return tc.inferBinaryOpType(e)
	case *parser.UnaryOp:
		return tc.inferUnaryOpType(e)
	case *parser.CallExpr:
		return "interface{}", nil
	case *parser. IndexExpr:
		return "interface{}", nil
	case *parser.NullCheckExpr:
		exprType, err := tc.inferExprType(e.Expr)
		if err != nil {
			return "", err
		}
		if ! tc.nullableVars[exprType] {
			return "", fmt.Errorf("cannot use null coalescing on non-nullable type: %s", exprType)
		}
		return exprType, nil
	case *parser.NullableExpr:
		return tc.inferExprType(e.Expr)
	case *parser.ArrayLiteral:
		if e.Type != "" {
			return "[]" + e.Type, nil
		}
		return "[]interface{}", nil
	case *parser.MapLiteral:
		return "map[string]interface{}", nil
	default:
		return "interface{}", nil
	}
}

func (tc *TypeChecker) inferBinaryOpType(expr *parser.BinaryOp) (string, error) {
	leftType, err := tc.inferExprType(expr.Left)
	if err != nil {
		return "", err
	}

	rightType, err := tc.inferExprType(expr.Right)
	if err != nil {
		return "", err
	}

	if expr.Op == "+" || expr.Op == "-" || expr.Op == "*" || expr. Op == "/" {
		if leftType != rightType {
			return "", fmt. Errorf("type mismatch in binary operation: %s %s %s", leftType, expr.Op, rightType)
		}
		return leftType, nil
	}

	if expr.Op == "==" || expr.Op == "! =" || expr.Op == "<" || expr.Op == "<=" || expr. Op == ">" || expr.Op == ">=" {
		return "bool", nil
	}

	if expr.Op == "&&" || expr.Op == "||" {
		if leftType != "bool" || rightType != "bool" {
			return "", fmt.Errorf("logical operator requires bool operands")
		}
		return "bool", nil
	}

	return "interface{}", nil
}

func (tc *TypeChecker) inferUnaryOpType(expr *parser.UnaryOp) (string, error) {
	operandType, err := tc.inferExprType(expr.Right)
	if err != nil {
		return "", err
	}

	if expr.Op == "!" {
		if operandType != "bool" {
			return "", fmt.Errorf("logical not requires bool operand")
		}
		return "bool", nil
	}

	if expr.Op == "-" || expr.Op == "+" {
		if operandType != "int" && operandType != "float64" {
			return "", fmt.Errorf("unary %s requires numeric operand", expr.Op)
		}
		return operandType, nil
	}

	return operandType, nil
}

func (tc *TypeChecker) isCompatible(targetType, sourceType string) bool {
	if targetType == sourceType {
		return true
	}
	if targetType == "interface{}" {
		return true
	}
	if sourceType == "nil" {
		return true
	}
	return false
}

func (tc *TypeChecker) defineVar(name, varType string) {
	if len(tc.scopes) > 0 {
		tc. scopes[len(tc.scopes)-1][name] = varType
	}
}

func (tc *TypeChecker) lookupVar(name string) string {
	for i := len(tc.scopes) - 1; i >= 0; i-- {
		if t, ok := tc.scopes[i][name]; ok {
			return t
		}
	}
	return ""
}

func (tc *TypeChecker) pushScope() {
	tc.scopes = append(tc.scopes, make(map[string]string))
}

func (tc *TypeChecker) popScope() {
	if len(tc.scopes) > 1 {
		tc.scopes = tc. scopes[:len(tc.scopes)-1]
	}
}
