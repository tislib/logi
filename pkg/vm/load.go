package vm

import (
	"fmt"
	logiAst "github.com/tislib/logi/pkg/ast/logi"
	macroAst "github.com/tislib/logi/pkg/ast/macro"
	"github.com/tislib/logi/pkg/parser/logi"
	"github.com/tislib/logi/pkg/parser/macro"
	"os"
)

func (v *vm) LoadMacroFiles(path ...string) error {
	for _, p := range path {
		data, err := os.ReadFile(p)

		if err != nil {
			return fmt.Errorf("error reading file: %v", err)
		}

		ast, err := macro.ParseMacroContent(string(data))

		if err != nil {
			return fmt.Errorf("error parsing macro content: %v", err)
		}

		v.Macros = append(v.Macros, ast.Macros...)
	}

	return nil
}

func (v *vm) LoadMacroContent(content ...string) error {
	for _, c := range content {
		ast, err := macro.ParseMacroContent(c)

		if err != nil {
			return fmt.Errorf("error parsing macro content: %v", err)
		}

		v.Macros = append(v.Macros, ast.Macros...)
	}

	return nil
}

func (v *vm) LoadMacroAst(ast ...macroAst.Ast) error {
	for _, a := range ast {
		v.Macros = append(v.Macros, a.Macros...)
	}

	return nil
}

func (v *vm) LoadLogiFiles(path ...string) ([]Definition, error) {
	var result []Definition

	for _, p := range path {
		data, err := os.ReadFile(p)

		if err != nil {
			return nil, fmt.Errorf("error reading file: %v", err)
		}

		ast, err := logi.Parse(string(data), v.Macros)

		if err != nil {
			return nil, fmt.Errorf("error parsing logi content: %v", err)
		}

		v.Logis = append(v.Logis, *ast)

		for _, astDefinition := range ast.Definitions {
			definition, err := v.prepareDefinition(astDefinition)

			if err != nil {
				return nil, fmt.Errorf("error preparing definition: %v", err)
			}

			v.Definitions = append(v.Definitions, definition)
			result = append(result, definition)
		}
	}

	return result, nil
}

func (v *vm) LoadLogiContent(content ...string) ([]Definition, error) {
	var result []Definition

	for _, c := range content {
		ast, err := logi.Parse(c, v.Macros)

		if err != nil {
			return nil, fmt.Errorf("error parsing logi content: %v", err)
		}

		v.Logis = append(v.Logis, *ast)

		for _, astDefinition := range ast.Definitions {
			definition, err := v.prepareDefinition(astDefinition)

			if err != nil {
				return nil, fmt.Errorf("error preparing definition: %v", err)
			}

			v.Definitions = append(v.Definitions, definition)
			result = append(result, definition)
		}
	}

	return result, nil
}

func (v *vm) LoadLogiAst(ast ...logiAst.Ast) ([]Definition, error) {
	var result []Definition

	for _, item := range ast {
		v.Logis = append(v.Logis, item)

		for _, astDefinition := range item.Definitions {
			definition, err := v.prepareDefinition(astDefinition)

			if err != nil {
				return nil, fmt.Errorf("error preparing definition: %v", err)
			}

			v.Definitions = append(v.Definitions, definition)
			result = append(result, definition)
		}
	}

	return result, nil
}
