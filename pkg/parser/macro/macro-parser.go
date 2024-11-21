package macro

import (
	"fmt"
	astMacro "github.com/tislib/logi/pkg/ast/macro"
	"regexp"
	"strings"
)

type yyMakroLexerProxy struct {
	lexer *macroLexer
	Node  yaccNode
	err   error
}

func (y *yyMakroLexerProxy) Lex(lval *yySymType) int {
	return y.lexer.Lex(lval)
}

func (y *yyMakroLexerProxy) Error(s string) {
	lastToken := y.lexer.lexer.GetLastToken()
	lastLocation := y.lexer.lexer.GetLastLocation()

	// syntax error: unexpected token_identifier, expecting MacroKeyword or Eol
	var unexpectedPattern = regexp.MustCompile(`syntax error: unexpected (?P<unexpected>.+), expecting (?P<expected>.+)`)

	if unexpectedPattern.MatchString(s) {
		matches := unexpectedPattern.FindStringSubmatch(s)
		unexpected := matches[1]
		expected := matches[2]

		unexpected = y.translateToken(unexpected)
		expected = y.translateToken(expected)

		y.err = newError(lastLocation.Line, lastLocation.Column, fmt.Sprintf("%s", lastToken.Value), fmt.Sprintf("unexpected %s \"%s\", expecting %s", unexpected, lastToken.Value, expected))
		return
	}

	y.err = newError(lastLocation.Line, lastLocation.Column, fmt.Sprintf("%s", lastToken.Value), s)
}

func (y *yyMakroLexerProxy) translateToken(token string) string {
	if strings.HasSuffix(token, " or Eol") {
		return y.translateToken(strings.TrimSuffix(token, " or Eol"))
	}

	if strings.HasPrefix(token, "token_") {
		return strings.TrimPrefix(token, "token_")
	}

	if strings.HasSuffix(token, "Keyword") {
		return strings.ToLower(strings.TrimSuffix(token, "Keyword")) + " keyword"
	}

	for _, te := range tokens() {
		if yyToknames[te.Id-yyPrivate+1] == token {
			if te.Equals != "" {
				return te.Equals
			}
		}
	}

	return token
}

func ParseMacroContent(d string, enableSourceMap bool) (*astMacro.Ast, error) {
	s := newMacroLexer(strings.NewReader(d), false)
	parser := yyNewParser()
	proxy := &yyMakroLexerProxy{lexer: s, Node: yaccNode{op: NodeOpFile}}

	parser.Parse(proxy)

	var c = &converter{enableSourceMap}
	var result, err = c.convertNodeToMacroAst(proxy.Node)

	if s.Err != nil {
		return result, s.Err
	}

	if proxy.err != nil {
		return result, proxy.err
	}

	return result, err
}

func init() {
	yyErrorVerbose = true
}
