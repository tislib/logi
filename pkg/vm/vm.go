package vm

import (
	"github.com/tislib/logi/pkg/ast/common"
	logiAst "github.com/tislib/logi/pkg/ast/logi"
	macroAst "github.com/tislib/logi/pkg/ast/macro"
)

type VM struct {
	Macros []macroAst.Ast
}

func NewVM() *VM {
	return &VM{}
}

func (v *VM) Run(ast logiAst.Ast) (*common.Value, error) {
	// locate main function

	for _, function := range ast.Functions {
		if function.Name == "main" {
			return v.runFunction(function)
		}
	}

	return nil, nil
}

func (v *VM) runFunction(function logiAst.Function) (*common.Value, error) {
	for _, statement := range function.CodeBlock.Statements {
		value, err := v.runStatement(statement)

		if err != nil {
			return nil, err
		}

		if value != nil {
			return value, nil
		}
	}

	return nil, nil
}

func (v *VM) runStatement(statement common.Statement) (*common.Value, error) {
	switch statement.Kind {
	case common.FuncCallStatementKind:
		return v.runFuncCallStatement(statement.FuncCall)
	}

	return nil, nil
}

func (v *VM) runFuncCallStatement(call *common.FunctionCallStatement) (*common.Value, error) {

	switch call.Call.Name {
	case "println":
		for _, arg := range call.Call.Arguments {
			value, err := v.runExpression(arg)

			if err != nil {
				return nil, err
			}

			println(value.ToDisplayName())
		}
	}
	return nil, nil
}

func (v *VM) runExpression(arg *common.Expression) (*common.Value, error) {
	switch arg.Kind {
	case common.LiteralKind:
		return &arg.Literal.Value, nil
	}

	return nil, nil
}
