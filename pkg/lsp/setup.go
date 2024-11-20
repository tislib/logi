package lsp

import (
	"github.com/tislib/logi/pkg/lsp/common"
	"github.com/tislib/logi/pkg/lsp/protocol"
)

func (h *handler) onInitialize(context *common.Context, params *protocol.InitializeParams) (any, error) {
	capabilities := h.CreateServerCapabilities()

	capabilities.SemanticTokensProvider = h.prepareSemanticTokensOptions()

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    "Logi",
			Version: pointer("1.0.0"),
		},
	}, nil
}

func (h *handler) onInitialized(context *common.Context, params *protocol.InitializedParams) error {
	return nil
}

func (h *handler) onShutdown(context *common.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)

	return nil
}

func (h *handler) onSetTrace(context *common.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)

	return nil
}
