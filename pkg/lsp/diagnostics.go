package lsp

import (
	"errors"
	"github.com/tislib/logi/pkg/lsp/common"
	"github.com/tislib/logi/pkg/lsp/protocol"
	"github.com/tislib/logi/pkg/parser/logi"
	"github.com/tislib/logi/pkg/parser/macro"
)

func (h *handler) publishDiagnosticMacroError(context *common.Context, uri string, err error) {
	var mErr *macro.Error
	if errors.As(err, &mErr) {
		context.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
			URI: uri,
			Diagnostics: []protocol.Diagnostic{
				{
					Range: protocol.Range{
						Start: protocol.Position{
							Line:      protocol.UInteger(mErr.Line - 1),
							Character: protocol.UInteger(mErr.Column - 1),
						},
						End: protocol.Position{
							Line:      protocol.UInteger(mErr.Line - 1),
							Character: protocol.UInteger(mErr.Column - 1 + len(mErr.At)),
						},
					},
					Severity:           pointer(protocol.DiagnosticSeverityError),
					Message:            mErr.Msg,
					Tags:               nil,
					RelatedInformation: nil,
					Data:               nil,
				},
			},
		})
	}
}

func (h *handler) publishDiagnosticLogiError(context *common.Context, uri string, err error) {
	var mErr *logi.Error
	if errors.As(err, &mErr) {
		context.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
			URI: uri,
			Diagnostics: []protocol.Diagnostic{
				{
					Range: protocol.Range{
						Start: protocol.Position{
							Line:      protocol.UInteger(mErr.Line - 1),
							Character: protocol.UInteger(mErr.Column - 1),
						},
						End: protocol.Position{
							Line:      protocol.UInteger(mErr.Line - 1),
							Character: protocol.UInteger(mErr.Column - 1 + len(mErr.At)),
						},
					},
					Severity:           pointer(protocol.DiagnosticSeverityError),
					Message:            mErr.Msg,
					Tags:               nil,
					RelatedInformation: nil,
					Data:               nil,
				},
			},
		})
	}
}

func (h *handler) publishDiagnosticNoErrors(context *common.Context, uri string) {
	context.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: []protocol.Diagnostic{},
	})
}
