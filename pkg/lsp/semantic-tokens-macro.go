package lsp

import (
	macroAst "github.com/tislib/logi/pkg/ast/macro"
	"github.com/tislib/logi/pkg/lsp/common"
)

func (h *handler) analyserMacro(context *common.Context, uri string) []Token {
	ast, err := h.parser.ParseMacroContent(h.fileContent[uri])
	if err != nil {
		h.publishDiagnosticMacroError(context, uri, err)
	} else {
		h.publishDiagnosticNoErrors(context, uri)
	}

	var tl []Token

	for _, macro := range ast.Macros {
		for key, sourceLocation := range macro.SourceMap {
			var length int
			var tokenType string
			var tokenModifier string

			switch key {
			case "name":
				length = len(macro.Name)
				tokenType = "keyword"
				tokenModifier = "declaration"
			case "macro":
				length = 5
				tokenType = "decorator"
				tokenModifier = "declaration"
			case "kind":
				length = 5
				tokenType = "keyword"
				tokenModifier = "declaration"
			case "types":
				length = 5
				tokenType = "keyword"
				tokenModifier = "declaration"
			case "syntax":
				length = 6
				tokenType = "keyword"
				tokenModifier = "declaration"
			default:
				continue
			}

			if sourceLocation.Line == 0 {
				continue
			}

			tl = append(tl, Token{
				Line:      sourceLocation.Line - 1,
				StartChar: sourceLocation.Column - 1,
				Length:    length,
				TokenType: tokenIdMap[tokenType],
				Modifiers: tokenModifierIdMap[tokenModifier],
			})
		}

		for _, statement := range macro.Syntax.Statements {
			tl = append(tl, h.tokenizeMacroStatement(statement)...)
		}

		for _, statement := range macro.Types.Types {
			tl = append(tl, h.tokenizeMacroType(statement)...)
		}
	}

	return tl
}

func (h *handler) tokenizeMacroStatement(statement macroAst.SyntaxStatement) []Token {
	var tl []Token

	for _, element := range statement.Elements {
		items := h.tokenizeMacroStatementElement(element)
		tl = append(tl, items...)
	}

	return tl
}

func (h *handler) tokenizeMacroStatementElement(element macroAst.SyntaxStatementElement) []Token {
	// 	SyntaxStatementElementKindKeyword         SyntaxStatementElementKind = "Keyword"
	//	SyntaxStatementElementKindTypeReference   SyntaxStatementElementKind = "TypeReference"
	//	SyntaxStatementElementKindVariableKeyword SyntaxStatementElementKind = "VariableKeyword"
	//	SyntaxStatementElementKindCombination     SyntaxStatementElementKind = "Combination"
	//	SyntaxStatementElementKindStructure       SyntaxStatementElementKind = "Structure"
	//	SyntaxStatementElementKindParameterList   SyntaxStatementElementKind = "ParameterList"
	//	SyntaxStatementElementKindArgumentList    SyntaxStatementElementKind = "ArgumentList"
	//	SyntaxStatementElementKindCodeBlock       SyntaxStatementElementKind = "CodeBlock"
	//	SyntaxStatementElementKindExpressionBlock SyntaxStatementElementKind = "ExpressionBlock"
	//	SyntaxStatementElementKindAttributeList   SyntaxStatementElementKind = "AttributeList"
	switch element.Kind {

	}

	return nil
}

func (h *handler) tokenizeMacroType(statement macroAst.TypeStatement) []Token {
	return nil
}
