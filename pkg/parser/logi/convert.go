package logi

import (
	"fmt"
	"github.com/tislib/logi/pkg/ast/common"
	"github.com/tislib/logi/pkg/ast/plain"
)

type converter struct {
	enableSourceMap bool
}

func (c *converter) convertNodeToLogiAst(node yaccNode) (*plain.Ast, error) {
	var res = new(plain.Ast)

	for _, child := range node.children {
		switch child.op {
		case NodeOpDefinition:
			definition, err := c.convertPlainDefinition(child)
			if err != nil {
				return res, fmt.Errorf("failed to convert definition: %w", err)
			}
			res.Definitions = append(res.Definitions, *definition)
		default:
			return res, fmt.Errorf("unexpected node op: %s", child.op)
		}
	}

	return res, nil
}

func (c *converter) convertPlainDefinition(child yaccNode) (*plain.Definition, error) {
	definition := new(plain.Definition)

	var signature = child.children[0]
	var body = child.children[1]

	definition.MacroName = signature.children[0].value.(string)
	definition.Name = signature.children[1].value.(string)

	for _, statement := range body.children {
		switch statement.op {
		case NodeOpStatements:
			for _, statementElement := range statement.children {
				definitionStatement, err := c.convertDefinitionStatementElement(statementElement)
				if err != nil {
					return definition, fmt.Errorf("failed to convert definition statement element: %w", err)
				}
				definition.Statements = append(definition.Statements, *definitionStatement)
			}
		default:
			return definition, fmt.Errorf("unexpected node op: %s", statement.op)
		}
	}

	if c.enableSourceMap {
		definition.NameSourceLocation = common.SourceLocation{
			Line:   signature.children[1].location.Line,
			Column: signature.children[1].location.Column,
		}
		definition.MacroNameSourceLocation = common.SourceLocation{
			Line:   signature.children[0].location.Line,
			Column: signature.children[0].location.Column,
		}
	}

	return definition, nil

}

func (c *converter) convertDefinitionStatementElement(element yaccNode) (*plain.DefinitionStatement, error) {
	definitionStatement := new(plain.DefinitionStatement)

	for _, child := range element.children {
		statementElement, err := c.convertStatementElement(child)
		if err != nil {
			return definitionStatement, fmt.Errorf("failed to convert statement: %w", err)
		}
		definitionStatement.Elements = append(definitionStatement.Elements, *statementElement)
	}

	if c.enableSourceMap {
		definitionStatement.SourceLocation = common.SourceLocation{
			Line:   element.location.Line,
			Column: element.location.Column,
		}
	}

	return definitionStatement, nil

}

func (c *converter) convertStatementElement(element yaccNode) (*plain.DefinitionStatementElement, error) {
	statementElement := new(plain.DefinitionStatementElement)

	switch element.op {
	case NodeOpIdentifier:
		identifier, err := c.convertIdentifier(element)
		if err != nil {
			return statementElement, fmt.Errorf("failed to convert identifier: %w", err)
		}
		statementElement.Kind = plain.DefinitionStatementElementKindIdentifier
		statementElement.Identifier = identifier
	case NodeOpValue:
		value, err := c.convertValue(element)
		if err != nil {
			return statementElement, fmt.Errorf("failed to convert value: %w", err)
		}
		statementElement.Kind = plain.DefinitionStatementElementKindValue
		statementElement.Value = &plain.DefinitionStatementElementValue{Value: *value}
	case NodeOpArray:
		array, err := c.convertArray(element)

		if err != nil {
			return statementElement, err
		}

		statementElement.Kind = plain.DefinitionStatementElementKindArray
		statementElement.Array = array
	case NodeOpAttributeList:
		attributeList, err := c.convertAttributeList(element)
		if err != nil {
			return statementElement, err
		}
		statementElement.Kind = plain.DefinitionStatementElementKindAttributeList
		statementElement.AttributeList = attributeList
	case NodeOpArgumentList:
		argumentList, err := c.convertArgumentList(element)

		if err != nil {
			return statementElement, err
		}

		statementElement.Kind = plain.DefinitionStatementElementKindArgumentList
		statementElement.ArgumentList = argumentList
	case NodeOpParameterList:
		parameterList, err := c.convertParameterList(element)

		if err != nil {
			return statementElement, err
		}

		statementElement.Kind = plain.DefinitionStatementElementKindParameterList
		statementElement.ParameterList = parameterList
	case NodeOpStruct:
		structure, err := c.convertStruct(element)

		if err != nil {
			return statementElement, err
		}

		statementElement.Kind = plain.DefinitionStatementElementKindStruct
		statementElement.Struct = structure
	case NodeOpJsonObject:
		jsonObject, err := c.convertJsonObject(element)

		if err != nil {
			return statementElement, err
		}

		statementElement.Kind = plain.DefinitionStatementElementKindValue
		statementElement.Value = &plain.DefinitionStatementElementValue{Value: jsonObject}
	case NodeOpFunctionCall:
		functionCall, err := c.convertFunctionCall(element)

		if err != nil {
			return statementElement, err
		}

		statementElement.Kind = plain.DefinitionStatementElementKindExpression
		statementElement.Expression = &common.Expression{
			FuncCall: functionCall,
		}
	case NodeOpExpression:
	default:
		return statementElement, fmt.Errorf("unexpected node op: %s", element.op)
	}

	if c.enableSourceMap {
		statementElement.SourceLocation = common.SourceLocation{
			Line:   element.location.Line,
			Column: element.location.Column,
		}
	}

	return statementElement, nil
}

