package logi

import (
	"fmt"
	"github.com/tislib/logi/pkg/ast/common"
	"github.com/tislib/logi/pkg/ast/plain"
)

func convertNodeToLogiAst(node yaccNode) (*plain.Ast, error) {
	var res = new(plain.Ast)

	for _, child := range node.children {
		switch child.op {
		case NodeOpDefinition:
			definition, err := convertPlainDefinition(child)
			if err != nil {
				return res, fmt.Errorf("failed to convert definition: %w", err)
			}
			res.Definitions = append(res.Definitions, *definition)
		case NodeOpFunctionDefinition:
			function, err := convertFunctionDefinition(child)
			if err != nil {
				return res, fmt.Errorf("failed to convert function definition: %w", err)
			}
			res.Functions = append(res.Functions, *function)
		default:
			return res, fmt.Errorf("unexpected node op: %s", child.op)
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
				definitionStatement, err := convertDefinitionStatementElement(statementElement)
				if err != nil {
					return definition, fmt.Errorf("failed to convert definition statement element: %w", err)
				}
				definition.Statements = append(definition.Statements, *definitionStatement)
			}
		default:
			return definition, fmt.Errorf("unexpected node op: %s", statement.op)
		}
	}

	definition.NameSourceLocation = common.SourceLocation{
		Line:   signature.children[1].location.Line,
		Column: signature.children[1].location.Column,
	}
	definition.MacroNameSourceLocation = common.SourceLocation{
		Line:   signature.children[0].location.Line,
		Column: signature.children[0].location.Column,
	}

	return definition, nil

}

func convertDefinitionStatementElement(element yaccNode) (*plain.DefinitionStatement, error) {
	definitionStatement := new(plain.DefinitionStatement)

	for _, child := range element.children {
		statementElement, err := convertStatementElement(child)
		if err != nil {
			return definitionStatement, fmt.Errorf("failed to convert statement: %w", err)
		}
		definitionStatement.Elements = append(definitionStatement.Elements, *statementElement)
	}

	definitionStatement.SourceLocation = common.SourceLocation{
		Line:   element.location.Line,
		Column: element.location.Column,
	}

	return definitionStatement, nil

}

func convertStatementElement(element yaccNode) (*plain.DefinitionStatementElement, error) {
	statementElement := new(plain.DefinitionStatementElement)

	switch element.op {
	case NodeOpIdentifier:
		identifier, err := convertIdentifier(element)
		if err != nil {
			return statementElement, fmt.Errorf("failed to convert identifier: %w", err)
		}
		statementElement.Kind = plain.DefinitionStatementElementKindIdentifier
		statementElement.Identifier = identifier
	case NodeOpValue:
		value, err := convertValue(element)
		if err != nil {
			return statementElement, fmt.Errorf("failed to convert value: %w", err)
		}
		statementElement.Kind = plain.DefinitionStatementElementKindValue
		statementElement.Value = &plain.DefinitionStatementElementValue{Value: *value}
	case NodeOpArray:
		array, err := convertArray(element)

		if err != nil {
			return statementElement, err
		}

		statementElement.Kind = plain.DefinitionStatementElementKindArray
		statementElement.Array = array
	case NodeOpAttributeList:
		attributeList, err := convertAttributeList(element)
		if err != nil {
			return statementElement, err
		}
		statementElement.Kind = plain.DefinitionStatementElementKindAttributeList
		statementElement.AttributeList = attributeList
	case NodeOpArgumentList:
		argumentList, err := convertArgumentList(element)

		if err != nil {
			return statementElement, err
		}

		statementElement.Kind = plain.DefinitionStatementElementKindArgumentList
		statementElement.ArgumentList = argumentList
	case NodeOpParameterList:
		parameterList, err := convertParameterList(element)

		if err != nil {
			return statementElement, err
		}

		statementElement.Kind = plain.DefinitionStatementElementKindParameterList
		statementElement.ParameterList = parameterList
	case NodeOpCodeBlock:
		codeBlock, err := convertCodeBlock(element)

		if err != nil {
			return statementElement, err
		}

		statementElement.Kind = plain.DefinitionStatementElementKindCodeBlock
		statementElement.CodeBlock = &plain.DefinitionStatementElementCodeBlock{CodeBlock: *codeBlock}
	case NodeOpStruct:
		structure, err := convertStruct(element)

		if err != nil {
			return statementElement, err
		}

		statementElement.Kind = plain.DefinitionStatementElementKindStruct
		statementElement.Struct = structure
	case NodeOpJsonObject:
		jsonObject, err := convertJsonObject(element)

		if err != nil {
			return statementElement, err
		}

		statementElement.Kind = plain.DefinitionStatementElementKindValue
		statementElement.Value = &plain.DefinitionStatementElementValue{Value: jsonObject}
	default:
		return statementElement, fmt.Errorf("unexpected node op: %s", element.op)
	}

	statementElement.SourceLocation = common.SourceLocation{
		Line:   element.location.Line,
		Column: element.location.Column,
	}

	return statementElement, nil
}

