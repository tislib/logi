package vm

import (
	"github.com/tislib/logi/pkg/ast/common"
	logiAst "github.com/tislib/logi/pkg/ast/logi"
	macroAst "github.com/tislib/logi/pkg/ast/macro"
)

type ImplementerFunc func(vm VirtualMachine, statement logiAst.Statement, passNext func(statement logiAst.Statement) error) error

type Implementer interface {
	Call(vm VirtualMachine, statement logiAst.Statement) error
}

type VirtualMachine interface {
	// loads macro files from the given paths
	LoadMacroFile(path ...string) error
	LoadMacroContent(content ...string) error
	LoadMacroAst(ast ...macroAst.Ast) error

	// loads logi files from the given paths
	LoadLogiFile(path ...string) ([]logiAst.Definition, error)
	LoadLogiContent(content ...string) ([]logiAst.Definition, error)
	LoadLogiAst(ast ...logiAst.Ast) ([]logiAst.Definition, error)
	GetMacroContent(name string) string
	GetMacros() []macroAst.Macro
	GetDefinitionByName(name string) (*logiAst.Definition, error)

	// VM functions
	Execute(def *logiAst.Definition, implementer Implementer) error
	Evaluate(expression common.Expression, vars map[string]common.Value, fns map[string]func(args ...common.Value) (common.Value, error)) (common.Value, error)
}
