package vm

import (
	logiAst "github.com/tislib/logi/pkg/ast/logi"
)

type simpleImplementer struct {
	fn ImplementerFunc
}

func (s simpleImplementer) Call(vm VirtualMachine, statement logiAst.Statement) error {
	return s.fn(vm, statement, func(next logiAst.Statement) error {
		return s.Call(vm, next)
	})
}

func NewImplementerFunc(fn ImplementerFunc) Implementer {
	return &simpleImplementer{fn: fn}
}
