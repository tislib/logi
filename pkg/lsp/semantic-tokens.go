package lsp

import (
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

func (h *handler) analyseDocument(context *common.Context, uri string) {
	h.fileSemanticTokens[uri] = h.analyser(context, uri)
}

func (h *handler) onTextDocumentSemanticTokensFull(context *common.Context, params *protocol.SemanticTokensParams) (*protocol.SemanticTokens, error) {
	var data = h.fileSemanticTokens[params.TextDocument.URI]
	return &protocol.SemanticTokens{
		Data: data,
	}, nil
}

func (h *handler) analyser(context *common.Context, uri string) []protocol.UInteger {
	var tl []Token
	if strings.HasSuffix(uri, ".lg") {
		tl = h.analyserLogi(context, uri)
	} else if strings.HasSuffix(uri, ".lgm") {
		tl = h.analyserMacro(context, uri)
	}

	return encodeTokens(tl)
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
