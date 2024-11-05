package logi

import (
	"fmt"
	"github.com/tislib/logi/pkg/ast/common"
	macroAst "github.com/tislib/logi/pkg/ast/macro"
	"github.com/tislib/logi/pkg/ast/plain"
)

func locateMacroDefinition(definition plain.Definition, ast macroAst.Ast) (*macroAst.Macro, error) {
	for _, macroDefinition := range ast.Macros {
		if macroDefinition.Name == definition.MacroName {
			return &macroDefinition, nil
		}
	}

	return nil, fmt.Errorf("macro definition not found")

}

func isValidAttributeList(plainStatementElementAttributes *plain.DefinitionStatementElementAttributeList, syntaxStatementElementAttributes *macroAst.SyntaxStatementElementAttributeList) error {
	for _, plainAttribute := range plainStatementElementAttributes.Attributes {
		found := false

		for _, syntaxAttribute := range syntaxStatementElementAttributes.Attributes {
			if plainAttribute.Name == syntaxAttribute.Name {
				found = true

				validTyped := isValidValueType(plainAttribute.Value, syntaxAttribute.Type)

				if !validTyped {
					return fmt.Errorf("attribute %s expected type %s, got %s", plainAttribute.Name, syntaxAttribute.Type.ToDisplayName(), plainAttribute.Value.ToDisplayName())
				}

				break
			}
		}

		if !found {
			return fmt.Errorf("attribute %s is not allowed", plainAttribute.Name)
		}
	}

	return nil
}

func isValidArgumentList(plainStatementElementArguments *plain.DefinitionStatementElementArgumentList, syntaxStatementElementArguments *macroAst.SyntaxStatementElementArgumentList) error {
	// check if the argument list is valid
	if syntaxStatementElementArguments.VarArgs {
		return nil
	}

	return nil
}

func isValidValueType(value *common.Value, typeDefinition common.TypeDefinition) bool {
	return true // TODO: implement type validation for values
}
