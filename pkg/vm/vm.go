package vm

import (
	"fmt"
	"github.com/tislib/logi/pkg/ast/common"
	logiAst "github.com/tislib/logi/pkg/ast/logi"
	macroAst "github.com/tislib/logi/pkg/ast/macro"
)

type vm struct {
	Macros      []macroAst.Macro
	Logis       []logiAst.Ast
	Definitions []Definition
	locals      map[string]interface{}
	vars        map[string]interface{}
	types       map[string]common.TypeDefinition
}

func (v *vm) SetLocals(locals map[string]interface{}) {
	for key, value := range locals {
		v.locals[key] = value
	}
}

func (v *vm) GetLocals() map[string]interface{} {
	return v.locals
}

func (v *vm) prepareDefinition(plainAst logiAst.Definition) (Definition, error) {
	var result = Definition{}

	result.Name = plainAst.Name
	result.Macro = plainAst.MacroName
	result.Data = make(map[string]interface{})

	for _, param := range plainAst.Parameters {
		v.mapParameters(param, &result)
		for _, attribute := range param.Attributes {
			result.Data[param.Name+"/"+attribute.Name] = attribute.Value.AsInterface()
		}

		if param.CodeBlock != nil {
			result.Data[param.Name] = v.executableFunc(*param.CodeBlock)
		}
	}

	for _, method := range plainAst.Methods {
		result.Data[method.Name] = v.executableFunc(method.CodeBlock)
		for _, parameter := range method.Parameters {
			result.Data[method.Name+"/"+parameter.Name] = parameter.Value.AsInterface()
		}
		for _, attribute := range method.Attributes {
			result.Data[method.Name+"/"+attribute.Name] = attribute.Value.AsInterface()
		}
		for _, argument := range method.Arguments {
			result.Data[method.Name+"/"+argument.Name] = argument.Type.AsValue().AsInterface()
		}
	}

	for _, methodSignature := range plainAst.MethodSignature {
		result.Data[methodSignature.Name] = nil
		for _, parameter := range methodSignature.Parameters {
			result.Data[methodSignature.Name+"/"+parameter.Name] = parameter.Value.AsInterface()
		}
		for _, attribute := range methodSignature.Attributes {
			result.Data[methodSignature.Name+"/"+attribute.Name] = attribute.Value.AsInterface()
		}
		for _, argument := range methodSignature.Arguments {
			result.Data[methodSignature.Name+"/"+argument.Name] = argument.Type.AsValue().AsInterface()
		}
	}

	return result, nil
}

func (v *vm) mapParameters(param logiAst.DefinitionParameter, result *Definition) {
	for _, parameter := range param.Parameters {
		if param.Name == parameter.Name {
			result.Data[parameter.Name] = parameter.Value.AsInterface()
		} else {
			result.Data[param.Name+"/"+parameter.Name] = parameter.Value.AsInterface()
		}
	}
}

func (v *vm) MapToStruct(definition Definition) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (v *vm) LocateCodeBlock(definition Definition, codeBlockPath string) (ExecutableFunc, error) {
	if definition.Data[codeBlockPath] == nil {
		return nil, fmt.Errorf("code block %s not found in definition %s", codeBlockPath, definition.Name)
	}

	fn, ok := definition.Data[codeBlockPath].(ExecutableFunc)

	if !ok {
		return nil, fmt.Errorf("code block %s is not a function in definition %s", codeBlockPath, definition.Name)
	}

	return fn, nil
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
		locals: make(map[string]interface{}),
		vars:   make(map[string]interface{}),
	}

	for _, o := range option {
		if err := o(v); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return v, nil
}
