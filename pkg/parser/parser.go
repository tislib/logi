package parser

import (
	logiAst "github.com/tislib/logi/pkg/ast/logi"
	macroAst "github.com/tislib/logi/pkg/ast/macro"
	"github.com/tislib/logi/pkg/ast/plain"
	"github.com/tislib/logi/pkg/parser/logi"
	"github.com/tislib/logi/pkg/parser/macro"
)

type Parser interface {
	ParseMacroContent(content string) (*macroAst.Ast, error)
	ParseLogiPlainContent(content string) (*plain.Ast, error)
	ParseLogiContent(content string, macros []macroAst.Macro) (*logiAst.Ast, error)
}

type parser struct {
	enableSourceMap bool
}

func (p parser) ParseMacroContent(content string) (*macroAst.Ast, error) {
	ast, err := macro.ParseMacroContent(content, p.enableSourceMap)

	if err != nil {
		return ast, err
	}

	return ast, nil
}

func (p parser) ParseLogiContent(content string, macros []macroAst.Macro) (*logiAst.Ast, error) {
	ast, err := logi.Parse(content, macros, p.enableSourceMap)

	if err != nil {
		return nil, err
	}

	return ast, nil
}

func (p parser) ParseLogiPlainContent(content string) (*plain.Ast, error) {
	ast, err := logi.ParsePlainContent(content, p.enableSourceMap)

	if err != nil {
		return ast, err
	}

	return ast, nil
}

func NewParser(enableSourceMap bool) Parser {
	return &parser{enableSourceMap}
}
