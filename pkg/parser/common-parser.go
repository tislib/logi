package parser

type NodeOp string

const (
	NodeOpFile                         = "file"
	NodeOpMacro                        = "macro"
	NodeOpSignature                    = "signature"
	NodeOpName                         = "name"
	NodeOpValue                        = "value"
	NodeOpBody                         = "body"
	NodeOpKind                         = "kind"
	NodeOpSyntax                       = "syntax"
	NodeOpTypeDef                      = "type_def"
	NodeOpDefinition                   = "definition"
	NodeOpSyntaxStatement              = "syntax_statement"
	NodeOpSyntaxKeywordElement         = "syntax_keyword_element"
	NodeOpSyntaxVariableKeywordElement = "syntax_variable_keyword_element"
	NodeOpSyntaxDefinitionElement      = "syntax_definition_element"
	NodeOpSyntaxParameterListElement   = "syntax_parameter_list_element"
	NodeOpSyntaxArgumentListElement    = "syntax_argument_list_element"
	NodeOpSyntaxCodeBlockElement       = "syntax_code_block_element"
	NodeOpSyntaxAttributeListElement   = "syntax_attribute_list_element"
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
	parser.(*yyMakroLexerProxy).Node.children = append(parser.(*yyMakroLexerProxy).Node.children, n)
}

func assertEqual(parser yyLexer, a, b interface{}, msg string) {
	if a != b {
		parser.Error(msg)
	}
}

func init() {
	yyErrorVerbose = true
	yyDebug = 5
}
