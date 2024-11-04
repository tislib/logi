package interpreter

import (
	macroAst "github.com/tislib/logi/pkg/ast/macro"
	"github.com/tislib/logi/pkg/registry"
	"log"
	"os"
)

type Interpreter struct {
	baseDirectory string
	registry      *registry.Registry
	macros        []macroAst.Ast
}

type Option func(*Interpreter)

func NewInterpreter(opts ...Option) *Interpreter {
	interpreter := &Interpreter{
		baseDirectory: currentDirectory(),
		registry:      registry.NewRegistry(),
	}
	for _, opt := range opts {
		opt(interpreter)
	}

	interpreter.registry.RegisterResolver("local", interpreter.resolver)

	return interpreter
}

func (i *Interpreter) resolver(fileName string) (string, error) {
	content, err := os.ReadFile(fileName)

	if err != nil {
		return "", err
	}

	return string(content), nil
}

func currentDirectory() string {
	wd, err := os.Getwd()

	if err != nil {
		log.Fatalf("error getting current directory: %v", err)
	}

	return wd
}
