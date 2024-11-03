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

	for _, plainFunction := range plainAst.Functions {
		function, err := prepareFunction(plainFunction)

		if err != nil {
			return nil, fmt.Errorf("failed to convert function: %w", err)
		}

		result.Functions = append(result.Functions, *function)
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

		asr := analyseStatement(plainStatement, macroSyntaxStatement, syntaxElementMatch)

		// check if the statement is a property
		if canBeProperty(asr) {
			property, err := prepareProperty(plainStatement, macroSyntaxStatement, syntaxElementMatch)

			if err != nil {
				return nil, fmt.Errorf("failed to convert property: %w", err)
			}

			definition.Properties = append(definition.Properties, *property)
		}
		if canBeMethodSignature(asr) {
			methodSignature, err := prepareMethodSignature(plainStatement, macroSyntaxStatement, syntaxElementMatch)

			if err != nil {
				return nil, fmt.Errorf("failed to convert method signature: %w", err)
			}

			definition.MethodSignature = append(definition.MethodSignature, *methodSignature)
		}
		if canBeMethod(asr) {
			method, err := prepareMethod(plainStatement, macroSyntaxStatement, syntaxElementMatch)

			if err != nil {
				return nil, fmt.Errorf("failed to convert method: %w", err)
			}

			definition.Methods = append(definition.Methods, *method)
		}

		definition.PlainStatements = append(definition.PlainStatements, plainStatement)
	}

	return definition, nil
}

func prepareFunction(plainFunction plain.Function) (*logi.Function, error) {
	function := new(logi.Function)

	function.Name = plainFunction.Name
	function.CodeBlock = plainFunction.CodeBlock

	for _, argument := range plainFunction.Arguments {
		function.Arguments = append(function.Arguments, logi.Argument{Name: argument.Name, Type: argument.Type})
	}

	return function, nil
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
				property.Parameters = append(property.Parameters, logi.Parameter{Name: syntaxStatementElement.VariableKeyword.Name, Value: common.PointerValue(property.Type.AsValue())})
			} else {
				nameParts = append(nameParts, element.Identifier.Identifier)
				property.Parameters = append(property.Parameters, logi.Parameter{Name: syntaxStatementElement.VariableKeyword.Name, Value: common.PointerValue(common.StringValue(element.Identifier.Identifier))})
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

func prepareMethodSignature(statement plain.DefinitionStatement, syntaxStatement *macroAst.SyntaxStatement, syntaxStatementMatch []int) (*logi.MethodSignature, error) {
	methodSignature := new(logi.MethodSignature)
	methodSignature.Type = common.TypeDefinition{Name: "void"} // default type

	var nameParts []string

	for ei, element := range statement.Elements {
		syntaxStatementElement := syntaxStatement.Elements[syntaxStatementMatch[ei]]

		switch syntaxStatementElement.Kind {
		case macroAst.SyntaxStatementElementKindKeyword:
			nameParts = append(nameParts, element.Identifier.Identifier)
		case macroAst.SyntaxStatementElementKindVariableKeyword:
			if syntaxStatementElement.VariableKeyword.Type.Name == "Type" {
				methodSignature.Type = common.TypeDefinition{Name: element.Identifier.Identifier}
				methodSignature.Parameters = append(methodSignature.Parameters, logi.Parameter{Name: syntaxStatementElement.VariableKeyword.Name, Value: common.PointerValue(methodSignature.Type.AsValue())})
			} else {
				nameParts = append(nameParts, element.Identifier.Identifier)
				methodSignature.Parameters = append(methodSignature.Parameters, logi.Parameter{Name: syntaxStatementElement.VariableKeyword.Name, Value: common.PointerValue(common.StringValue(element.Identifier.Identifier))})
			}
		case macroAst.SyntaxStatementElementKindAttributeList:
			for _, attribute := range element.AttributeList.Attributes {
				methodSignature.Attributes = append(methodSignature.Attributes, logi.Attribute{Name: attribute.Name, Value: attribute.Value})
			}
		case macroAst.SyntaxStatementElementKindArgumentList:
			for _, argument := range element.ArgumentList.Arguments {
				methodSignature.Arguments = append(methodSignature.Arguments, logi.Argument{Name: argument.Name, Type: argument.Type})
			}
		}
	}

	methodSignature.Name = camelCaseFromNameParts(nameParts)

	return methodSignature, nil
}

func prepareMethod(statement plain.DefinitionStatement, syntaxStatement *macroAst.SyntaxStatement, syntaxStatementMatch []int) (*logi.Method, error) {
	method := new(logi.Method)
	method.Type = common.TypeDefinition{Name: "void"} // default type

	var nameParts []string

	for ei, element := range statement.Elements {
		syntaxStatementElement := syntaxStatement.Elements[syntaxStatementMatch[ei]]

		switch syntaxStatementElement.Kind {
		case macroAst.SyntaxStatementElementKindKeyword:
			nameParts = append(nameParts, element.Identifier.Identifier)
		case macroAst.SyntaxStatementElementKindVariableKeyword:
			if syntaxStatementElement.VariableKeyword.Type.Name == "Type" {
				method.Type = common.TypeDefinition{Name: element.Identifier.Identifier}
				method.Parameters = append(method.Parameters, logi.Parameter{Name: syntaxStatementElement.VariableKeyword.Name, Value: common.PointerValue(method.Type.AsValue())})
			} else {
				nameParts = append(nameParts, element.Identifier.Identifier)
				method.Parameters = append(method.Parameters, logi.Parameter{Name: syntaxStatementElement.VariableKeyword.Name, Value: common.PointerValue(common.StringValue(element.Identifier.Identifier))})
			}
		case macroAst.SyntaxStatementElementKindAttributeList:
			for _, attribute := range element.AttributeList.Attributes {
				method.Attributes = append(method.Attributes, logi.Attribute{Name: attribute.Name, Value: attribute.Value})
			}
		case macroAst.SyntaxStatementElementKindArgumentList:
			for _, argument := range element.ArgumentList.Arguments {
				method.Arguments = append(method.Arguments, logi.Argument{Name: argument.Name, Type: argument.Type})
			}
		case macroAst.SyntaxStatementElementKindCodeBlock:
			method.CodeBlock = element.CodeBlock.CodeBlock
		}
	}

	method.Name = camelCaseFromNameParts(nameParts)

	return method, nil
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

type analyseStatementResult struct {
	hasName         bool
	hasType         bool
	hasCodeBlock    bool
	hasArgumentList bool
}

func analyseStatement(statement plain.DefinitionStatement, syntaxStatement *macroAst.SyntaxStatement, syntaxStatementMatch []int) analyseStatementResult {
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

	return analyseStatementResult{
		hasName:         hasName,
		hasType:         hasType,
		hasCodeBlock:    hasCodeBlock,
		hasArgumentList: hasArgumentList,
	}
}

func canBeProperty(analyseStatementResult analyseStatementResult) bool {
	// check if the statement has a name, has a type and has not any code block or argument list
	return analyseStatementResult.hasName && analyseStatementResult.hasType && !analyseStatementResult.hasCodeBlock && !analyseStatementResult.hasArgumentList
}

func canBeMethodSignature(asr analyseStatementResult) bool {
	return asr.hasName && asr.hasArgumentList && !asr.hasCodeBlock
}

func canBeMethod(asr analyseStatementResult) bool {
	return asr.hasName && asr.hasArgumentList && asr.hasCodeBlock
}
