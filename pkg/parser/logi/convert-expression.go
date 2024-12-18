package logi

import (
	"fmt"
	"github.com/tislib/logi/pkg/ast/common"
)

func (c *converter) convertExpression(element yaccNode) (*common.Expression, error) {
	expression := new(common.Expression)

	switch element.op {
	case NodeOpBinaryExpression:
		binaryExpression, err := c.convertBinaryExpression(element)

		if err != nil {
			return nil, err
		}

		expression.Kind = common.BinaryExprKind
		expression.BinaryExpr = binaryExpression
	case NodeOpLiteral:
		literal, err := c.convertValue(element)

		if err != nil {
			return nil, err
		}

		expression.Kind = common.LiteralKind
		expression.Literal = &common.Literal{Value: *literal}
	case NodeOpVariable:
		variable, err := c.convertVariable(element)

		if err != nil {
			return nil, err
		}

		expression.Kind = common.VariableKind
		expression.Variable = variable
	case NodeOpFunctionCall:
		functionCall, err := c.convertFunctionCall(element)

		if err != nil {
			return nil, err
		}

		expression.Kind = common.FuncCallKind
		expression.FuncCall = functionCall
	case NodeOpExpression:
		subExpression, err := c.convertExpression(element.children[0])

		if err != nil {
			return nil, err
		}

		return subExpression, nil
	default:
		return nil, fmt.Errorf("unexpected node op: %s", element.op)
	}

	return expression, nil
}

func (c *converter) convertFunctionCall(element yaccNode) (*common.FunctionCall, error) {
	functionCall := new(common.FunctionCall)

	functionCall.Name = element.value.(string)

	if len(element.children) > 0 {
		arguments, err := c.convertArguments(element.children[0])

		if err != nil {
			return nil, err
		}

		functionCall.Arguments = arguments
	}

	return functionCall, nil
}

func (c *converter) convertArguments(element yaccNode) ([]*common.Expression, error) {
	var arguments []*common.Expression

	for _, child := range element.children {
		argument, err := c.convertExpression(child)

		if err != nil {
			return nil, err
		}

		arguments = append(arguments, argument)
	}

	return arguments, nil
}

func (c *converter) convertBinaryExpression(element yaccNode) (*common.BinaryExpression, error) {
	binaryExpression := new(common.BinaryExpression)

	left, err := c.convertExpression(element.children[0])

	if err != nil {
		return nil, err
	}

	binaryExpression.Left = left

	right, err := c.convertExpression(element.children[1])

	if err != nil {
		return nil, err
	}

	binaryExpression.Right = right

	binaryExpression.Operator = element.children[2].value.(string)

	return binaryExpression, nil
}

func (c *converter) convertVariable(element yaccNode) (*common.Variable, error) {
	variable := new(common.Variable)

	variable.Name = element.value.(string)

	return variable, nil
}
