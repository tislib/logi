package logi

import (
	"fmt"
	"logi/pkg/ast/common"
	"logi/pkg/ast/logi"
	macroAst "logi/pkg/ast/macro"
	"logi/pkg/ast/plain"
	"logi/pkg/parser/macro"
	"strings"
)

func ParseFullWithMacro(logiInput string, macroInput string) (*logi.Ast, error) {
	mAst, err := macro.ParseMacroContent(macroInput)

	if err != nil {
		return nil, fmt.Errorf("failed to parse macro: %w", err)
	}

	plainAst, err := ParsePlainContent(logiInput)

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

	for _, plainStatement := range plainDefinition.Statements {
		// locate matching macro syntax for the statement
		macroSyntaxStatement, syntaxElementMatch, err := locateMacroSyntaxStatement(plainStatement, macroDefinition)

		if err != nil {
			return nil, fmt.Errorf("failed to locate syntax statement: %w", err)
		}

		// check if the statement is a property
		if canBeProperty(plainStatement, macroSyntaxStatement, syntaxElementMatch) {
			property, err := prepareProperty(plainStatement, macroSyntaxStatement, syntaxElementMatch)

			if err != nil {
				return nil, fmt.Errorf("failed to convert property: %w", err)
			}

			definition.Properties = append(definition.Properties, *property)
		}

		definition.PlainStatements = append(definition.PlainStatements, plainStatement)
	}

	return definition, nil

}

func prepareProperty(statement plain.DefinitionStatement, syntaxStatement *macroAst.SyntaxStatement, syntaxStatementMatch []int) (*logi.Property, error) {
	property := new(logi.Property)

	var nameParts []string

	for ei, element := range statement.Elements {
		syntaxStatementElement := syntaxStatement.Elements[syntaxStatementMatch[ei]]

		switch syntaxStatementElement.Kind {
		case macroAst.SyntaxStatementElementKindKeyword:
			nameParts = append(nameParts, element.Identifier.Identifier)
		case macroAst.SyntaxStatementElementKindVariableKeyword:
			if syntaxStatementElement.VariableKeyword.Type.Name == "Type" {
				property.Type = common.TypeDefinition{Name: element.Identifier.Identifier}
			} else {
				nameParts = append(nameParts, element.Identifier.Identifier)
			}
		case macroAst.SyntaxStatementElementKindAttributeList:
			for _, attribute := range element.AttributeList.Attributes {
				property.Attributes = append(property.Attributes, logi.Attribute{Name: attribute.Name, Value: attribute.Value})
			}
		}
	}

	property.Name = camelCaseFromNameParts(nameParts)

	return property, nil

}

func camelCaseFromNameParts(parts []string) string {
	var result string

	for i, part := range parts {
		if i == 0 {
			result += part
		} else {
			result += strings.ToUpper(part[:1]) + part[1:]
		}
	}

	return result
}

func canBeProperty(statement plain.DefinitionStatement, syntaxStatement *macroAst.SyntaxStatement, syntaxStatementMatch []int) bool {
	// check if the statement has a name, has a type and has not any code block or argument list

	var hasName = false
	var hasType = false
	var hasCodeBlock = false
	var hasArgumentList = false

	for ei, element := range statement.Elements {
		syntaxStatementElement := syntaxStatement.Elements[syntaxStatementMatch[ei]]

		switch syntaxStatementElement.Kind {
		case macroAst.SyntaxStatementElementKindKeyword:
			hasName = true
		case macroAst.SyntaxStatementElementKindVariableKeyword:
			if syntaxStatementElement.VariableKeyword.Type.Name == "Type" {
				hasType = true
			} else {
				hasName = true
			}
		}

		if element.Kind == plain.DefinitionStatementElementKindCodeBlock {
			hasCodeBlock = true
		}
		if element.Kind == plain.DefinitionStatementElementKindArgumentList {
			hasArgumentList = true
		}
	}

	return hasName && hasType && !hasCodeBlock && !hasArgumentList
}
