package codegen

import (
	"fmt"
	"strings"

	"github.com/MistyPigeon/lingo/pkg/parser"
)

type CodeGen struct {
	output    strings.Builder
	indent    int
	imports   map[string]bool
	nullSafe  bool
}

func New() *CodeGen {
	return &CodeGen{
		imports: make(map[string]bool),
	}
}

func (cg *CodeGen) Generate(program *parser.Program) (string, error) {
	cg.output.Reset()
	cg.imports["fmt"] = false

	cg.emitln("package main")
	cg.emitln("")

	for _, item := range program.Items {
		switch node := item.(type) {
		case *parser.PackageDecl:
			// Skip - already emitted
		case *parser.ImportDecl:
			cg.imports[node.Path] = true
		case *parser.FuncDecl:
			cg. generateFunc(node)
		case *parser.VarDecl:
			cg.generateVar(node)
		case *parser.ConstDecl:
			cg. generateConst(node)
		case *parser.TypeDecl:
			cg. generateType(node)
		}
	}

	// Emit imports
	if len(cg.imports) > 0 {
		importStr := cg.generateImports()
		return importStr + cg.output.String(), nil
	}

	return cg. output.String(), nil
}

func (cg *CodeGen) generateImports() string {
	if len(cg.imports) == 0 {
		return ""
	}

	var imports strings.Builder
	imports.WriteString("import (\n")
	for path := range cg.imports {
		if path != "" {
			imports.WriteString(fmt.Sprintf(`	"%s"` + "\n", path))
		}
	}
	imports. WriteString(")\n\n")
	return imports.String()
}

func (cg *CodeGen) generateFunc(fn *parser.FuncDecl) {
	cg.emit("func ")

	if fn. Receiver != nil {
		cg.emit("(r " + fn.Receiver. Type + ") ")
	}

	cg.emit(fn.Name + "(")

	for i, param := range fn.Params {
		if i > 0 {
			cg.emit(", ")
		}
		cg.emit(param.Name + " " + param.Type)
	}

	cg.emit(")")

	if len(fn.Returns) > 0 {
		if len(fn.Returns) == 1 {
			cg.emit(" " + fn.Returns[0])
		} else {
			cg.emit(" (")
			for i, ret := range fn.Returns {
				if i > 0 {
					cg.emit(", ")
				}
				cg.emit(ret)
			}
			cg.emit(")")
		}
	}

	cg.emitln(" {")
	cg.indent++

	for _, stmt := range fn.Body {
		cg.generateStatement(stmt)
	}

	cg.indent--
	cg. emitln("}")
	cg.emitln("")
}

func (cg *CodeGen) generateVar(v *parser.VarDecl) {
	cg.emit("var " + v.Name)

	if v.Type != "" {
		cg. emit(" " + v.Type)
	}

	if v. Value != nil {
		cg. emit(" = ")
		cg.generateExpr(v.Value)
	}

	cg.emitln("")
}

func (cg *CodeGen) generateConst(c *parser.ConstDecl) {
	cg.emit("const " + c.Name)

	if c.Type != "" {
		cg. emit(" " + c.Type)
	}

	cg.emit(" = ")
	cg.generateExpr(c.Value)
	cg.emitln("")
}

func (cg *CodeGen) generateType(t *parser.TypeDecl) {
	cg.emit("type " + t.Name + " " + t.Type)
	cg.emitln("")
}

func (cg *CodeGen) generateStatement(stmt interface{}) {
	switch s := stmt.(type) {
	case *parser.VarDecl:
		cg.generateVar(s)
	case *parser. ConstDecl:
		cg.generateConst(s)
	case *parser.ReturnStmt:
		cg.generateReturn(s)
	case *parser. IfStmt:
		cg.generateIf(s)
	case *parser.ForStmt:
		cg. generateFor(s)
	case *parser.AssignStmt:
		cg.generateAssign(s)
	case *parser.ShortAssignStmt:
		cg.generateShortAssign(s)
	case *parser.CallExpr:
		cg. emit(cg.getIndent())
		cg.generateExpr(s)
		cg.emitln("")
	case *parser.MethodCall:
		cg.emit(cg.getIndent())
		cg.generateExpr(s)
		cg.emitln("")
	case *parser.DeferStmt:
		cg. emit(cg.getIndent() + "defer ")
		cg.generateExpr(s. Call)
		cg.emitln("")
	case *parser. GoStmt:
		cg. emit(cg.getIndent() + "go ")
		cg.generateExpr(s. Call)
		cg.emitln("")
	case *parser.PanicStmt:
		cg.emit(cg.getIndent() + "panic(")
		cg.generateExpr(s. Expr)
		cg.emitln(")")
	}
}

func (cg *CodeGen) generateReturn(ret *parser.ReturnStmt) {
	cg.emit(cg.getIndent() + "return")

	for i, val := range ret.Values {
		if i == 0 {
			cg.emit(" ")
		} else {
			cg.emit(", ")
		}
		cg.generateExpr(val)
	}

	cg.emitln("")
}

