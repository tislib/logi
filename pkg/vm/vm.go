package vm

import (
	"fmt"
	"github.com/tislib/logi/pkg/ast/common"
	logiAst "github.com/tislib/logi/pkg/ast/logi"
	macroAst "github.com/tislib/logi/pkg/ast/macro"
)

type vm struct {
	Macros          []macroAst.Macro
	MacroContents   map[string]string
	Definitions     []logiAst.Definition
	locals          map[string]interface{}
	vars            map[string]interface{}
	types           map[string]common.TypeDefinition
	enableSourceMap bool
}

func (v *vm) GetMacros() []macroAst.Macro {
	return v.Macros
}

func (v *vm) GetMacroContent(name string) string {
	return v.MacroContents[name]
}

func (v *vm) SetLocals(locals map[string]interface{}) {
	for key, value := range locals {
		v.locals[key] = value
	}
}

func (v *vm) GetLocals() map[string]interface{} {
	return v.locals
}

func (v *vm) MapToStruct(definition logiAst.Definition) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (v *vm) GetDefinitionByName(name string) (*logiAst.Definition, error) {
	for _, definition := range v.Definitions {
		if definition.Name == name {
			return &definition, nil
		}
	}

	return nil, fmt.Errorf("logiAst.Definition %s not found", name)
}

type Option func(vm *vm) error

func New(option ...Option) (VirtualMachine, error) {
	v := &vm{
		locals:        make(map[string]interface{}),
		vars:          make(map[string]interface{}),
		MacroContents: make(map[string]string),
	}

	for _, o := range option {
		if err := o(v); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return v, nil
}
