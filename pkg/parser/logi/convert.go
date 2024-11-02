package macro

import (
	"fmt"
	"logi/pkg/ast"
)

func convertNodeToMacroAst(node yaccNode) (*ast.MacroAst, error) {
	var res = new(ast.MacroAst)

	if node.op != NodeOpFile {
		return nil, ErrUnexpectedNode
	}

	for _, child := range node.children {
		macro, err := convertMacro(child)

		if err != nil {
			return nil, fmt.Errorf("failed to convert syntax macro: %w", err)
		}

		res.Macros = append(res.Macros, *macro)
	}

	return res, nil
}

func convertMacro(macroNode yaccNode) (*ast.Macro, error) {
	var signature = macroNode.children[0]
	var name = signature.children[0]
	var body = macroNode.children[1]
	var kind = body.children[0].value.(string)

	var result = new(ast.Macro)
	if !NamePattern.MatchString(name.value.(string)) {
		return nil, fmt.Errorf("unexpected name value: %s", name.value)
	}

	result.Name = name.value.(string)
	switch kind {
	case "Syntax":
		result.Kind = ast.MacroKindSyntax
	default:
		return nil, fmt.Errorf("unexpected kind value: %s", kind)
	}

	for _, child := range body.children {
		switch child.op {
		case NodeOpSyntax:
			if result.Kind != ast.MacroKindSyntax {
				return nil, fmt.Errorf("syntax defined for macro of kind %s; but expected Syntax", result.Kind)
			}
			if len(child.children) != 0 {
				syntaxBody, err := convertSyntaxBody(child.children[0])

				if err != nil {
					return nil, fmt.Errorf("failed to convert syntax: %w", err)
				}

				result.Syntax = ast.Syntax{Statements: syntaxBody}
			}
		case NodeOpDefinition:
			if result.Kind != ast.MacroKindSyntax {
				return nil, fmt.Errorf("definition defined for macro of kind %s; but expected Syntax", result.Kind)
			}
			if len(child.children) != 0 {
				syntaxBody, err := convertSyntaxBody(child.children[0])

				if err != nil {
					return nil, fmt.Errorf("failed to convert definition: %w", err)
				}

				result.Definition = ast.Definition{Statements: syntaxBody}
			}
		}
	}

	return result, nil
}

func convertSyntaxBody(syntaxNode yaccNode) ([]ast.SyntaxStatement, error) {
	if syntaxNode.children == nil {
		return nil, nil
	}

	var result []ast.SyntaxStatement

	for _, child := range syntaxNode.children {
		statement, err := convertSyntaxStatement(child)

		if err != nil {
			return nil, fmt.Errorf("failed to convert syntax statement: %w", err)
		}

		result = append(result, *statement)
	}

	return result, nil
}

func convertSyntaxStatement(child yaccNode) (*ast.SyntaxStatement, error) {
	var result = new(ast.SyntaxStatement)

	for _, elementNode := range child.children {
		element, err := convertSyntaxStatementElement(elementNode)

		if err != nil {
			return nil, fmt.Errorf("failed to convert syntax statement element: %w", err)
		}

		result.Elements = append(result.Elements, *element)
	}

	return result, nil
}

