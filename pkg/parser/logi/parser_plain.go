package logi

import (
	"fmt"
	"github.com/tislib/logi/pkg/ast/plain"
	"regexp"
	"strings"
)

type yyLogiLexerProxy struct {
	lexer *LogiLexer
	Node  yaccNode
	err   error
}

func (y *yyLogiLexerProxy) Lex(lval *yySymType) int {
	return y.lexer.Lex(lval)
}

func (y *yyLogiLexerProxy) Error(s string) {
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

func (y *yyLogiLexerProxy) translateToken(token string) string {
	if strings.HasSuffix(token, " or Eol") {
		return y.translateToken(strings.TrimSuffix(token, " or Eol"))
	}

	if strings.HasPrefix(token, "token_") {
		return strings.TrimPrefix(token, "token_")
	}

	if strings.HasSuffix(token, "Keyword") {
		return strings.ToLower(strings.TrimSuffix(token, "Keyword")) + " keyword"
	}

	for _, te := range Tokens() {
		if yyToknames[te.Id-yyPrivate+1] == token {
			if te.Equals != "" {
				return te.Equals
			}
		}
	}

	return token
}

func ParsePlainContent(d string) (*plain.Ast, error) {
	s := NewLogiLexer(strings.NewReader(d), false)
	parser := yyNewParser()
	proxy := &yyLogiLexerProxy{lexer: s, Node: yaccNode{op: NodeOpFile}}

	parser.Parse(proxy)

	ast, err := convertNodeToLogiAst(proxy.Node)

	if s.Err != nil {
		return ast, s.Err
	}

	if proxy.err != nil {
		return ast, proxy.err
	}

	return ast, err
}

func init() {
	yyErrorVerbose = true
}
