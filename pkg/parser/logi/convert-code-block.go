package logi

import (
	"fmt"
	"logi/pkg/ast/plain"
)

func convertCodeBlock(element yaccNode) (*plain.CodeBlock, error) {
	codeBlock := new(plain.CodeBlock)

	var statementsElement = element.children[0].children

	for _, child := range statementsElement {
		switch child.op {
		case NodeOpIf:
			ifStatement, err := convertIfStatement(child)

			if err != nil {
				return nil, err
			}

			codeBlock.Statements = append(codeBlock.Statements, plain.Statement{
				Kind: plain.IfStatementKind,

				IfStmt: ifStatement,
			})
		case NodeOpReturn:
			returnStatement, err := convertReturnStatement(child)

			if err != nil {
				return nil, err
			}

			codeBlock.Statements = append(codeBlock.Statements, plain.Statement{
				Kind: plain.ReturnStatementKind,

				ReturnStmt: returnStatement,
			})

		case NodeOpVariableDeclaration:
			varDeclaration, err := convertVariableDeclaration(child)

			if err != nil {
				return nil, err
			}

			codeBlock.Statements = append(codeBlock.Statements, plain.Statement{
				Kind: plain.VarDeclKind,

				VarDecl: varDeclaration,
			})

		default:
			return nil, fmt.Errorf("unexpected node op: %s", child.op)
		}
	}

	return codeBlock, nil
}

func convertReturnStatement(element yaccNode) (*plain.ReturnStatement, error) {
	returnStatement := new(plain.ReturnStatement)

	expression, err := convertExpression(element.children[0])

	if err != nil {
		return nil, err
	}

	returnStatement.Result = expression

	return returnStatement, nil
}

func convertVariableDeclaration(child yaccNode) (*plain.VarDeclaration, error) {
	varDeclaration := new(plain.VarDeclaration)

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

func convertIfStatement(element yaccNode) (*plain.IfStatement, error) {
	ifStatement := new(plain.IfStatement)

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

func convertExpression(element yaccNode) (*plain.Expression, error) {
	expression := new(plain.Expression)

	switch element.op {
	case NodeOpBinaryExpression:
		binaryExpression, err := convertBinaryExpression(element)

		if err != nil {
			return nil, err
		}

		expression.Kind = plain.BinaryExprKind
		expression.BinaryExpr = binaryExpression
	case NodeOpLiteral:
		literal, err := convertValue(element)

		if err != nil {
			return nil, err
		}

		expression.Kind = plain.LiteralKind
		expression.Literal = &plain.Literal{Value: *literal}
	case NodeOpVariable:
		variable, err := convertVariable(element)

		if err != nil {
			return nil, err
		}

		expression.Kind = plain.VariableKind
		expression.Variable = variable
	default:
		return nil, fmt.Errorf("unexpected node op: %s", element.op)
	}

	return expression, nil
}

func convertBinaryExpression(element yaccNode) (*plain.BinaryExpression, error) {
	binaryExpression := new(plain.BinaryExpression)

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

func convertVariable(element yaccNode) (*plain.Variable, error) {
	variable := new(plain.Variable)

	variable.Name = element.value.(string)

	return variable, nil
}
