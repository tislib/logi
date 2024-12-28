package macro

import "github.com/tislib/logi/pkg/parser/lexer"

type NodeOp string

const (
	NodeOpFile                         = "file"
	NodeOpMacro                        = "macro"
	NodeOpSignature                    = "signature"
	NodeOpName                         = "name"
	NodeOpValueIdentifier              = "value_identifier"
	NodeOpValueNumber                  = "value_number"
	NodeOpValueString                  = "value_string"
	NodeOpValueBool                    = "value_bool"
	NodeOpValueArray                   = "value_array"
	NodeOpValueArrayItem               = "value_array_item"
	NodeOpBody                         = "body"
	NodeOpKind                         = "kind"
	NodeOpSyntax                       = "syntax"
	NodeOpTypeDef                      = "type_def"
	NodeOpTypes                        = "types"
	NodeOpSyntaxStatement              = "syntax_statement"
	NodeOpSyntaxExample                = "syntax_example"
	NodeOpSyntaxExamples               = "syntax_examples"
	NodeOpSyntaxElements               = "syntax_elements"
	NodeOpSyntaxKeywordElement         = "syntax_keyword_element"
	NodeOpSyntaxVariableKeywordElement = "syntax_variable_keyword_element"
	NodeOpTypesStatement               = "types_statement"
	NodeOpSyntaxParameterListElement   = "syntax_parameter_list_element"
	NodeOpSyntaxArgumentListElement    = "syntax_argument_list_element"
	NodeOpSyntaxAttributeListElement   = "syntax_attribute_list_element"
	NodeOpSyntaxCombinationElement     = "syntax_combination_element"
	NodeOpSyntaxTypeReferenceElement   = "syntax_type_reference_element"
	NodeOpScopes                       = "scopes"
	NodeOpScopesItem                   = "scopes_item"
	NodeOpSyntaxScopeElement           = "syntax_scope_element"
	NodeOpSyntaxSymbolElement          = "syntax_symbol_element"
)

var emptyToken = lexer.Token{}
var emptyLocation = lexer.Location{}

type yaccNode struct {
	op       NodeOp
	children []yaccNode
	value    interface{}
	location lexer.Location
	token    lexer.Token
}

func appendNode(nodeOp NodeOp, children ...yaccNode) yaccNode {
	return yaccNode{op: nodeOp, children: children}
}
func appendNodeX(nodeOp NodeOp, children ...yaccNode) yaccNode {
	return yaccNode{op: nodeOp, children: children}
}

func newNode(nodeOp NodeOp, value interface{}, token lexer.Token, location lexer.Location, children ...yaccNode) yaccNode {
	return yaccNode{op: nodeOp, value: value, children: children, token: token, location: location}
}

func appendNodeTo(node *yaccNode, child yaccNode) yaccNode {
	node.children = append(node.children, child)

	return *node
}

func registerRootNode(parser yyLexer, n yaccNode) {
	parser.(*yyMakroLexerProxy).Node.children = append(parser.(*yyMakroLexerProxy).Node.children, n)
}

func assertEqual(parser yyLexer, a, b interface{}, msg string) {
	if a != b {
		parser.Error(msg)
	}
}
