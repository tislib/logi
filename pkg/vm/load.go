package vm

import (
	"fmt"
	logiAst "github.com/tislib/logi/pkg/ast/logi"
	macroAst "github.com/tislib/logi/pkg/ast/macro"
	"github.com/tislib/logi/pkg/parser/logi"
	"github.com/tislib/logi/pkg/parser/macro"
	"os"
)

func (v *vm) LoadMacroFile(path ...string) error {
	for _, p := range path {
		data, err := os.ReadFile(p)

		if err != nil {
			return fmt.Errorf("error reading file: %v", err)
		}

		ast, err := macro.ParseMacroContent(string(data), v.enableSourceMap)

		if err != nil {
			return fmt.Errorf("error parsing macro content: %v", err)
		}

		v.Macros = append(v.Macros, ast.Macros...)
		v.MacroContents[p] = string(data)
	}

	return nil
}

func (v *vm) LoadMacroContent(content ...string) error {
	for _, c := range content {
		ast, err := macro.ParseMacroContent(c, v.enableSourceMap)

		if err != nil {
			return fmt.Errorf("error parsing macro content: %v", err)
		}

		v.Macros = append(v.Macros, ast.Macros...)
		v.MacroContents[c] = c
	}

	return nil
}

func (v *vm) LoadMacroAst(ast ...macroAst.Ast) error {
	for _, a := range ast {
		v.Macros = append(v.Macros, a.Macros...)
	}

	return nil
}

func (v *vm) LoadLogiFile(path ...string) ([]logiAst.Definition, error) {
	var result []logiAst.Definition

	for _, p := range path {
		data, err := os.ReadFile(p)

		if err != nil {
			return nil, fmt.Errorf("error reading file: %v", err)
		}

		ast, err := logi.Parse(string(data), v.Macros, false)

		if err != nil {
			return nil, fmt.Errorf("error parsing logi content: %v", err)
		}

		v.Definitions = append(v.Definitions, ast.Definitions...)
	}

	return result, nil
}

func (v *vm) LoadLogiContent(content ...string) ([]logiAst.Definition, error) {
	var result []logiAst.Definition

	for _, c := range content {
		ast, err := logi.Parse(c, v.Macros, false)

		if err != nil {
			return nil, fmt.Errorf("error parsing logi content: %v", err)
		}

		v.Definitions = append(v.Definitions, ast.Definitions...)
		result = append(result, ast.Definitions...)
	}

	return result, nil
}

func (v *vm) LoadLogiAst(ast ...logiAst.Ast) ([]logiAst.Definition, error) {
	var result []logiAst.Definition

	for _, item := range ast {
		v.Definitions = append(v.Definitions, item.Definitions...)
	}

	return result, nil
}
