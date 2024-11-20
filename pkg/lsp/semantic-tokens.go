package lsp

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	common2 "github.com/tislib/logi/pkg/ast/common"
	"github.com/tislib/logi/pkg/ast/plain"
	"github.com/tislib/logi/pkg/lsp/common"
	"github.com/tislib/logi/pkg/lsp/protocol"
	"github.com/tislib/logi/pkg/parser/lexer"
	"github.com/tislib/logi/pkg/parser/logi"
	"strings"
)

var tokenConfigs = logi.Tokens()
var tokenConfigsIdMap = make(map[int]lexer.TokenConfig)

func init() {
	for _, tokenConfig := range tokenConfigs {
		tokenConfigsIdMap[tokenConfig.Id] = tokenConfig
	}
}

// Token represents a code element with its position, type, and modifiers.
type Token struct {
	Line      int
	StartChar int
	Length    int
	TokenType int // Corresponding to SemanticTokenTypes
	Modifiers int // Corresponding to SemanticTokenModifiers
}

var tokens = []string{
	"namespace",
	"type",
	"class",
	"enum",
	"interface",
	"struct",
	"typeParameter",
	"parameter",
	"variable",
	"property",
	"enumMember",
	"event",
	"function",
	"method",
	"macro",
	"keyword",
	"modifier",
	"comment",
	"string",
	"number",
	"regexp",
	"operator",
	"decorator",
}

var tokenModifiers = []string{
	"declaration",
	"definition",
	"readonly",
	"static",
	"deprecated",
	"abstract",
	"async",
	"modification",
	"documentation",
	"defaultLibrary",
}

var tokenIdMap = make(map[string]int)
var tokenModifierIdMap = make(map[string]int)

func init() {
	for i, token := range tokens {
		tokenIdMap[token] = i
	}

	for i, tokenModifier := range tokenModifiers {
		tokenModifierIdMap[tokenModifier] = i
	}
}

func (h *handler) prepareSemanticTokensOptions() *protocol.SemanticTokensOptions {
	return &protocol.SemanticTokensOptions{
		Legend: protocol.SemanticTokensLegend{
			TokenTypes:     tokens,
			TokenModifiers: tokenModifiers,
		},
		Full: true,
	}
}

func (h *handler) onTextDocumentSemanticTokensFull(context *common.Context, params *protocol.SemanticTokensParams) (*protocol.SemanticTokens, error) {
	tl, err := h.analyser(params)
	if err != nil {
		return nil, err
	}

	return &protocol.SemanticTokens{
		Data: encodeTokens(tl),
	}, nil
}

func (h *handler) analyser(params *protocol.SemanticTokensParams) ([]Token, error) {
	if strings.HasSuffix(params.TextDocument.URI, ".lg") {
		return h.analyserLogi(params)
	}

	return nil, nil
}

func (h *handler) analyserLogi(params *protocol.SemanticTokensParams) ([]Token, error) {
	ast, err := h.parser.ParseLogiPlainContent(h.fileContent[params.TextDocument.URI])
	if err != nil {
		log.Error(err)
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

	return tl, nil
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
	case plain.DefinitionStatementElementKindCodeBlock:
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

// encodeTokens encodes the tokens into the LSP format.
func encodeTokens(tokens []Token) []protocol.UInteger {
	var data []protocol.UInteger
	prevLine := 0
	prevChar := 0
	for _, token := range tokens {
		deltaLine := token.Line - prevLine
		deltaStart := token.StartChar
		if deltaLine == 0 {
			deltaStart = token.StartChar - prevChar
		}
		data = append(data, protocol.UInteger(deltaLine), protocol.UInteger(deltaStart), protocol.UInteger(token.Length), protocol.UInteger(token.TokenType), protocol.UInteger(token.Modifiers))
		prevLine = token.Line
		prevChar = token.StartChar
	}
	return data
}
