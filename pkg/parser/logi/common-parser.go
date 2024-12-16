package logi

import "github.com/tislib/logi/pkg/parser/lexer"

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
	NodeOpLiteral             = "literal"
	NodeOpVariable            = "variable"
	NodeOpBinaryExpression    = "binary_expression"
	NodeOpFunctionCall        = "function_call"
	NodeOpFunctionParams      = "function_params"
	NodeOpCodeBlock           = "code_block"
	NodeOpOperator            = "operator"
	NodeOpArray               = "array"
	NodeOpParameterList       = "parameter_list"
	NodeOpStruct              = "struct"
	NodeOpJsonObject          = "json_object"
	NodeOpJsonObjectItem      = "json_object_item"
	NodeOpJsonObjectItemValue = "json_object_item_value"
	NodeOpJsonArray           = "json_array"
	NodeOpJsonIdentifier      = "json_identifier"
)

type yaccNode struct {
	op       NodeOp
	children []yaccNode
	value    interface{}
	token    lexer.Token
	location lexer.Location
}

func appendNode(nodeOp NodeOp, children ...yaccNode) yaccNode {
	result := yaccNode{op: nodeOp, children: children}

	if len(children) > 0 {
		result.location = children[0].location
	}

	return result
}

func newNode(nodeOp NodeOp, value interface{}, token lexer.Token, location lexer.Location, children ...yaccNode) yaccNode {
	return yaccNode{op: nodeOp, value: value, children: children, token: token, location: location}
}

func appendNodeTo(node *yaccNode, child yaccNode) yaccNode {
	node.children = append(node.children, child)

	return *node
}

func registerRootNode(parser yyLexer, n yaccNode) {
	parser.(*yyLogiLexerProxy).Node.children = append(parser.(*yyLogiLexerProxy).Node.children, n)
}

func init() {
	yyErrorVerbose = true
}
