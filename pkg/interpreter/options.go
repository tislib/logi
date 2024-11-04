package interpreter

import "github.com/tislib/logi/pkg/registry"

func WithBaseDirectory(baseDirectory string) Option {
	return func(i *Interpreter) {
		i.baseDirectory = baseDirectory
	}
}

func WithRegistry(registry *registry.Registry) Option {
	return func(i *Interpreter) {
		i.registry = registry
	}
}
