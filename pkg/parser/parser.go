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
}

func (p parser) ParseMacroContent(content string) (*macroAst.Ast, error) {
	ast, err := macro.ParseMacroContent(content)

	if err != nil {
		return nil, err
	}

	return ast, nil
}

func (p parser) ParseLogiContent(content string, macros []macroAst.Macro) (*logiAst.Ast, error) {
	ast, err := logi.Parse(content, macros)

	if err != nil {
		return nil, err
	}

	return ast, nil
}

func (p parser) ParseLogiPlainContent(content string) (*plain.Ast, error) {
	ast, err := logi.ParsePlainContent(content)

	if err != nil {
		return ast, err
	}

	return ast, nil
}

func NewParser() Parser {
	return &parser{}
}
