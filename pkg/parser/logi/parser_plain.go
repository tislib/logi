package logi

import (
	"fmt"
	"github.com/tislib/logi/pkg/ast/plain"
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
	y.lexer.Error(fmt.Sprintf("at %s[%s]", y.lexer.GetReadString(), s))
}

func ParsePlainContent(d string) (*plain.Ast, error) {
	s := newLogiLexer(strings.NewReader(d), false)
	parser := yyNewParser()
	proxy := &yyLogiLexerProxy{lexer: s, Node: yaccNode{op: NodeOpFile}}

	parser.Parse(proxy)

	if s.Err != nil {
		return nil, s.Err
	}

	return convertNodeToLogiAst(proxy.Node)
}

func init() {
	yyErrorVerbose = true
}
