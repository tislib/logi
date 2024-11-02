package parser

import (
	"logi/pkg/ast"
	"strings"
)

type NodeOp string

const (
	NodeOpFile      = "file"
	NodeOpMacro     = "macro"
	NodeOpSignature = "signature"
	NodeOpName      = "name"
	NodeOpBody      = "body"
	NodeOpKind      = "kind"
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

func newNode(nodeOp NodeOp, value interface{}) yaccMacroNode {
	return yaccMacroNode{op: nodeOp, value: value}
}

func registerRootNode(parser yyLexer, n yaccMacroNode) {
	children := parser.(*yyMakroLexerProxy).Node.children
	parser.(*yyMakroLexerProxy).Node.children = append([]yaccMacroNode{n}, children...)
}

func assertEqual(parser yyLexer, a, b interface{}, msg string) {
	if a != b {
		parser.Error(msg)
	}
}

func init() {
	yyErrorVerbose = true
	//yyDebug = 5
}
