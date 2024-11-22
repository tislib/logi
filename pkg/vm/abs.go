package vm

import (
	logiAst "github.com/tislib/logi/pkg/ast/logi"
	macroAst "github.com/tislib/logi/pkg/ast/macro"
)

type Definition struct {
	Macro string                            `json:"macro"`
	Name  string                            `json:"name"`
	Data  map[string]map[string]interface{} `json:"data"`
}

type ExecutableFunc func(args ...interface{}) (interface{}, error)

type VirtualMachine interface {
	// loads macro files from the given paths
	LoadMacroFile(path ...string) error
	LoadMacroContent(content ...string) error
	LoadMacroAst(ast ...macroAst.Ast) error

	// loads logi files from the given paths
	LoadLogiFile(path ...string) ([]Definition, error)
	LoadLogiContent(content ...string) ([]Definition, error)
	LoadLogiAst(ast ...logiAst.Ast) ([]Definition, error)

	MapToStruct(definition Definition) (string, error)

	LocateCodeBlock(definition Definition, codeBlockPath string) (ExecutableFunc, error)
	GetDefinitionByName(name string) (*Definition, error)
	SetLocals(locals map[string]interface{})
	GetLocals() map[string]interface{}
	Execute(def *Definition, s string) (interface{}, error)
}
