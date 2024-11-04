package logi

type NodeOp string

const (
	NodeOpFile                = "file"
	NodeOpSignature           = "signature"
	NodeOpMacro               = "macro"
	NodeOpName                = "name"
	NodeOpValue               = "value"
	NodeOpBody                = "body"
	NodeOpDefinition          = "definition"
	NodeOpStatements          = "statements"
	NodeOpStatement           = "statement"
	NodeOpIdentifier          = "identifier"
	NodeOpAttributeList       = "attributeList"
	NodeOpAttribute           = "attribute"
	NodeOpArgumentList        = "argumentList"
	NodeOpArgument            = "argument"
	NodeOpTypeDef             = "type_def"
	NodeOpExpression          = "expression"
	NodeOpAssignment          = "assignment"
	NodeOpIf                  = "if"
	NodeOpIfElse              = "if_else"
	NodeOpReturn              = "return"
	NodeOpVariableDeclaration = "var_declaration"
	NodeOpLiteral             = "literal"
	NodeOpVariable            = "variable"
	NodeOpBinaryExpression    = "binary_expression"
	NodeOpFunctionCall        = "function_call"
	NodeOpFunctionParams      = "function_params"
	NodeOpCodeBlock           = "code_block"
	NodeOpOperator            = "operator"
	NodeOpFunctionDefinition  = "function_definition"
	NodeOpArray               = "array"
	NodeOpParameterList       = "parameter_list"
	NodeOpParameter           = "parameter"
)

type yaccNode struct {
	op       NodeOp
	children []yaccNode
	value    interface{}
}

func appendNode(nodeOp NodeOp, children ...yaccNode) yaccNode {
	return yaccNode{op: nodeOp, children: children}
}

func newNode(nodeOp NodeOp, value interface{}, children ...yaccNode) yaccNode {
	return yaccNode{op: nodeOp, value: value, children: children}
}

func appendNodeTo(node *yaccNode, child yaccNode) yaccNode {
	node.children = append(node.children, child)

	return *node
}

func registerRootNode(parser yyLexer, n yaccNode) {
	parser.(*yyLogiLexerProxy).Node.children = append(parser.(*yyLogiLexerProxy).Node.children, n)
}

func assertEqual(parser yyLexer, a, b interface{}, msg string) {
	if a != b {
		parser.Error(msg)
	}
}

func init() {
	yyErrorVerbose = true
}
