package interpreter

import (
	"fmt"
	"github.com/tislib/logi/pkg/ast/common"
	logiAst "github.com/tislib/logi/pkg/ast/logi"
	"github.com/tislib/logi/pkg/vm"
)

func (i *Interpreter) RunSourceFile(file string) (*common.Value, error) {
	ast, err := i.parse(file)

	if err != nil {
		return nil, fmt.Errorf("error parsing %s: %w", file, err)
	}

	return i.runAst(ast)
}

func (i *Interpreter) runAst(ast *logiAst.Ast) (*common.Value, error) {
	instance := vm.NewVM()
	instance.Macros = i.macros

	return instance.Run(*ast)
}
