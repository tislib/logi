package macro

import (
	"fmt"
	"github.com/tislib/logi/pkg/ast/common"
	astMacro "github.com/tislib/logi/pkg/ast/macro"
)

func convertNodeToMacroAst(node yaccNode) (*astMacro.Ast, error) {
	var res = new(astMacro.Ast)

	if node.op != NodeOpFile {
		return nil, ErrUnexpectedNode
	}

	for _, child := range node.children {
		macro, err := convertMacro(child)

		if err != nil {
			return nil, err
		}

		res.Macros = append(res.Macros, *macro)
	}

	return res, nil
}

func convertMacro(macroNode yaccNode) (*astMacro.Macro, error) {
	var signature = macroNode.children[0]
	var name = signature.children[0]
	var body = macroNode.children[1]
	var kind = body.children[0].value.(string)

	var result = new(astMacro.Macro)
	if !NamePattern.MatchString(name.value.(string)) {
		return nil, fmt.Errorf("unexpected name value: %s", name.value)
	}

	result.Name = name.value.(string)
	switch kind {
	case "Syntax":
		result.Kind = astMacro.KindSyntax
	default:
		return nil, newErrorFromNode(body.children[0], fmt.Sprintf("unexpected kind value: \"%s\", expecting \"Syntax\"", kind))
	}

	for _, child := range body.children {
		switch child.op {
		case NodeOpSyntax:
			if result.Kind != astMacro.KindSyntax {
				return nil, fmt.Errorf("syntax defined for macro of kind %s; but expected Syntax", result.Kind)
			}
			if len(child.children) != 0 {
				syntaxBody, err := convertSyntaxBody(child.children[0])

				if err != nil {
					return nil, err
				}

				result.Syntax = astMacro.Syntax{Statements: syntaxBody}
			}
		case NodeOpTypes:
			if result.Kind != astMacro.KindSyntax {
				return nil, fmt.Errorf("types defined for macro of kind %s; but expected Syntax", result.Kind)
			}
			if len(child.children) != 0 {
				types, err := convertTypes(child.children[0])

				if err != nil {
					return nil, err
				}

				result.Types = *types
			}
		}
	}

	return result, nil
}

func newErrorFromNode(name yaccNode, msg string) error {
	return newError(name.location.Line, name.location.Column, name.token.Value, msg)
}

func convertSyntaxBody(syntaxNode yaccNode) ([]astMacro.SyntaxStatement, error) {
	if syntaxNode.children == nil {
		return nil, nil
	}

	var result []astMacro.SyntaxStatement

	for _, child := range syntaxNode.children {
		statement, err := convertSyntaxStatement(child)

		if err != nil {
			return nil, err
		}

		result = append(result, *statement)
	}

	return result, nil
}

func convertTypes(typesNode yaccNode) (*astMacro.Types, error) {
	if typesNode.children == nil {
		return nil, nil
	}

	var result []astMacro.TypeStatement

	for _, child := range typesNode.children {
		statement, err := convertTypeStatement(child)

		if err != nil {
			return nil, err
		}

		result = append(result, *statement)
	}

	return &astMacro.Types{Types: result}, nil
}

func convertTypeStatement(node yaccNode) (*astMacro.TypeStatement, error) {
	var result = new(astMacro.TypeStatement)
	var name = node.children[0].value.(string)
	var items = node.children[1].children

	result.Name = name

	for _, child := range items {
		element, err := convertSyntaxStatementElement(child)

		if err != nil {
			return nil, err
		}

		result.Elements = append(result.Elements, *element)
	}

	return result, nil
}

func convertSyntaxStatement(child yaccNode) (*astMacro.SyntaxStatement, error) {
	var result = new(astMacro.SyntaxStatement)

	for _, elementNode := range child.children {
		element, err := convertSyntaxStatementElement(elementNode)

		if err != nil {
			return nil, err
		}

		result.Elements = append(result.Elements, *element)
	}

	return result, nil
}

