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
	Logis           []logiAst.Ast
	Definitions     []Definition
	locals          map[string]interface{}
	vars            map[string]interface{}
	types           map[string]common.TypeDefinition
	enableSourceMap bool
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

func (v *vm) prepareDefinition(ast logiAst.Definition) (Definition, error) {
	var definition = Definition{}

	definition.Name = ast.Name
	definition.Macro = ast.MacroName
	definition.Data = make(map[string]map[string]interface{})

	for key, value := range ast.Dynamic {
		definition.Data[key] = make(map[string]interface{})

		for dk, dv := range value {
			definition.Data[key][dk] = dv

			if dk == "code" {
				definition.Data[key]["exec"] = v.executableFunc(definition.Data[key][dk].(common.CodeBlock))
			}
		}
	}

	return definition, nil
}

func (v *vm) MapToStruct(definition Definition) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (v *vm) GetDefinitionByName(name string) (*Definition, error) {
	for _, definition := range v.Definitions {
		if definition.Name == name {
			return &definition, nil
		}
	}

	return nil, fmt.Errorf("definition %s not found", name)
}

type Option func(vm *vm) error

func WithLocals(locals map[string]interface{}) Option {
	return func(vm *vm) error {
		for k, v := range locals {
			vm.locals[k] = v
		}

		return nil
	}
}

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

func (v *vm) SetLocal(key string, value interface{}) {
	v.locals[key] = value
}
