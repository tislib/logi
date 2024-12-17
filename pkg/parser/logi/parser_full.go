package logi

import (
	"fmt"
	"github.com/tislib/logi/pkg/ast/logi"
	macroAst "github.com/tislib/logi/pkg/ast/macro"
	"github.com/tislib/logi/pkg/ast/plain"
	"github.com/tislib/logi/pkg/parser/macro"
)

func ParseFullWithMacro(logiInput string, macroInput string, enableSourceMap bool) (*logi.Ast, error) {
	mAst, err := macro.ParseMacroContent(macroInput, enableSourceMap)

	if err != nil {
		return nil, fmt.Errorf("failed to parse macro: %w", err)
	}

	plainAst, err := ParsePlainContent(logiInput, enableSourceMap)

	if err != nil {
		return nil, fmt.Errorf("failed to parse logi: %w", err)
	}

	ast, err := prepareAst(*plainAst, *mAst)

	return ast, err
}

func Parse(logiInput string, macros []macroAst.Macro, enableSourceMap bool) (*logi.Ast, error) {
	var mAst = new(macroAst.Ast)

	mAst.Macros = macros

	plainAst, err := ParsePlainContent(logiInput, enableSourceMap)

	if err != nil {
		return nil, fmt.Errorf("failed to parse logi: %w", err)
	}

	ast, err := prepareAst(*plainAst, *mAst)

	return ast, err
}

func prepareAst(plainAst plain.Ast, macroAst macroAst.Ast) (*logi.Ast, error) {
	var result = new(logi.Ast)

	for _, plainDefinition := range plainAst.Definitions {
		// locate matching macro
		macroDefinition, err := locateMacroDefinition(plainDefinition, macroAst)

		if err != nil {
			return nil, fmt.Errorf("failed to locate macro definition: %w", err)
		}

		definition, err := prepareDefinition(plainDefinition, macroDefinition)

		if err != nil {
			return nil, fmt.Errorf("failed to convert definition: %w", err)
		}

		result.Definitions = append(result.Definitions, *definition)
	}

	return result, nil
}

func prepareDefinition(plainDefinition plain.Definition, macroDefinition *macroAst.Macro) (*logi.Definition, error) {
	definition := new(logi.Definition)

	definition.MacroName = plainDefinition.MacroName
	definition.Name = plainDefinition.Name
	definition.PlainStatements = plainDefinition.Statements

	for _, plainStatement := range plainDefinition.Statements {
		// locate matching macro syntax for the statement
		rsp := recursiveStatementParser{
			plainStatement:  plainStatement,
			macroDefinition: macroDefinition,
		}

		err := rsp.parse("")

		if err != nil {
			return nil, fmt.Errorf("failed to parse statement: %w", err)
		}

		definition.Statements = append(definition.Statements, rsp.statement)
	}

	return definition, nil
}