func (cg *CodeGen) generateIf(ifStmt *parser.IfStmt) {
	cg. emit(cg.getIndent() + "if ")
	cg.generateExpr(ifStmt.Condition)
	cg.emitln(" {")
	cg.indent++

	for _, stmt := range ifStmt.Then {
		cg.generateStatement(stmt)
	}

	cg.indent--
	cg.emit(cg.getIndent())

	if len(ifStmt. Else) > 0 {
		cg.emitln("} else {")
		cg.indent++

		for _, stmt := range ifStmt.Else {
			cg.generateStatement(stmt)
		}

		cg.indent--
		cg.emitln(cg.getIndent() + "}")
	} else {
		cg.emitln("}")
	}
}

func (cg *CodeGen) generateFor(forStmt *parser.ForStmt) {
	cg.emitln(cg.getIndent() + "for {")
	cg.indent++

	for _, stmt := range forStmt.Body {
		cg.generateStatement(stmt)
	}

	cg.indent--
	cg.emitln(cg.getIndent() + "}")
}

func (cg *CodeGen) generateAssign(assign *parser.AssignStmt) {
	cg.emit(cg.getIndent() + assign.Name + " = ")
	cg.generateExpr(assign.Value)
	cg.emitln("")
}

func (cg *CodeGen) generateShortAssign(assign *parser.ShortAssignStmt) {
	cg.emit(cg.getIndent() + assign.Name + " := ")
	cg.generateExpr(assign.Value)
	cg.emitln("")
}

func (cg *CodeGen) generateExpr(expr interface{}) {
	switch e := expr.(type) {
	case *parser.LiteralInt:
		cg.emit(e.Value)
	case *parser.LiteralFloat:
		cg.emit(e. Value)
	case *parser. LiteralString:
		cg.emit(`"` + e.Value + `"`)
	case *parser. LiteralBool:
		if e.Value {
			cg.emit("true")
		} else {
			cg.emit("false")
		}
	case *parser.LiteralNull:
		cg.emit("nil")
	case *parser. Identifier:
		cg.emit(e.Name)
	case *parser.BinaryOp:
		cg.generateBinaryOp(e)
	case *parser.UnaryOp:
		cg.generateUnaryOp(e)
	case *parser.CallExpr:
		cg. emit(e. Func + "(")
		for i, arg := range e.Args {
			if i > 0 {
				cg.emit(", ")
			}
			cg.generateExpr(arg)
		}
		cg.emit(")")
	case *parser.MethodCall:
		cg.emit(e.Receiver + "." + e.Method + "(")
		for i, arg := range e.Args {
			if i > 0 {
				cg.emit(", ")
			}
			cg.generateExpr(arg)
		}
		cg.emit(")")
	case *parser.IndexExpr:
		cg. generateExpr(e. Expr)
		cg.emit("[")
		cg.generateExpr(e.Index)
		cg.emit("]")
	case *parser.NullCheckExpr:
		cg. generateNullCheck(e)
	case *parser.ArrayLiteral:
		cg.emit("[]" + e.Type + "{")
		for i, elem := range e.Elements {
			if i > 0 {
				cg.emit(", ")
			}
			cg.generateExpr(elem)
		}
		cg.emit("}")
	case *parser.MapLiteral:
		cg.emit("map[string]interface{}{")
		first := true
		for key, val := range e.Pairs {
			if !first {
				cg.emit(", ")
			}
			cg.emit(`"` + key + `": `)
			cg.generateExpr(val)
			first = false
		}
		cg.emit("}")
	}
}

func (cg *CodeGen) generateBinaryOp(expr *parser.BinaryOp) {
	cg. emit("(")
	cg.generateExpr(expr.Left)
	cg.emit(" " + expr.Op + " ")
	cg.generateExpr(expr.Right)
	cg.emit(")")
}

func (cg *CodeGen) generateUnaryOp(expr *parser.UnaryOp) {
	cg.emit(expr.Op)
	cg.generateExpr(expr.Right)
}

func (cg *CodeGen) generateNullCheck(expr *parser.NullCheckExpr) {
	// Generate a nil check with default value
	cg.emit("func() interface{} { if ")
	cg.generateExpr(expr.Expr)
	cg.emit(" == nil { return ")
	cg.generateExpr(expr.DefaultExpr)
	cg.emit(" }; return ")
	cg.generateExpr(expr.Expr)
	cg.emit(" }()")
}

func (cg *CodeGen) emit(s string) {
	cg. output. WriteString(s)
}

func (cg *CodeGen) emitln(s string) {
	cg.output.WriteString(s)
	cg. output.WriteString("\n")
}

func (cg *CodeGen) getIndent() string {
	return strings.Repeat("\t", cg.indent)
}
