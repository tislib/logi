package logi

import (
	"encoding/json"
	"fmt"
	"github.com/tislib/logi/pkg/ast/common"
	"github.com/tislib/logi/pkg/ast/plain"
)

func convertNodeToLogiAst(node yaccNode) (*plain.Ast, error) {
	var res = new(plain.Ast)

	str, _ := json.MarshalIndent(node, "", "  ")

	fmt.Println(string(str))

	for _, child := range node.children {
		switch child.op {
		case NodeOpDefinition:
			definition, err := convertPlainDefinition(child)
			if err != nil {
				return nil, fmt.Errorf("failed to convert definition: %w", err)
			}
			res.Definitions = append(res.Definitions, *definition)
		case NodeOpFunctionDefinition:
			function, err := convertFunctionDefinition(child)
			if err != nil {
				return nil, fmt.Errorf("failed to convert function definition: %w", err)
			}
			res.Functions = append(res.Functions, *function)
		default:
			return nil, fmt.Errorf("unexpected node op: %s", child.op)
		}
	}

	return res, nil
}

func convertPlainDefinition(child yaccNode) (*plain.Definition, error) {
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
	case NodeOpArray:
		array, err := convertArray(element)

		if err != nil {
			return nil, err
		}

		statement.Kind = plain.DefinitionStatementElementKindArray
		statement.Array = array
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
	case NodeOpParameterList:
		parameterList, err := convertParameterList(element)

		if err != nil {
			return nil, err
		}

		statement.Kind = plain.DefinitionStatementElementKindParameterList
		statement.ParameterList = parameterList
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
		argument.Type = *typeDef
	}

	return argument, nil
}

func convertParameterList(element yaccNode) (*plain.DefinitionStatementElementParameterList, error) {
	parameterList := new(plain.DefinitionStatementElementParameterList)

	for _, child := range element.children {
		switch child.op {
		case NodeOpParameter:
			parameter, err := convertParameter(child)
			if err != nil {
				return nil, fmt.Errorf("failed to convert parameter: %w", err)
			}
			parameterList.Parameters = append(parameterList.Parameters, *parameter)
		}
	}

	return parameterList, nil
}

func convertParameter(element yaccNode) (*plain.DefinitionStatementElementParameter, error) {
	parameter := new(plain.DefinitionStatementElementParameter)

	value, err := statementToValue(element.value.(yaccNode))

	if err != nil {
		return nil, fmt.Errorf("failed to convert value: %w", err)
	}

	parameter.Value = *value

	return parameter, nil
}

func statementToValue(node yaccNode) (*common.Value, error) {
	statement, err := convertDefinitionStatementElement(node)

	if err != nil {
		return nil, fmt.Errorf("failed to convert statement: %w", err)
	}

	value := statement.AsValue()

	return &value, nil
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

func convertArray(element yaccNode) (*plain.DefinitionStatementElementArray, error) {
	array := new(plain.DefinitionStatementElementArray)

	for _, child := range element.children {
		st, err := convertDefinitionStatementElement(child)

		if err != nil {
			return nil, err
		}

		array.Items = append(array.Items, *st)
	}

	return array, nil
}

func convertFunctionDefinition(node yaccNode) (*plain.Function, error) {
	function := new(plain.Function)

	var nameNode = node.children[0]
	var argumentListNode = node.children[1]
	var codeBlockNode = node.children[2]

	function.Name = nameNode.value.(string)

	for _, argument := range argumentListNode.children {
		arg, err := convertArgument(argument)
		if err != nil {
			return nil, fmt.Errorf("failed to convert argument: %w", err)
		}
		function.Arguments = append(function.Arguments, *arg)
	}

	codeBlock, err := convertCodeBlock(codeBlockNode)
	if err != nil {
		return nil, fmt.Errorf("failed to convert code block: %w", err)
	}
	function.CodeBlock = *codeBlock

	return function, nil
}
