package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/tislib/logi/cmd/lsp/server"
	"github.com/tislib/logi/pkg/lsp"
	"github.com/tislib/logi/pkg/lsp/common"
	"github.com/tislib/logi/pkg/lsp/protocol"
	"strings"
)

const lsName = "my language"

var (
	version string = "0.0.1"
	handler protocol.Handler
)

var activeDocument protocol.TextDocumentItem

func main() {
	hand := lsp.NewHandler()

	srv := server.NewServer(hand, lsName, false)

	srv.RunTCP(":7998")
}

func textDocumentDidChange(context *common.Context, params *protocol.DidChangeTextDocumentParams) error {
	//runDiagnostics(context, params.TextDocument.URI)
	return nil
}

func textDocumentDidOpen(context *common.Context, params *protocol.DidOpenTextDocumentParams) error {
	//runDiagnostics(context, params.TextDocument.URI)

	activeDocument = params.TextDocument
	return nil
}

func runDiagnostics(context *common.Context, uri protocol.DocumentUri) {
	var sev2 = protocol.DiagnosticSeverityInformation
	context.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
		URI: uri,
		Diagnostics: []protocol.Diagnostic{
			{
				Range: protocol.Range{
					Start: protocol.Position{
						Line:      3,
						Character: 0,
					},
					End: protocol.Position{
						Line:      3,
						Character: 13,
					},
				},
				Severity:           pointer(protocol.DiagnosticSeverityError),
				Code:               pointer(protocol.IntegerOrString{Value: "12332111"}),
				Source:             pointer("12332111 sad jasd sajkdjsa kdas"),
				Message:            "Unknown error not occoured, you can be happy",
				Tags:               nil,
				RelatedInformation: nil,
				Data:               nil,
			},
			{
				Range: protocol.Range{
					Start: protocol.Position{
						Line:      4,
						Character: 3,
					},
					End: protocol.Position{
						Line:      5,
						Character: 13,
					},
				},
				Severity:           &sev2,
				Code:               nil,
				CodeDescription:    nil,
				Source:             nil,
				Message:            "Unknown error not occoured, you can be happy",
				Tags:               nil,
				RelatedInformation: nil,
				Data:               nil,
			},
		},
	})
}

func pointer[T any](severityError T) *T {
	return &severityError
}

func initialize(context *common.Context, params *protocol.InitializeParams) (any, error) {
	capabilities := handler.CreateServerCapabilities()

	capabilities.SemanticTokensProvider = &protocol.SemanticTokensOptions{
		Legend: protocol.SemanticTokensLegend{
			TokenTypes: []string{
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
			},
			TokenModifiers: []string{
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
			},
		},
		Full: true,
	}

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    lsName,
			Version: &version,
		},
	}, nil
}

func initialized(context *common.Context, params *protocol.InitializedParams) error {
	return nil
}

func shutdown(context *common.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func setTrace(context *common.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}

// analyzeDocument analyzes the document's content and returns a slice of Token structs.
func analyzeDocument(source string) []Token {
	tokens := []Token{}
	lines := strings.Split(source, "\n")
	for lineNum, line := range lines {
		words := strings.Fields(line) // Split into words
		for _, word := range words {
			startChar := strings.Index(line, word)
			switch word {
			case "func", "var", "if", "else", "return":
				tokens = append(tokens, Token{Line: lineNum, StartChar: startChar, Length: len(word), TokenType: 1, Modifiers: 0}) // Keyword
			case "main", "analyzeDocument", "encodeTokens", "handleSemanticTokensFull", "handleInitialize":
				tokens = append(tokens, Token{Line: lineNum, StartChar: startChar, Length: len(word), TokenType: 2, Modifiers: 0}) // Function
			default:
				if strings.HasPrefix(word, "(") && strings.HasSuffix(word, ")") {
					tokens = append(tokens, Token{Line: lineNum, StartChar: startChar, Length: len(word), TokenType: 3, Modifiers: 0}) // Parameter
				} else {
					tokens = append(tokens, Token{Line: lineNum, StartChar: startChar, Length: len(word), TokenType: 5, Modifiers: 0}) // Variable
				}
			}
		}
	}
	return tokens
}

// handleSemanticTokensFull handles the `textDocument/semanticTokens/full` request.
func handleSemanticTokensFull(ctx *common.Context, params *protocol.SemanticTokensParams) (*protocol.SemanticTokens, error) {
	log.Printf("Received semanticTokens/full request: %#v", params)

	// Analyze the document and generate tokens
	tokens := analyzeDocument(activeDocument.Text)

	tokens = []Token{}

	for i := 0; i < 200; i++ {
		tokens = append(tokens, Token{Line: i, StartChar: 0, Length: 4, TokenType: i % 20, Modifiers: i / 20})
	}

	// Encode the token data
	data := encodeTokens(tokens)

	return &protocol.SemanticTokens{
		Data: data,
	}, nil
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

// Token represents a code element with its position, type, and modifiers.
type Token struct {
	Line      int
	StartChar int
	Length    int
	TokenType int // Corresponding to SemanticTokenTypes
	Modifiers int // Corresponding to SemanticTokenModifiers
}
