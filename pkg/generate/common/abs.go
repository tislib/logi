package common

import macroAst "github.com/tislib/logi/pkg/ast/macro"

type CodeGenerator interface {
	Generate(ast *macroAst.Ast, pkg string) (string, error)
}