func (c *converter) convertJsonObject(element yaccNode) (common.Value, error) {
	var result = common.Value{
		Kind: common.ValueKindMap,
		Map:  make(map[string]common.Value),
	}

	for _, child := range element.children {
		value, err := c.convertJsonValue(child.children[0])

		if err != nil {
			return result, err
		}
		result.Map[child.value.(string)] = value
	}

	return result, nil
}

func (c *converter) convertJsonArray(element yaccNode) (common.Value, error) {
	var result = common.Value{
		Kind:  common.ValueKindArray,
		Array: make([]common.Value, 0),
	}

	for _, child := range element.children {
		value, err := c.convertJsonValue(child)

		if err != nil {
			return common.Value{}, err
		}
		result.Array = append(result.Array, value)
	}

	return result, nil
}

func (c *converter) convertJsonValue(node yaccNode) (common.Value, error) {
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
		return c.convertJsonArray(node)
	case NodeOpJsonObject:
		return c.convertJsonObject(node)
	case NodeOpJsonIdentifier:
		var valStr = node.value.(string)
		switch valStr {
		case "null":
			return common.NullValue(), nil
		default:
			return common.StringValue(valStr), fmt.Errorf("unexpected json identifier: %s", valStr)
		}
	default:
		panic(fmt.Sprintf("unexpected node op: %s", node.op))
	}
}

func (c *converter) convertArgumentList(element yaccNode) (*plain.DefinitionStatementElementArgumentList, error) {
	argumentList := new(plain.DefinitionStatementElementArgumentList)

	for _, child := range element.children {
		switch child.op {
		case NodeOpArgument:
			argument, err := c.convertArgument(child)
			if err != nil {
				return argumentList, fmt.Errorf("failed to convert argument: %w", err)
			}
			argumentList.Arguments = append(argumentList.Arguments, *argument)
		}
	}

	return argumentList, nil
}

func (c *converter) convertArgument(element yaccNode) (*plain.DefinitionStatementElementArgument, error) {
	argument := new(plain.DefinitionStatementElementArgument)

	argument.Name = element.value.(string)

	if len(element.children) > 0 {
		typeDef, err := c.convertTypeDefinition(element.children[0])
		if err != nil {
			return argument, fmt.Errorf("failed to convert type def: %w", err)
		}
		argument.Type = *typeDef
	}

	return argument, nil
}

func (c *converter) convertParameterList(element yaccNode) (*plain.DefinitionStatementElementParameterList, error) {
	parameterList := new(plain.DefinitionStatementElementParameterList)

	for _, child := range element.children {
		expr, err := c.convertExpression(child)
		if err != nil {
			return parameterList, fmt.Errorf("failed to convert parameter: %w", err)
		}
		parameterList.Parameters = append(parameterList.Parameters, *expr)
	}

	return parameterList, nil
}

func (c *converter) statementToValue(node yaccNode) (*common.Value, error) {
	statement, err := c.convertDefinitionStatementElement(node)

	value := statement.AsValue()

	if err != nil {
		return &value, fmt.Errorf("failed to convert statement: %w", err)
	}

	return &value, nil
}

func (c *converter) convertTypeDefinition(node yaccNode) (*common.TypeDefinition, error) {
	var result = new(common.TypeDefinition)
	result.Name = node.value.(string)

	if len(node.children) > 0 {
		for _, child := range node.children {
			subType, err := c.convertTypeDefinition(child)

			if err != nil {
				return nil, fmt.Errorf("failed to convert type definition: %w", err)
			}

			result.SubTypes = append(result.SubTypes, *subType)
		}
	}

	return result, nil
}

func (c *converter) convertIdentifier(element yaccNode) (*plain.DefinitionStatementElementIdentifier, error) {
	identifier := new(plain.DefinitionStatementElementIdentifier)

	identifier.Identifier = element.value.(string)

	return identifier, nil
}

func (c *converter) convertAttributeList(element yaccNode) (*plain.DefinitionStatementElementAttributeList, error) {
	attributeList := new(plain.DefinitionStatementElementAttributeList)

	for _, child := range element.children {
		switch child.op {
		case NodeOpAttribute:
			attribute, err := c.convertAttribute(child)
			if err != nil {
				return nil, fmt.Errorf("failed to convert attribute: %w", err)
			}
			attributeList.Attributes = append(attributeList.Attributes, *attribute)
		}
	}

	return attributeList, nil

}

func (c *converter) convertAttribute(element yaccNode) (*plain.DefinitionStatementElementAttribute, error) {
	attribute := new(plain.DefinitionStatementElementAttribute)

	attribute.Name = element.value.(string)

	if len(element.children) > 0 {
		value, err := c.convertValue(element.children[0])
		if err != nil {
			return nil, fmt.Errorf("failed to convert value: %w", err)
		}
		attribute.Value = value
	}

	return attribute, nil

}

func (c *converter) convertValue(element yaccNode) (*common.Value, error) {
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

func (c *converter) convertArray(element yaccNode) (*plain.DefinitionStatementElementArray, error) {
	array := new(plain.DefinitionStatementElementArray)

	for _, child := range element.children {
		st, err := c.convertDefinitionStatementElement(child)

		if err != nil {
			return nil, err
		}

		array.Items = append(array.Items, *st)
	}

	return array, nil
}

func (c *converter) convertStruct(element yaccNode) (*plain.DefinitionStatementElementStruct, error) {
	result := new(plain.DefinitionStatementElementStruct)

	for _, child := range element.children[0].children[0].children {
		st, err := c.convertDefinitionStatementElement(child)

		if err != nil {
			return nil, err
		}

		result.Statements = append(result.Statements, *st)
	}

	return result, nil
}