func convertJsonObject(element yaccNode) (common.Value, error) {
	var result = common.Value{
		Kind: common.ValueKindMap,
		Map:  make(map[string]common.Value),
	}

	for _, child := range element.children {
		value, err := convertJsonValue(child.children[0])

		if err != nil {
			return result, err
		}
		result.Map[child.value.(string)] = value
	}

	return result, nil
}

func convertJsonArray(element yaccNode) (common.Value, error) {
	var result = common.Value{
		Kind:  common.ValueKindArray,
		Array: make([]common.Value, 0),
	}

	for _, child := range element.children {
		value, err := convertJsonValue(child)

		if err != nil {
			return common.Value{}, err
		}
		result.Array = append(result.Array, value)
	}

	return result, nil
}

func convertJsonValue(node yaccNode) (common.Value, error) {
	switch node.op {
	case NodeOpJsonObjectItemValue:
		switch node.value.(type) {
		case string:
			return common.StringValue(node.value.(string)), nil
		case int:
			return common.IntegerValue(int64(node.value.(int))), nil
		case float64:
			return common.FloatValue(node.value.(float64)), nil
		case bool:
			return common.BooleanValue(node.value.(bool)), nil
		default:
			panic(fmt.Sprintf("unexpected value type: %T", node.value))
		}
	case NodeOpJsonArray:
		return convertJsonArray(node)
	case NodeOpJsonObject:
		return convertJsonObject(node)
	default:
		panic(fmt.Sprintf("unexpected node op: %s", node.op))
	}
}

func convertArgumentList(element yaccNode) (*plain.DefinitionStatementElementArgumentList, error) {
	argumentList := new(plain.DefinitionStatementElementArgumentList)

	for _, child := range element.children {
		switch child.op {
		case NodeOpArgument:
			argument, err := convertArgument(child)
			if err != nil {
				return argumentList, fmt.Errorf("failed to convert argument: %w", err)
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
			return argument, fmt.Errorf("failed to convert type def: %w", err)
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
				return parameterList, fmt.Errorf("failed to convert parameter: %w", err)
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
		return parameter, fmt.Errorf("failed to convert value: %w", err)
	}

	parameter.Value = *value

	return parameter, nil
}

func statementToValue(node yaccNode) (*common.Value, error) {
	statement, err := convertDefinitionStatementElement(node)

	value := statement.AsValue()

	if err != nil {
		return &value, fmt.Errorf("failed to convert statement: %w", err)
	}

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
	case float64:
		value = common.FloatValue(element.value.(float64))
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

func convertStruct(element yaccNode) (*plain.DefinitionStatementElementStruct, error) {
	result := new(plain.DefinitionStatementElementStruct)

	for _, child := range element.children[0].children[0].children {
		st, err := convertDefinitionStatementElement(child)

		if err != nil {
			return nil, err
		}

		result.Statements = append(result.Statements, *st)
	}

	return result, nil
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
