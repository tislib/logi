package logi

import (
	"logi/pkg/ast"
	"strings"
)

type yyLogiLexerProxy struct {
	lexer *logiLexer
	Node  yaccNode
}

func (y *yyLogiLexerProxy) Lex(lval *yySymType) int {
	return y.lexer.Lex(lval)
}

func (y *yyLogiLexerProxy) Error(s string) {
	y.lexer.Error(s)
}

func ParseLogiPlainContent(d string) (*ast.LogiPlainAst, error) {
	s := newLogiLexer(strings.NewReader(d), true)
	parser := yyNewParser()
	proxy := &yyLogiLexerProxy{lexer: s, Node: yaccNode{op: NodeOpFile}}

	parser.Parse(proxy)

	if s.Err != nil {
		return nil, s.Err
	}

	return convertNodeToLogiAst(proxy.Node)
}
