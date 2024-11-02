package parser

import (
	"fmt"
	"logi/pkg/ast"
)

func convertNodeToMacroAst(node yaccMacroNode) (*ast.MacroAst, error) {
	var res = new(ast.MacroAst)

	if node.op != NodeOpFile {
		return nil, ErrUnexpectedNode
	}

	for _, child := range node.children {
		macro, err := convertMacro(child)

		if err != nil {
			return nil, fmt.Errorf("failed to convert syntax macro: %w", err)
		}

		res.Macros = append(res.Macros, *macro)
	}

	return res, nil
}

func convertMacro(child yaccMacroNode) (*ast.Macro, error) {
	var signature = child.children[0]
	var name = signature.children[0]
	var body = signature.children[1]
	var kind = body.children[0].value.(string)

	var result = new(ast.Macro)
	if !NamePattern.MatchString(name.value.(string)) {
		return nil, fmt.Errorf("unexpected name value: %s", name.value)
	}

	result.Name = name.value.(string)
	switch kind {
	case "Syntax":
		result.Kind = ast.MacroKindSyntax
	default:
		return nil, fmt.Errorf("unexpected kind value: %s", kind)
	}

	return result, nil
}
