package interpreter

import (
	"fmt"
	logiAst "github.com/tislib/logi/pkg/ast/logi"
	"github.com/tislib/logi/pkg/parser/logi"
)

func (i *Interpreter) parse(file string) (*logiAst.Ast, error) {
	content, err := i.registry.ImportFile(i.baseDirectory + "/" + file)

	if err != nil {
		return nil, fmt.Errorf("error importing %s: %w", file, err)
	}

	ast, err := logi.Parse(content, i.macros)

	if err != nil {
		return nil, fmt.Errorf("error parsing %s: %w", file, err)
	}

	return ast, nil
}
