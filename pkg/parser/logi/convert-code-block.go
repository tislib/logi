package logi

import (
	"fmt"
	"github.com/tislib/logi/pkg/ast/common"
)

func (c *converter) convertCodeBlock(element yaccNode) (*common.CodeBlock, error) {
	codeBlock := new(common.CodeBlock)

	var statementsElement = element.children[0].children

	for _, child := range statementsElement {
		switch child.op {
		case NodeOpIf:
			ifStatement, err := c.convertIfStatement(child)

			if err != nil {
				return nil, err
			}

			codeBlock.Statements = append(codeBlock.Statements, common.Statement{
				Kind: common.IfStatementKind,

				IfStmt: ifStatement,
			})
		case NodeOpReturn:
			returnStatement, err := c.convertReturnStatement(child)

			if err != nil {
				return nil, err
			}

			codeBlock.Statements = append(codeBlock.Statements, common.Statement{
				Kind: common.ReturnStatementKind,

				ReturnStmt: returnStatement,
			})

		case NodeOpVariableDeclaration:
			varDeclaration, err := c.convertVariableDeclaration(child)

			if err != nil {
				return nil, err
			}

			codeBlock.Statements = append(codeBlock.Statements, common.Statement{
				Kind: common.VarDeclKind,

				VarDecl: varDeclaration,
			})
		case NodeOpFunctionCall:
			functionCall, err := c.convertFunctionCallStatement(child)

			if err != nil {
				return nil, err
			}

			codeBlock.Statements = append(codeBlock.Statements, common.Statement{
				Kind: common.FuncCallStatementKind,

				FuncCall: functionCall,
			})
		case NodeOpExpression:
			expression, err := c.convertExpression(child)

			if err != nil {
				return nil, err
			}

			codeBlock.Statements = append(codeBlock.Statements, common.Statement{
				Kind: common.ExprStatementKind,
				ExprStmt: &common.ExpressionStatement{
					Expr: expression,
				},
			})
		default:
			return nil, fmt.Errorf("unexpected node op: %s", child.op)
		}
	}

	return codeBlock, nil
}

func (c *converter) convertReturnStatement(element yaccNode) (*common.ReturnStatement, error) {
	returnStatement := new(common.ReturnStatement)

	expression, err := c.convertExpression(element.children[0])

	if err != nil {
		return nil, err
	}

	returnStatement.Result = expression

	return returnStatement, nil
}

func (c *converter) convertVariableDeclaration(child yaccNode) (*common.VarDeclaration, error) {
	varDeclaration := new(common.VarDeclaration)

	varDeclaration.Name = child.children[0].value.(string)

	typeDef, err := c.convertTypeDefinition(child.children[0])

	if err != nil {
		return nil, err
	}

	varDeclaration.Type = *typeDef

	if len(child.children) > 2 {
		expression, err := c.convertExpression(child.children[1])

		if err != nil {
			return nil, err
		}

		varDeclaration.Value = expression
	}

	return varDeclaration, nil

}

func (c *converter) convertIfStatement(element yaccNode) (*common.IfStatement, error) {
	ifStatement := new(common.IfStatement)

	condition, err := c.convertExpression(element.children[0])

	if err != nil {
		return nil, err
	}

	ifStatement.Condition = condition

	codeBlock, err := c.convertCodeBlock(element.children[1])

	if err != nil {
		return nil, err
	}

	ifStatement.ThenBlock = codeBlock

	if len(element.children) > 2 {
		elseBlock, err := c.convertCodeBlock(element.children[2])

		if err != nil {
			return nil, err
		}

		ifStatement.ElseBlock = elseBlock
	}

	return ifStatement, nil
}

func (c *converter) convertFunctionCallStatement(element yaccNode) (*common.FunctionCallStatement, error) {
	functionCall := new(common.FunctionCallStatement)

	expression, err := c.convertExpression(element.children[0])

	if err != nil {
		return nil, err
	}

	functionCall.Call = expression.FuncCall

	return functionCall, nil
}

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
