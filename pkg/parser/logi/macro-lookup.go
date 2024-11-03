package logi

import (
	"fmt"
	"logi/pkg/ast/common"
	macroAst "logi/pkg/ast/macro"
	"logi/pkg/ast/plain"
)

func locateMacroDefinition(definition plain.Definition, ast macroAst.Ast) (*macroAst.Macro, error) {
	for _, macroDefinition := range ast.Macros {
		if macroDefinition.Name == definition.MacroName {
			return &macroDefinition, nil
		}
	}

	return nil, fmt.Errorf("macro definition not found")

}

func locateMacroSyntaxStatement(statement plain.DefinitionStatement, macro *macroAst.Macro) (*macroAst.SyntaxStatement, []int, error) {
	maxMatch := 0
	var bestMatch *macroAst.SyntaxStatement
	var mismatchCause string
	var syntaxElementMatch []int

	var reportMismatch = func(matchedUntil int, syntaxStatement macroAst.SyntaxStatement, reason string) {
		if matchedUntil > maxMatch {
			maxMatch = matchedUntil
			bestMatch = &syntaxStatement
			mismatchCause = reason
		}
	}

	for _, syntaxStatement := range macro.Syntax.Statements {
		ei := 0
		match := len(syntaxStatement.Elements) > 0
		syntaxElementMatch = make([]int, len(statement.Elements))

		for i, syntaxStatementElement := range syntaxStatement.Elements {
			// check if the statement is always required
			alwaysRequired := isSyntaxElementAlwaysRequired(syntaxStatementElement)
			currentElementExists := len(statement.Elements) > ei

			var currentElement *plain.DefinitionStatementElement

			if alwaysRequired && !currentElementExists {
				reportMismatch(i, syntaxStatement, "statement is shorter than syntax")

				match = false
				break
			}
			if currentElementExists {
				currentElement = &statement.Elements[ei]
			}

			if currentElement != nil {
				switch syntaxStatementElement.Kind {
				case macroAst.SyntaxStatementElementKindKeyword:
					if currentElement.Kind != plain.DefinitionStatementElementKindIdentifier {
						reportMismatch(i, syntaxStatement, fmt.Sprintf("expected keyword (%s), got %s (%v)", syntaxStatementElement.KeywordDef.Name, currentElement.Kind, currentElement))
						match = false
						break
					} else if currentElement.Identifier.Identifier != syntaxStatementElement.KeywordDef.Name {
						reportMismatch(i, syntaxStatement, fmt.Sprintf("expected keyword (%s), got %s", syntaxStatementElement.KeywordDef.Name, currentElement.Identifier.Identifier))
						match = false
						break
					}
				case macroAst.SyntaxStatementElementKindVariableKeyword:
					if currentElement.Kind != plain.DefinitionStatementElementKindIdentifier {
						reportMismatch(i, syntaxStatement, fmt.Sprintf("expected variable keyword (%s %s), got %s (%v)", syntaxStatementElement.VariableKeyword.Name, syntaxStatementElement.VariableKeyword.Type.ToDisplayName(), currentElement.Kind, currentElement))
						match = false
						break
					}
				case macroAst.SyntaxStatementElementKindAttributeList:
					if currentElement.Kind != plain.DefinitionStatementElementKindAttributeList {
						reportMismatch(i, syntaxStatement, fmt.Sprintf("expected attribute list, got %s", currentElement.Kind))
						match = false
						break
					}

					// check if the attribute list is valid
					err := isValidAttributeList(currentElement.AttributeList, syntaxStatementElement.AttributeList)

					if err != nil {
						reportMismatch(i, syntaxStatement, fmt.Sprintf(err.Error()))
						match = false
						break
					}
				case macroAst.SyntaxStatementElementKindArgumentList:
					if currentElement.Kind != plain.DefinitionStatementElementKindArgumentList {
						reportMismatch(i, syntaxStatement, fmt.Sprintf("expected argument list, got %s", currentElement.Kind))
						match = false
						break
					}

					// check if the argument list is valid
					err := isValidArgumentList(currentElement.ArgumentList, syntaxStatementElement.ArgumentList)

					if err != nil {
						reportMismatch(i, syntaxStatement, fmt.Sprintf(err.Error()))
						match = false
						break
					}
				default:
					reportMismatch(i, syntaxStatement, fmt.Sprintf("unexpected syntax element kind: %s", syntaxStatementElement.Kind))
				}

				syntaxElementMatch[ei] = i
			}

			ei++
		}

		if match {
			return &syntaxStatement, syntaxElementMatch, nil
		}
	}

	return nil, nil, fmt.Errorf("no matching syntax found, best match: %v, cause: %s", bestMatch, mismatchCause)

}

func isValidAttributeList(plainStatementElementAttributes *plain.DefinitionStatementElementAttributeList, syntaxStatementElementAttributes *macroAst.SyntaxStatementElementAttributeList) error {
	for _, plainAttribute := range plainStatementElementAttributes.Attributes {
		found := false

		for _, syntaxAttribute := range syntaxStatementElementAttributes.Attributes {
			if plainAttribute.Name == syntaxAttribute.Name {
				found = true

				validTyped := isValidValueType(plainAttribute.Value, syntaxAttribute.Type)

				if !validTyped {
					return fmt.Errorf("attribute %s expected type %s, got %s", plainAttribute.Name, syntaxAttribute.Type.ToDisplayName(), plainAttribute.Value.ToDisplayName())
				}

				break
			}
		}

		if !found {
			return fmt.Errorf("attribute %s is not allowed", plainAttribute.Name)
		}
	}

	return nil
}

func isValidArgumentList(plainStatementElementArguments *plain.DefinitionStatementElementArgumentList, syntaxStatementElementArguments *macroAst.SyntaxStatementElementArgumentList) error {
	// check if the argument list is valid
	if syntaxStatementElementArguments.VarArgs {
		return nil
	}

	return nil
}

func isValidValueType(value *common.Value, typeDefinition common.TypeDefinition) bool {
	return true // TODO: implement type validation for values
}

func isSyntaxElementAlwaysRequired(element macroAst.SyntaxStatementElement) bool {
	switch element.Kind {
	case macroAst.SyntaxStatementElementKindAttributeList:
		return false
	}

	return true
}
