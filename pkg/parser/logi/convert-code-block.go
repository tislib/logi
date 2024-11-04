package logi

import (
	"fmt"
	"github.com/tislib/logi/pkg/ast/common"
)

func convertCodeBlock(element yaccNode) (*common.CodeBlock, error) {
	codeBlock := new(common.CodeBlock)

	var statementsElement = element.children[0].children

	for _, child := range statementsElement {
		switch child.op {
		case NodeOpIf:
			ifStatement, err := convertIfStatement(child)

			if err != nil {
				return nil, err
			}

			codeBlock.Statements = append(codeBlock.Statements, common.Statement{
				Kind: common.IfStatementKind,

				IfStmt: ifStatement,
			})
		case NodeOpReturn:
			returnStatement, err := convertReturnStatement(child)

			if err != nil {
				return nil, err
			}

			codeBlock.Statements = append(codeBlock.Statements, common.Statement{
				Kind: common.ReturnStatementKind,

				ReturnStmt: returnStatement,
			})

		case NodeOpVariableDeclaration:
			varDeclaration, err := convertVariableDeclaration(child)

			if err != nil {
				return nil, err
			}

			codeBlock.Statements = append(codeBlock.Statements, common.Statement{
				Kind: common.VarDeclKind,

				VarDecl: varDeclaration,
			})
		case NodeOpFunctionCall:
			functionCall, err := convertFunctionCallStatement(child)

			if err != nil {
				return nil, err
			}

			codeBlock.Statements = append(codeBlock.Statements, common.Statement{
				Kind: common.FuncCallStatementKind,

				FuncCall: functionCall,
			})
		default:
			return nil, fmt.Errorf("unexpected node op: %s", child.op)
		}
	}

	return codeBlock, nil
}

func convertReturnStatement(element yaccNode) (*common.ReturnStatement, error) {
	returnStatement := new(common.ReturnStatement)

	expression, err := convertExpression(element.children[0])

	if err != nil {
		return nil, err
	}

	returnStatement.Result = expression

	return returnStatement, nil
}

func convertVariableDeclaration(child yaccNode) (*common.VarDeclaration, error) {
	varDeclaration := new(common.VarDeclaration)

	varDeclaration.Name = child.children[0].value.(string)

	typeDef, err := convertTypeDefinition(child.children[0])

	if err != nil {
		return nil, err
	}

	varDeclaration.Type = *typeDef

	if len(child.children) > 2 {
		expression, err := convertExpression(child.children[1])

		if err != nil {
			return nil, err
		}

		varDeclaration.Value = expression
	}

	return varDeclaration, nil

}

func convertIfStatement(element yaccNode) (*common.IfStatement, error) {
	ifStatement := new(common.IfStatement)

	condition, err := convertExpression(element.children[0])

	if err != nil {
		return nil, err
	}

	ifStatement.Condition = condition

	codeBlock, err := convertCodeBlock(element.children[1])

	if err != nil {
		return nil, err
	}

	ifStatement.ThenBlock = codeBlock

	if len(element.children) > 2 {
		elseBlock, err := convertCodeBlock(element.children[2])

		if err != nil {
			return nil, err
		}

		ifStatement.ElseBlock = elseBlock
	}

	return ifStatement, nil
}

func convertFunctionCallStatement(element yaccNode) (*common.FunctionCallStatement, error) {
	functionCall := new(common.FunctionCallStatement)

	expression, err := convertExpression(element.children[0])

	if err != nil {
		return nil, err
	}

	functionCall.Call = expression.FuncCall

	return functionCall, nil
}

func convertExpression(element yaccNode) (*common.Expression, error) {
	expression := new(common.Expression)

	switch element.op {
	case NodeOpBinaryExpression:
		binaryExpression, err := convertBinaryExpression(element)

		if err != nil {
			return nil, err
		}

		expression.Kind = common.BinaryExprKind
		expression.BinaryExpr = binaryExpression
	case NodeOpLiteral:
		literal, err := convertValue(element)

		if err != nil {
			return nil, err
		}

		expression.Kind = common.LiteralKind
		expression.Literal = &common.Literal{Value: *literal}
	case NodeOpVariable:
		variable, err := convertVariable(element)

		if err != nil {
			return nil, err
		}

		expression.Kind = common.VariableKind
		expression.Variable = variable
	case NodeOpFunctionCall:
		functionCall, err := convertFunctionCall(element)

		if err != nil {
			return nil, err
		}

		expression.Kind = common.FuncCallKind
		expression.FuncCall = functionCall
	default:
		return nil, fmt.Errorf("unexpected node op: %s", element.op)
	}

	return expression, nil
}

func convertFunctionCall(element yaccNode) (*common.FunctionCall, error) {
	functionCall := new(common.FunctionCall)

	functionCall.Name = element.value.(string)

	if len(element.children) > 0 {
		arguments, err := convertArguments(element.children[0])

		if err != nil {
			return nil, err
		}

		functionCall.Arguments = arguments
	}

	return functionCall, nil
}

func convertArguments(element yaccNode) ([]*common.Expression, error) {
	var arguments []*common.Expression

	for _, child := range element.children {
		argument, err := convertExpression(child)

		if err != nil {
			return nil, err
		}

		arguments = append(arguments, argument)
	}

	return arguments, nil
}

func convertBinaryExpression(element yaccNode) (*common.BinaryExpression, error) {
	binaryExpression := new(common.BinaryExpression)

	left, err := convertExpression(element.children[0])

	if err != nil {
		return nil, err
	}

	binaryExpression.Left = left

	right, err := convertExpression(element.children[1])

	if err != nil {
		return nil, err
	}

	binaryExpression.Right = right

	binaryExpression.Operator = element.children[2].value.(string)

	return binaryExpression, nil
}

func convertVariable(element yaccNode) (*common.Variable, error) {
	variable := new(common.Variable)

	variable.Name = element.value.(string)

	return variable, nil
}
