package lsp

import (
	"fmt"
	common2 "github.com/tislib/logi/pkg/ast/common"
	"github.com/tislib/logi/pkg/ast/plain"
	"github.com/tislib/logi/pkg/lsp/common"
)

func (h *handler) analyserLogi(context *common.Context, uri string) []Token {
	ast, err := h.parser.ParseLogiPlainContent(h.fileContent[uri])
	if err != nil {
		h.publishDiagnosticLogiError(context, uri, err)
	} else {
		h.publishDiagnosticNoErrors(context, uri)
	}

	var tl []Token

	for _, definition := range ast.Definitions {
		tl = append(tl, Token{
			Line:      definition.MacroNameSourceLocation.Line - 1,
			StartChar: definition.MacroNameSourceLocation.Column - 1,
			Length:    len(definition.MacroName),
			TokenType: tokenIdMap["decorator"],
			Modifiers: tokenModifierIdMap["declaration"],
		})
		tl = append(tl, Token{
			Line:      definition.NameSourceLocation.Line - 1,
			StartChar: definition.NameSourceLocation.Column - 1,
			Length:    len(definition.Name),
			TokenType: tokenIdMap["class"],
			Modifiers: tokenModifierIdMap["declaration"],
		})

		for _, statement := range definition.Statements {
			tl = append(tl, h.tokenizeStatement(statement)...)
		}
	}

	return tl
}

func (h *handler) tokenizeStatement(statement plain.DefinitionStatement) []Token {
	var tl []Token

	for i, element := range statement.Elements {
		items := h.tokenizeStatementElements(element, i)
		tl = append(tl, items...)
	}
	return tl
}

func (h *handler) tokenizeStatementElements(element plain.DefinitionStatementElement, i int) []Token {
	var tl []Token

	var tokenTypeStr string
	var tokenModifierStr string
	var length int

	switch element.Kind {
	case plain.DefinitionStatementElementKindIdentifier:
		if i == 0 {
			tokenTypeStr = "property"
		} else {
			tokenTypeStr = "keyword"
		}
		tokenModifierStr = "declaration"
		length = len(element.Identifier.Identifier)
	case plain.DefinitionStatementElementKindValue:
		switch element.Value.Value.Kind {
		case common2.ValueKindString:
			tokenTypeStr = "string"
			tokenModifierStr = "declaration"
			length = len(*element.Value.Value.String) + 2
		case common2.ValueKindInteger:
			tokenTypeStr = "number"
			tokenModifierStr = "declaration"
			length = len(fmt.Sprintf("%d", *element.Value.Value.Integer))
		case common2.ValueKindFloat:
			tokenTypeStr = "number"
			tokenModifierStr = "declaration"
			length = len(fmt.Sprintf("%f", *element.Value.Value.Float))
		case common2.ValueKindBoolean:
			tokenTypeStr = "keyword"
			tokenModifierStr = "declaration"
			length = len(fmt.Sprintf("%t", *element.Value.Value.Boolean))
		case common2.ValueKindMap:

		default:
			panic("unexpected value kind: " + element.Value.Value.Kind)
		}
	case plain.DefinitionStatementElementKindArray:
		for _, item := range element.Array.Items {
			tl = append(tl, h.tokenizeStatement(item)...)
		}
		return tl
	case plain.DefinitionStatementElementKindStruct:
		for _, item := range element.Struct.Statements {
			tl = append(tl, h.tokenizeStatement(item)...)
		}
		return tl
	case plain.DefinitionStatementElementKindAttributeList:
	case plain.DefinitionStatementElementKindArgumentList:
	case plain.DefinitionStatementElementKindParameterList:
	}

	tl = append(tl, Token{
		Line:      element.SourceLocation.Line - 1,
		StartChar: element.SourceLocation.Column - 1,
		Length:    length,
		TokenType: tokenIdMap[tokenTypeStr],
		Modifiers: tokenModifierIdMap[tokenModifierStr],
	})
	return tl
}
