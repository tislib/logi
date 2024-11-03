package common

type CodeBlock struct {
	Statements []Statement `json:"statements"`
}

// Statement represents different types of statements, with a Kind field to specify the type.
type Statement struct {
	Kind StatementKind `json:"kind"`

	ExprStmt   *ExpressionStatement `json:"exprStmt,omitempty"`
	AssignStmt *AssignmentStatement `json:"assignStmt,omitempty"`
	IfStmt     *IfStatement         `json:"ifStmt,omitempty"`
	ReturnStmt *ReturnStatement     `json:"returnStmt,omitempty"`
	ForStmt    *ForStatement        `json:"forStmt,omitempty"`
	SwitchStmt *SwitchStatement     `json:"switchStmt,omitempty"`
	VarDecl    *VarDeclaration      `json:"varDecl,omitempty"`
}

// StatementKind is an enum-like type representing the kind of statement.
type StatementKind string

const (
	ExprStatementKind   StatementKind = "expression"
	AssignStatementKind StatementKind = "assignment"
	IfStatementKind     StatementKind = "if"
	ReturnStatementKind StatementKind = "return"
	ForStatementKind    StatementKind = "for"
	SwitchStatementKind StatementKind = "switch"
	VarDeclKind         StatementKind = "var_declaration"
)

// ExpressionStatement represents a standalone expression, like a function call.
type ExpressionStatement struct {
	Expr *Expression `json:"expr"`
}

// AssignmentStatement represents an assignment, e.g., `x = 5`.
type AssignmentStatement struct {
	Left  *Expression `json:"left"`
	Right *Expression `json:"right"`
}

// IfStatement represents an if statement with an optional else block.
type IfStatement struct {
	Condition *Expression `json:"condition"`
	ThenBlock *CodeBlock  `json:"thenBlock"`
	ElseBlock *CodeBlock  `json:"elseBlock,omitempty"`
}

// ReturnStatement represents a return statement, e.g., `return x`.
type ReturnStatement struct {
	Result *Expression `json:"result"`
}

// ForStatement represents a for-loop.
type ForStatement struct {
	Init      *Statement  `json:"init,omitempty"`
	Condition *Expression `json:"condition,omitempty"`
	Post      *Statement  `json:"post,omitempty"`
	Body      *CodeBlock  `json:"body"`
}

// SwitchStatement represents a switch-case structure.
type SwitchStatement struct {
	Expression *Expression      `json:"expression"`
	Cases      []*CaseStatement `json:"cases"`
	Default    *CodeBlock       `json:"default,omitempty"`
}

// CaseStatement represents each case in a switch.
type CaseStatement struct {
	Condition *Expression `json:"condition,omitempty"`
	Body      *CodeBlock  `json:"body"`
}

// VarDeclaration represents a variable declaration, e.g., `var x int = 5`.
type VarDeclaration struct {
	Name  string         `json:"name"`
	Type  TypeDefinition `json:"type,omitempty"`
	Value *Expression    `json:"value,omitempty"`
}

// Expressions

// Expression represents different types of expressions, with a Kind field to specify the type.
type Expression struct {
	Kind       ExpressionKind    `json:"kind"`
	Literal    *Literal          `json:"literal,omitempty"`
	Variable   *Variable         `json:"variable,omitempty"`
	BinaryExpr *BinaryExpression `json:"binaryExpr,omitempty"`
	FuncCall   *FunctionCall     `json:"funcCall,omitempty"`
}

// ExpressionKind is an enum-like type representing the kind of expression.
type ExpressionKind string

const (
	LiteralKind    ExpressionKind = "literal"
	VariableKind   ExpressionKind = "variable"
	BinaryExprKind ExpressionKind = "binary_expression"
	FuncCallKind   ExpressionKind = "function_call"
)

// Literal represents basic values like integers, strings, etc.
type Literal struct {
	Value Value `json:"value"`
}

// Variable represents a variable, e.g., `x`.
type Variable struct {
	Name string `json:"name"`
}

// BinaryExpression represents binary operations, e.g., `x + y`.
type BinaryExpression struct {
	Left     *Expression `json:"left"`
	Operator string      `json:"operator"` // e.g., "+", "-", "*", "/"
	Right    *Expression `json:"right"`
}

// FunctionCall represents a function call, e.g., `foo(x, y)`.
type FunctionCall struct {
	Name      string        `json:"name"`
	Arguments []*Expression `json:"arguments"`
}
