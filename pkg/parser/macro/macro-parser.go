package macro

import (
	"fmt"
	astMacro "github.com/tislib/logi/pkg/ast/macro"
	"strings"
)

type yyMakroLexerProxy struct {
	lexer *macroLexer
	Node  yaccNode
}

func (y *yyMakroLexerProxy) Lex(lval *yySymType) int {
	return y.lexer.Lex(lval)
}

func (y *yyMakroLexerProxy) Error(s string) {
	y.lexer.Error(fmt.Sprintf("at %s[%s]", y.lexer.GetReadString(), s))
}

func ParseMacroContent(d string) (*astMacro.Ast, error) {
	s := newMacroLexer(strings.NewReader(d), false)
	parser := yyNewParser()
	proxy := &yyMakroLexerProxy{lexer: s, Node: yaccNode{op: NodeOpFile}}

	parser.Parse(proxy)

	if s.Err != nil {
		return nil, s.Err
	}

	return convertNodeToMacroAst(proxy.Node)
}

func init() {
	yyErrorVerbose = true
}
