package parser

import (
	"logi/pkg/ast"
	"strings"
)

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

type yaccMacroNode struct {
	op       NodeOp
	children []yaccMacroNode
	value    interface{}
}

type yyMakroLexerProxy struct {
	lexer *macroLexer
	Node  yaccMacroNode
}

func (y *yyMakroLexerProxy) Lex(lval *yySymType) int {
	return y.lexer.Lex(lval)
}

func (y *yyMakroLexerProxy) Error(s string) {
	y.lexer.Error(s)
}

func ParseMacroContent(d string) (*ast.MacroAst, error) {
	s := newMacroLexer(strings.NewReader(d), true)
	parser := yyNewParser()
	proxy := &yyMakroLexerProxy{lexer: s, Node: yaccMacroNode{op: NodeOpFile}}

	parser.Parse(proxy)

	if s.Err != nil {
		return nil, s.Err
	}

	return convertNodeToMacroAst(proxy.Node)
}

func appendNode(nodeOp NodeOp, children ...yaccMacroNode) yaccMacroNode {
	return yaccMacroNode{op: nodeOp, children: children}
}

func newNode(nodeOp NodeOp, value interface{}, children ...yaccMacroNode) yaccMacroNode {
	return yaccMacroNode{op: nodeOp, value: value, children: children}
}

func appendNodeTo(node *yaccMacroNode, child yaccMacroNode) yaccMacroNode {
	node.children = append(node.children, child)

	return *node
}

func registerRootNode(parser yyLexer, n yaccMacroNode) {
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
