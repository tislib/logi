package logi

import (
	"fmt"
	"github.com/tislib/logi/pkg/ast/common"
	macroAst "github.com/tislib/logi/pkg/ast/macro"
	"github.com/tislib/logi/pkg/ast/plain"
	"log"
)

func locateMacroDefinition(definition plain.Definition, ast macroAst.Ast) (*macroAst.Macro, error) {
	for _, macroDefinition := range ast.Macros {
		if macroDefinition.Name == definition.MacroName {
			return &macroDefinition, nil
		}
	}

	return nil, fmt.Errorf("macro definition not found")

}

func matchSyntaxStatement(statement plain.DefinitionStatement, syntaxStatement macroAst.SyntaxStatement, reportMismatch func(matchedUntil int, syntaxStatement macroAst.SyntaxStatement, reason string), macro *macroAst.Macro) []int {
	ei := 0
	syntaxElementMatch := make([]int, len(statement.Elements))

	var mtvi = 0

	for i, syntaxStatementElement := range syntaxStatement.Elements {
		// check if the statement is always required
		alwaysRequired := isSyntaxElementAlwaysRequired(syntaxStatementElement)
		currentElementExists := len(statement.Elements) > ei

		var currentElement *plain.DefinitionStatementElement

		if alwaysRequired && !currentElementExists {
			reportMismatch(i, syntaxStatement, "statement is shorter than syntax")

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
					break
				} else if currentElement.Identifier.Identifier != syntaxStatementElement.KeywordDef.Name {
					reportMismatch(i, syntaxStatement, fmt.Sprintf("expected keyword (%s), got %s", syntaxStatementElement.KeywordDef.Name, currentElement.Identifier.Identifier))
					break
				}
			case macroAst.SyntaxStatementElementKindVariableKeyword:
				switch currentElement.Kind {
				case plain.DefinitionStatementElementKindIdentifier:
					if currentElement.Identifier.Identifier != syntaxStatementElement.VariableKeyword.Name {
						reportMismatch(i, syntaxStatement, fmt.Sprintf("expected variable keyword (%s), got %s", syntaxStatementElement.VariableKeyword.Name, currentElement.Identifier.Identifier))
						break
					}
				case plain.DefinitionStatementElementKindValue:
					for _, typeStatement := range macro.Types.Types {
						if typeStatement.Name == syntaxStatementElement.VariableKeyword.Type.Name {
							// start matching the value
						}
					}
				}
				// todo type check
				log.Print("todo type check")
			case macroAst.SyntaxStatementElementKindAttributeList:
				if currentElement.Kind != plain.DefinitionStatementElementKindAttributeList {
					reportMismatch(i, syntaxStatement, fmt.Sprintf("expected attribute list, got %s", currentElement.Kind))

					break
				}

				// check if the attribute list is valid
				err := isValidAttributeList(currentElement.AttributeList, syntaxStatementElement.AttributeList)

				if err != nil {
					reportMismatch(i, syntaxStatement, fmt.Sprintf(err.Error()))

					break
				}
			case macroAst.SyntaxStatementElementKindArgumentList:
				if currentElement.Kind != plain.DefinitionStatementElementKindArgumentList {
					reportMismatch(i, syntaxStatement, fmt.Sprintf("expected argument list, got %s", currentElement.Kind))

					break
				}

				// check if the argument list is valid
				err := isValidArgumentList(currentElement.ArgumentList, syntaxStatementElement.ArgumentList)

				if err != nil {
					reportMismatch(i, syntaxStatement, fmt.Sprintf(err.Error()))

					break
				}
			case macroAst.SyntaxStatementElementKindCodeBlock:
				if currentElement.Kind != plain.DefinitionStatementElementKindCodeBlock {
					reportMismatch(i, syntaxStatement, fmt.Sprintf("expected code block, got %s", currentElement.Kind))

					break
				}
			default:
				reportMismatch(i, syntaxStatement, fmt.Sprintf("unexpected syntax element kind: %s", syntaxStatementElement.Kind))
			}
		}

		ei++
		if mtvi > 0 {
			ei += mtvi
			mtvi = 0
		}
	}
	return syntaxElementMatch
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
