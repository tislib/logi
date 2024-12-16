package common

// Expression represents different types of expressions, with a Kind field to specify the type.
type Expression struct {
	Kind       ExpressionKind    `json:"kind"`
	Literal    *Literal          `json:"literal,omitempty"`
	Variable   *Variable         `json:"variable,omitempty"`
	BinaryExpr *BinaryExpression `json:"binaryExpr,omitempty"`
	FuncCall   *FunctionCall     `json:"funcCall,omitempty"`
}

func (e Expression) AsValue() Value {
	switch e.Kind {
	case LiteralKind:
		return e.Literal.Value
	case VariableKind:
		return StringValue(e.Variable.Name)
	case BinaryExprKind:
		left := e.BinaryExpr.Left.AsValue()
		right := e.BinaryExpr.Right.AsValue()
		return MapValue(map[string]Value{
			"left":     left,
			"operator": StringValue(e.BinaryExpr.Operator),
			"right":    right,
		})
	case FuncCallKind:
		args := make([]Value, len(e.FuncCall.Arguments))
		for i, arg := range e.FuncCall.Arguments {
			args[i] = arg.AsValue()
		}
		return MapValue(map[string]Value{
			"name":      StringValue(e.FuncCall.Name),
			"arguments": ArrayValue(args...),
		})
	default:
		panic("unknown expression kind")
	}
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