func convertSyntaxStatementElement(node yaccNode) (*astMacro.SyntaxStatementElement, error) {
	var result = new(astMacro.SyntaxStatementElement)

	switch node.op {
	case NodeOpSyntaxKeywordElement:
		result.Kind = astMacro.SyntaxStatementElementKindKeyword
		result.KeywordDef = &astMacro.SyntaxStatementElementKeywordDef{Name: node.value.(string)}
	case NodeOpSyntaxTypeReferenceElement:
		result.Kind = astMacro.SyntaxStatementElementKindTypeReference
		result.TypeReference = &astMacro.SyntaxStatementElementTypeReference{Name: node.value.(string)}
	case NodeOpSyntaxVariableKeywordElement:
		result.Kind = astMacro.SyntaxStatementElementKindVariableKeyword
		var varName = node.children[0].value.(string)
		typeDef, err := convertTypeDefinition(node.children[1])

		if err != nil {
			return nil, err
		}

		result.VariableKeyword = &astMacro.SyntaxStatementElementVariableKeyword{Name: varName, Type: *typeDef}
	case NodeOpSyntaxCombinationElement:
		result.Kind = astMacro.SyntaxStatementElementKindCombination

		var elements []astMacro.SyntaxStatementElement
		for _, elementNode := range node.children {
			element, err := convertSyntaxStatementElement(elementNode)

			if err != nil {
				return nil, err
			}

			elements = append(elements, *element)
		}

		result.Combination = &astMacro.SyntaxStatementElementCombination{Elements: elements}
	case NodeOpSyntaxStructureElement:
		result.Kind = astMacro.SyntaxStatementElementKindStructure

		var statements []astMacro.SyntaxStatement
		for _, elementNode := range node.children {
			statement, err := convertSyntaxStatement(elementNode)

			if err != nil {
				return nil, err
			}

			statements = append(statements, *statement)
		}

		result.Structure = &astMacro.SyntaxStatementElementStructure{Statements: statements}
	case NodeOpSyntaxParameterListElement:
		result.Kind = astMacro.SyntaxStatementElementKindParameterList

		var parameters []astMacro.SyntaxStatementElementParameter

		for _, parameterNode := range node.children {
			parameter, err := convertSyntaxStatementElementParameter(parameterNode)

			if err != nil {
				return nil, err
			}

			parameters = append(parameters, *parameter)
		}

		result.ParameterList = &astMacro.SyntaxStatementElementParameterList{Parameters: parameters}
	case NodeOpSyntaxArgumentListElement:
		result.Kind = astMacro.SyntaxStatementElementKindArgumentList

		var arguments []astMacro.SyntaxStatementElementArgument

		for _, argumentNode := range node.children {
			argument, err := convertSyntaxStatementElementArgument(argumentNode)

			if err != nil {
				return nil, err
			}

			arguments = append(arguments, *argument)
		}

		result.ArgumentList = &astMacro.SyntaxStatementElementArgumentList{Arguments: arguments, VarArgs: true}
	case NodeOpSyntaxCodeBlockElement:
		result.Kind = astMacro.SyntaxStatementElementKindCodeBlock
		result.CodeBlock = &astMacro.SyntaxStatementElementCodeBlock{}
	case NodeOpSyntaxExpressionBlockElement:
		result.Kind = astMacro.SyntaxStatementElementKindExpressionBlock
		result.ExpressionBlock = &astMacro.SyntaxStatementElementExpressionBlock{}
	case NodeOpSyntaxAttributeListElement:
		result.Kind = astMacro.SyntaxStatementElementKindAttributeList

		var attributes []astMacro.SyntaxStatementElementAttribute

		for _, attributeNode := range node.children {
			attribute, err := convertSyntaxStatementElementAttribute(attributeNode)

			if err != nil {
				return nil, err
			}

			attributes = append(attributes, *attribute)
		}

		result.AttributeList = &astMacro.SyntaxStatementElementAttributeList{Attributes: attributes}
	default:
		return nil, fmt.Errorf("unexpected syntax statement element op: %s", node.op)
	}

	return result, nil
}

func convertSyntaxStatementElementParameter(node yaccNode) (*astMacro.SyntaxStatementElementParameter, error) {
	var result = new(astMacro.SyntaxStatementElementParameter)
	var varName = node.children[0].value.(string)
	typeDef, err := convertTypeDefinition(node.children[1])

	if err != nil {
		return nil, err
	}

	result.Name = varName
	result.Type = *typeDef

	return result, nil
}

func convertSyntaxStatementElementArgument(node yaccNode) (*astMacro.SyntaxStatementElementArgument, error) {
	var result = new(astMacro.SyntaxStatementElementArgument)
	var varName = node.children[0].value.(string)
	typeDef, err := convertTypeDefinition(node.children[1])

	if err != nil {
		return nil, err
	}

	result.Name = varName
	result.Type = *typeDef

	return result, nil
}

func convertSyntaxStatementElementAttribute(node yaccNode) (*astMacro.SyntaxStatementElementAttribute, error) {
	var result = new(astMacro.SyntaxStatementElementAttribute)
	var varName = node.value.(string)

	result.Name = varName

	if len(node.children) > 0 {
		typeDef, err := convertTypeDefinition(node.children[0])

		if err != nil {
			return nil, err
		}

		result.Type = *typeDef
	}

	return result, nil
}

func convertTypeDefinition(node yaccNode) (*common.TypeDefinition, error) {
	var result = new(common.TypeDefinition)
	result.Name = node.value.(string)

	if len(node.children) > 0 {
		for _, child := range node.children {
			subType, err := convertTypeDefinition(child)

			if err != nil {
				return nil, err
			}

			result.SubTypes = append(result.SubTypes, *subType)
		}
	}

	return result, nil
}