func convertSyntaxStatementElement(node yaccNode) (*ast.SyntaxStatementElement, error) {
	var result = new(ast.SyntaxStatementElement)

	switch node.op {
	case NodeOpSyntaxKeywordElement:
		result.Kind = ast.SyntaxStatementElementKindKeyword
		result.KeywordDef = &ast.SyntaxStatementElementKeywordDef{Name: node.value.(string)}
	case NodeOpSyntaxVariableKeywordElement:
		result.Kind = ast.SyntaxStatementElementKindVariableKeyword
		var varName = node.children[0].value.(string)
		typeDef, err := convertTypeDefinition(node.children[1])

		if err != nil {
			return nil, fmt.Errorf("failed to convert type definition: %w", err)
		}

		result.VariableKeyword = &ast.SyntaxStatementElementVariableKeyword{Name: varName, Type: *typeDef}
	case NodeOpSyntaxParameterListElement:
		result.Kind = ast.SyntaxStatementElementKindParameterList

		var parameters []ast.SyntaxStatementElementParameter

		for _, parameterNode := range node.children {
			parameter, err := convertSyntaxStatementElementParameter(parameterNode)

			if err != nil {
				return nil, fmt.Errorf("failed to convert syntax statement element parameter: %w", err)
			}

			parameters = append(parameters, *parameter)
		}

		result.ParameterList = &ast.SyntaxStatementElementParameterList{Parameters: parameters}
	case NodeOpSyntaxArgumentListElement:
		result.Kind = ast.SyntaxStatementElementKindArgumentList

		var arguments []ast.SyntaxStatementElementArgument

		for _, argumentNode := range node.children {
			argument, err := convertSyntaxStatementElementArgument(argumentNode)

			if err != nil {
				return nil, fmt.Errorf("failed to convert syntax statement element parameter: %w", err)
			}

			arguments = append(arguments, *argument)
		}

		result.ArgumentList = &ast.SyntaxStatementElementArgumentList{Arguments: arguments, VarArgs: true}
	case NodeOpSyntaxCodeBlockElement:
		result.Kind = ast.SyntaxStatementElementKindCodeBlock
		result.CodeBlock = &ast.SyntaxStatementElementCodeBlock{}

		if len(node.children) > 0 {
			returnType, err := convertTypeDefinition(node.children[0])

			if err != nil {
				return nil, fmt.Errorf("failed to convert type definition: %w", err)
			}

			result.CodeBlock.ReturnType = *returnType
		}
	case NodeOpSyntaxAttributeListElement:
		result.Kind = ast.SyntaxStatementElementKindAttributeList

		var attributes []ast.SyntaxStatementElementAttribute

		for _, attributeNode := range node.children {
			attribute, err := convertSyntaxStatementElementAttribute(attributeNode)

			if err != nil {
				return nil, fmt.Errorf("failed to convert syntax statement element attribute: %w", err)
			}

			attributes = append(attributes, *attribute)
		}

		result.AttributeList = &ast.SyntaxStatementElementAttributeList{Attributes: attributes}
	default:
		return nil, fmt.Errorf("unexpected syntax statement element op: %s", node.op)
	}

	return result, nil
}

func convertSyntaxStatementElementParameter(node yaccNode) (*ast.SyntaxStatementElementParameter, error) {
	var result = new(ast.SyntaxStatementElementParameter)
	var varName = node.children[0].value.(string)
	typeDef, err := convertTypeDefinition(node.children[1])

	if err != nil {
		return nil, fmt.Errorf("failed to convert type definition: %w", err)
	}

	result.Name = varName
	result.Type = *typeDef

	return result, nil
}

func convertSyntaxStatementElementArgument(node yaccNode) (*ast.SyntaxStatementElementArgument, error) {
	var result = new(ast.SyntaxStatementElementArgument)
	var varName = node.children[0].value.(string)
	typeDef, err := convertTypeDefinition(node.children[1])

	if err != nil {
		return nil, fmt.Errorf("failed to convert type definition: %w", err)
	}

	result.Name = varName
	result.Type = *typeDef

	return result, nil
}

func convertSyntaxStatementElementAttribute(node yaccNode) (*ast.SyntaxStatementElementAttribute, error) {
	var result = new(ast.SyntaxStatementElementAttribute)
	var varName = node.value.(string)

	result.Name = varName

	if len(node.children) > 0 {
		typeDef, err := convertTypeDefinition(node.children[0])

		if err != nil {
			return nil, fmt.Errorf("failed to convert type definition: %w", err)
		}

		result.Type = *typeDef
	}

	return result, nil
}

func convertTypeDefinition(node yaccNode) (*ast.TypeDefinition, error) {
	var result = new(ast.TypeDefinition)
	result.Name = node.value.(string)

	if len(node.children) > 0 {
		for _, child := range node.children {
			subType, err := convertTypeDefinition(child)

			if err != nil {
				return nil, fmt.Errorf("failed to convert type definition: %w", err)
			}

			result.SubTypes = append(result.SubTypes, *subType)
		}
	}

	return result, nil
}
