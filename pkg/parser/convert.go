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
		baseMacro, err := convertBaseMacro(child)

		if err != nil {
			return nil, fmt.Errorf("failed to convert base macro: %w", err)
		}

		switch baseMacro.Kind {
		case ast.MacroKindSyntax:
			syntaxMacro, err := convertSyntaxMacro(child, baseMacro)

			if err != nil {
				return nil, fmt.Errorf("failed to convert syntax macro: %w", err)
			}
			res.SyntaxMacros = append(res.SyntaxMacros, *syntaxMacro)
		default:
			return nil, fmt.Errorf("unexpected macro kind: %s", baseMacro.Kind)
		}

	}

	return res, nil
}

func convertSyntaxMacro(child yaccMacroNode, macro *ast.BaseMacro) (*ast.SyntaxMacro, error) {
	var result = new(ast.SyntaxMacro)

	result.Name = macro.Name
	result.Kind = macro.Kind

	return result, nil
}

func convertBaseMacro(child yaccMacroNode) (*ast.BaseMacro, error) {
	var signature = child.children[0]
	var name = signature.children[0]
	var body = signature.children[1]
	var kind = body.children[0].value.(string)

	var result = new(ast.BaseMacro)
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
