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

func convertMacro(macroNode yaccMacroNode) (*ast.Macro, error) {
	var signature = macroNode.children[0]
	var name = signature.children[0]
	var body = macroNode.children[1]
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

	for _, child := range body.children {
		switch child.op {
		case NodeOpSyntax:
			if len(child.children) != 0 {
				syntaxBody, err := convertSyntaxBody(child.children[0])

				if err != nil {
					return nil, fmt.Errorf("failed to convert syntax: %w", err)
				}

				result.Syntax = ast.Syntax{Statements: syntaxBody}
			}
		case NodeOpDefinition:
			if len(child.children) != 0 {
				syntaxBody, err := convertSyntaxBody(child.children[0])

				if err != nil {
					return nil, fmt.Errorf("failed to convert definition: %w", err)
				}

				result.Definition = ast.Definition{Statements: syntaxBody}
			}
		}
	}

	return result, nil
}

func convertSyntaxBody(syntaxNode yaccMacroNode) ([]ast.SyntaxStatement, error) {
	if syntaxNode.children == nil {
		return nil, nil
	}

	var result []ast.SyntaxStatement

	for _, child := range syntaxNode.children {
		statement, err := convertSyntaxStatement(child)

		if err != nil {
			return nil, fmt.Errorf("failed to convert syntax statement: %w", err)
		}

		result = append(result, *statement)
	}

	return result, nil
}

func convertSyntaxStatement(child yaccMacroNode) (*ast.SyntaxStatement, error) {
	var result = new(ast.SyntaxStatement)

	for _, elementNode := range child.children {
		element, err := convertSyntaxStatementElement(elementNode)

		if err != nil {
			return nil, fmt.Errorf("failed to convert syntax statement element: %w", err)
		}

		result.Elements = append(result.Elements, *element)
	}

	return result, nil
}

func convertSyntaxStatementElement(node yaccMacroNode) (*ast.SyntaxStatementElement, error) {
	var result = new(ast.SyntaxStatementElement)

	switch node.op {
	case NodeOpSyntaxKeywordElement:
		result.Kind = ast.SyntaxStatementElementKindKeyword
		result.KeywordDef = &ast.SyntaxStatementElementKeywordDef{Name: node.value.(string)}
	default:
		return nil, fmt.Errorf("unexpected syntax statement element op: %s", node.op)
	}

	return result, nil
}
