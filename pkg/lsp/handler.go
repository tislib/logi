package lsp

import (
	"github.com/tislib/logi/pkg/lsp/common"
	"github.com/tislib/logi/pkg/lsp/protocol"
	"github.com/tislib/logi/pkg/parser"
)

type Handler interface {
	Handle(context *common.Context) (r any, validMethod bool, validParams bool, err error)
}

type handler struct {
	protocol.Handler
	parser             parser.Parser
	fileContent        map[string]string
	fileSemanticTokens map[string][]protocol.UInteger
}

func NewHandler() Handler {
	var h = &handler{
		fileContent:        make(map[string]string),
		fileSemanticTokens: make(map[string][]protocol.UInteger),
	}

	h.Initialize = h.onInitialize
	h.Initialized = h.onInitialized
	h.Shutdown = h.onShutdown
	h.SetTrace = h.onSetTrace
	h.TextDocumentDidOpen = h.onTextDocumentDidOpen
	h.TextDocumentDidClose = h.onTextDocumentDidClose
	h.TextDocumentDidChange = h.onTextDocumentDidChange
	h.TextDocumentHover = h.onTextDocumentHover
	h.TextDocumentSemanticTokensFull = h.onTextDocumentSemanticTokensFull

	h.parser = parser.NewParser(true)

	return h
}
