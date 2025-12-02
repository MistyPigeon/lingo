package parser

type ASTNode interface {
	astNode()
}

type Program struct {
	Items []ASTNode
}

func (p *Program) astNode() {}

type PackageDecl struct {
	Name string
}

func (p *PackageDecl) astNode() {}

type ImportDecl struct {
	Path  string
	Alias string
}

func (i *ImportDecl) astNode() {}

type FuncDecl struct {
	Name    string
	Receiver *Param
	Params  []*Param
	Returns []string
	Body    []ASTNode
}

func (f *FuncDecl) astNode() {}

type Param struct {
	Name string
	Type string
}

type VarDecl struct {
	Name         string
	Type         string
	Value        ASTNode
	IsNullable   bool
	Initializer  ASTNode
}

func (v *VarDecl) astNode() {}

type ConstDecl struct {
	Name  string
	Type  string
	Value ASTNode
}

func (c *ConstDecl) astNode() {}

type StructDecl struct {
	Name   string
	Fields []*StructField
}

func (s *StructDecl) astNode() {}

type StructField struct {
	Name       string
	Type       string
	IsNullable bool
	Tag        string
}

type InterfaceDecl struct {
	Name    string
	Methods []*Param
}

func (i *InterfaceDecl) astNode() {}

type TypeDecl struct {
	Name  string
	Type  string
	IsNullable bool
}

func (t *TypeDecl) astNode() {}

type ReturnStmt struct {
	Values []ASTNode
}

func (r *ReturnStmt) astNode() {}

type IfStmt struct {
	Condition ASTNode
	Then      []ASTNode
	Else      []ASTNode
}

func (i *IfStmt) astNode() {}

type ForStmt struct {
	Init      ASTNode
	Condition ASTNode
	Post      ASTNode
	Body      []ASTNode
}

func (f *ForStmt) astNode() {}

type ForRangeStmt struct {
	Key   string
	Value string
	Expr  ASTNode
	Body  []ASTNode
}

func (f *ForRangeStmt) astNode() {}

type AssignStmt struct {
	Name  string
	Value ASTNode
}

func (a *AssignStmt) astNode() {}

type ShortAssignStmt struct {
	Name  string
	Value ASTNode
}

func (s *ShortAssignStmt) astNode() {}

type CallExpr struct {
	Func string
	Args []ASTNode
}

func (c *CallExpr) astNode() {}

type MethodCall struct {
	Receiver string
	Method   string
	Args     []ASTNode
}

func (m *MethodCall) astNode() {}

type BinaryOp struct {
	Left  ASTNode
	Op    string
	Right ASTNode
}

func (b *BinaryOp) astNode() {}

type UnaryOp struct {
	Op    string
	Right ASTNode
}

func (u *UnaryOp) astNode() {}

type LiteralInt struct {
	Value string
}

func (l *LiteralInt) astNode() {}

type LiteralFloat struct {
	Value string
}

func (l *LiteralFloat) astNode() {}

type LiteralString struct {
	Value string
}

func (l *LiteralString) astNode() {}

type LiteralBool struct {
	Value bool
}

func (l *LiteralBool) astNode() {}

type LiteralNull struct{}

func (l *LiteralNull) astNode() {}

type Identifier struct {
	Name string
}

func (i *Identifier) astNode() {}

type NullableExpr struct {
	Expr ASTNode
}

func (n *NullableExpr) astNode() {}

type NullCheckExpr struct {
	Expr        ASTNode
	DefaultExpr ASTNode
}

func (n *NullCheckExpr) astNode() {}

type IndexExpr struct {
	Expr  ASTNode
	Index ASTNode
}

func (i *IndexExpr) astNode() {}

type SliceExpr struct {
	Expr  ASTNode
	Start ASTNode
	End   ASTNode
}

func (s *SliceExpr) astNode() {}

type MapLiteral struct {
	KeyType   string
	ValueType string
	Pairs     map[string]ASTNode
}

func (m *MapLiteral) astNode() {}

type ArrayLiteral struct {
	Type     string
	Elements []ASTNode
}

func (a *ArrayLiteral) astNode() {}

type StructLiteral struct {
	Type   string
	Fields map[string]ASTNode
}

func (s *StructLiteral) astNode() {}

type ChanOp struct {
	Op    string
	Expr  ASTNode
	Value ASTNode
}

func (c *ChanOp) astNode() {}

type GoStmt struct {
	Call *CallExpr
}

func (g *GoStmt) astNode() {}

type SelectStmt struct {
	Cases []*SelectCase
}

func (s *SelectStmt) astNode() {}

type SelectCase struct {
	ChanOp *ChanOp
	Body   []ASTNode
}

type DeferStmt struct {
	Call *CallExpr
}

func (d *DeferStmt) astNode() {}

type PanicStmt struct {
	Expr ASTNode
}

func (p *PanicStmt) astNode() {}

type RecoverExpr struct{}

func (r *RecoverExpr) astNode() {}
