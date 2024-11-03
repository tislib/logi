package logi

import (
	"fmt"
	"logi/pkg/ast/common"
	"logi/pkg/ast/plain"
)

func convertNodeToLogiAst(node yaccNode) (*plain.Ast, error) {
	var res = new(plain.Ast)

	for _, child := range node.children {
		switch child.op {
		case NodeOpDefinition:
			definition, err := convertDefinition(child)
			if err != nil {
				return nil, fmt.Errorf("failed to convert definition: %w", err)
			}
			res.Definitions = append(res.Definitions, *definition)
		}
	}

	return res, nil
}

func convertDefinition(child yaccNode) (*plain.Definition, error) {
	definition := new(plain.Definition)

	var signature = child.children[0]
	var body = child.children[1]

	definition.MacroName = signature.children[0].value.(string)
	definition.Name = signature.children[1].value.(string)

	for _, statement := range body.children {
		switch statement.op {
		case NodeOpStatements:
			for _, statementElement := range statement.children {
				definitionStatementElement, err := convertDefinitionStatementElement(statementElement)
				if err != nil {
					return nil, fmt.Errorf("failed to convert definition statement element: %w", err)
				}
				definition.Statements = append(definition.Statements, *definitionStatementElement)
			}
		default:
			return nil, fmt.Errorf("unexpected node op: %s", statement.op)
		}
	}

	return definition, nil

}

func convertDefinitionStatementElement(element yaccNode) (*plain.DefinitionStatement, error) {
	definitionStatement := new(plain.DefinitionStatement)

	for _, child := range element.children {
		statementElement, err := convertStatementElement(child)
		if err != nil {
			return nil, fmt.Errorf("failed to convert statement: %w", err)
		}
		definitionStatement.Elements = append(definitionStatement.Elements, *statementElement)
	}

	return definitionStatement, nil

}

func convertStatementElement(element yaccNode) (*plain.DefinitionStatementElement, error) {
	statement := new(plain.DefinitionStatementElement)

	switch element.op {
	case NodeOpIdentifier:
		identifier, err := convertIdentifier(element)
		if err != nil {
			return nil, fmt.Errorf("failed to convert identifier: %w", err)
		}
		statement.Kind = plain.DefinitionStatementElementKindIdentifier
		statement.Identifier = identifier
	case NodeOpValue:
		value, err := convertValue(element)
		if err != nil {
			return nil, fmt.Errorf("failed to convert value: %w", err)
		}
		statement.Kind = plain.DefinitionStatementElementKindValue
		statement.Value = &plain.DefinitionStatementElementValue{Value: *value}
	case NodeOpAttributeList:
		attributeList, err := convertAttributeList(element)
		if err != nil {
			return nil, err
		}
		statement.Kind = plain.DefinitionStatementElementKindAttributeList
		statement.AttributeList = attributeList
	case NodeOpArgumentList:
		argumentList, err := convertArgumentList(element)

		if err != nil {
			return nil, err
		}

		statement.Kind = plain.DefinitionStatementElementKindArgumentList
		statement.ArgumentList = argumentList
	case NodeOpCodeBlock:
		codeBlock, err := convertCodeBlock(element)

		if err != nil {
			return nil, err
		}

		statement.Kind = plain.DefinitionStatementElementKindCodeBlock
		statement.CodeBlock = &plain.DefinitionStatementElementCodeBlock{CodeBlock: *codeBlock}
	default:
		return nil, fmt.Errorf("unexpected node op: %s", element.op)
	}

	return statement, nil

}

func convertArgumentList(element yaccNode) (*plain.DefinitionStatementElementArgumentList, error) {
	argumentList := new(plain.DefinitionStatementElementArgumentList)

	for _, child := range element.children {
		switch child.op {
		case NodeOpArgument:
			argument, err := convertArgument(child)
			if err != nil {
				return nil, fmt.Errorf("failed to convert argument: %w", err)
			}
			argumentList.Arguments = append(argumentList.Arguments, *argument)
		}
	}

	return argumentList, nil

}

func convertArgument(element yaccNode) (*plain.DefinitionStatementElementArgument, error) {
	argument := new(plain.DefinitionStatementElementArgument)

	argument.Name = element.value.(string)

	if len(element.children) > 0 {
		typeDef, err := convertTypeDefinition(element.children[0])
		if err != nil {
			return nil, fmt.Errorf("failed to convert type def: %w", err)
		}
		argument.TypeDefinition = typeDef
	}

	return argument, nil

}

func convertTypeDefinition(node yaccNode) (*common.TypeDefinition, error) {
	var result = new(common.TypeDefinition)
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

func convertIdentifier(element yaccNode) (*plain.DefinitionStatementElementIdentifier, error) {
	identifier := new(plain.DefinitionStatementElementIdentifier)

	identifier.Identifier = element.value.(string)

	return identifier, nil
}

func convertAttributeList(element yaccNode) (*plain.DefinitionStatementElementAttributeList, error) {
	attributeList := new(plain.DefinitionStatementElementAttributeList)

	for _, child := range element.children {
		switch child.op {
		case NodeOpAttribute:
			attribute, err := convertAttribute(child)
			if err != nil {
				return nil, fmt.Errorf("failed to convert attribute: %w", err)
			}
			attributeList.Attributes = append(attributeList.Attributes, *attribute)
		}
	}

	return attributeList, nil

}

func convertAttribute(element yaccNode) (*plain.DefinitionStatementElementAttribute, error) {
	attribute := new(plain.DefinitionStatementElementAttribute)

	attribute.Name = element.value.(string)

	if len(element.children) > 0 {
		value, err := convertValue(element.children[0])
		if err != nil {
			return nil, fmt.Errorf("failed to convert value: %w", err)
		}
		attribute.Value = value
	}

	return attribute, nil

}

func convertValue(element yaccNode) (*common.Value, error) {
	var value common.Value

	switch element.value.(type) {
	case string:
		value = common.StringValue(element.value.(string))
	case int:
		value = common.IntegerValue(int64(element.value.(int)))
	case bool:
		value = common.BooleanValue(element.value.(bool))
	default:
		return nil, fmt.Errorf("unexpected value type: %T", element.value)
	}

	return &value, nil

}
