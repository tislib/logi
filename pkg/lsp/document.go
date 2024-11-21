package lsp

import (
	"github.com/tislib/logi/pkg/lsp/common"
	"github.com/tislib/logi/pkg/lsp/protocol"
	"strings"
)

func (h *handler) onTextDocumentDidOpen(context *common.Context, params *protocol.DidOpenTextDocumentParams) error {
	h.fileContent[params.TextDocument.URI] = params.TextDocument.Text

	go h.analyseDocument(context, params.TextDocument.URI)

	return nil
}

func (h *handler) onTextDocumentDidClose(context *common.Context, params *protocol.DidCloseTextDocumentParams) error {
	delete(h.fileContent, params.TextDocument.URI)

	return nil
}

func (h *handler) onTextDocumentDidChange(context *common.Context, params *protocol.DidChangeTextDocumentParams) error {
	if len(params.ContentChanges) != 1 {
		return nil
	}

	var change = params.ContentChanges[0].(protocol.TextDocumentContentChangeEventWhole)

	h.fileContent[params.TextDocument.URI] = change.Text

	go h.analyseDocument(context, params.TextDocument.URI)

	return nil
}

func splitLines(existing string) []string {
	existing = strings.ReplaceAll(existing, "\r\n", "\n")
	existing = strings.ReplaceAll(existing, "\r", "\n")
	return strings.Split(existing, "\n")
}

func (h *handler) onTextDocumentHover(context *common.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	return nil, nil
}
