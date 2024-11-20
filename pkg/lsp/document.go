package lsp

import (
	"fmt"
	"github.com/tislib/logi/pkg/lsp/common"
	"github.com/tislib/logi/pkg/lsp/protocol"
	"log"
	"strings"
)

func (h *handler) onTextDocumentDidOpen(context *common.Context, params *protocol.DidOpenTextDocumentParams) error {
	h.fileContent[params.TextDocument.URI] = params.TextDocument.Text
	return nil
}

func (h *handler) onTextDocumentDidClose(context *common.Context, params *protocol.DidCloseTextDocumentParams) error {
	delete(h.fileContent, params.TextDocument.URI)

	return nil
}

func (h *handler) onTextDocumentDidChange(context *common.Context, params *protocol.DidChangeTextDocumentParams) error {
	existing, ok := h.fileContent[params.TextDocument.URI]
	if !ok {
		return fmt.Errorf("document not found: %s", params.TextDocument.URI)
	}

	lines := splitLines(existing)

	for _, change := range params.ContentChanges {
		switch typed := change.(type) {
		case protocol.TextDocumentContentChangeEvent:
			startLine := typed.Range.Start.Line
			startChar := typed.Range.Start.Character
			endLine := typed.Range.End.Line
			endChar := typed.Range.End.Character

			if startLine == endLine {
				// Single line edit
				lines[startLine] = lines[startLine][:startChar] + typed.Text + lines[startLine][endChar:]
			} else {
				// Multi-line edit
				if int(startLine) > len(lines) || int(endLine) > len(lines) {
					return fmt.Errorf("invalid line range")
				}

				// Replace the first line with the beginning and the new text
				lines[startLine] = lines[startLine][:startChar] + typed.Text

				// Append the end of the last line to the modified first line
				lines[startLine] += lines[endLine][endChar:]

				// Remove the lines in between (including the last line)
				lines = append(lines[:startLine+1], lines[endLine+1:]...)
			}

		case protocol.TextDocumentContentChangeEventWhole:
			lines = splitLines(typed.Text)

		default:
			return fmt.Errorf("unknown content change type: %T", change)
		}

		log.Println("Change:", change)
	}

	h.fileContent[params.TextDocument.URI] = strings.Join(lines, "\n")
	log.Println("Updated content:", h.fileContent[params.TextDocument.URI])

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
