package main

import (
	"context"
	"github.com/TobiasYin/go-lsp/logs"
	"github.com/TobiasYin/go-lsp/lsp"
	"github.com/TobiasYin/go-lsp/lsp/defines"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var lspCmd = &cobra.Command{
	Use:   "lsp",
	Short: "lsp - Language Server Protocol for Logi",
	Long:  `lsp - Language Server Protocol for Logi, listens on :7052 and provides LSP for Logi files`,
	RunE: func(cmd *cobra.Command, args []string) error {
		//<-jsonrpc2.NewConn(
		//	cmd.Context(),
		//	jsonrpc2.NewBufferedStream(stdrwc{}, jsonrpc2.VSCodeObjectCodec{}),
		//	lsp.NewHandler(),
		//).DisconnectNotify()

		logs.Init(log.Default())

		server := lsp.NewServer(&lsp.Options{
			CompletionProvider: &defines.CompletionOptions{
				TriggerCharacters: &[]string{"."},
			},
			Network: "tcp",
		})
		server.OnHover(func(ctx context.Context, req *defines.HoverParams) (result *defines.Hover, err error) {
			logs.Println(req)
			return &defines.Hover{Contents: defines.MarkupContent{Kind: defines.MarkupKindPlainText, Value: "hello world"}}, nil
		})

		var hello = "Hello"

		server.OnCompletion(func(ctx context.Context, req *defines.CompletionParams) (result *[]defines.CompletionItem, err error) {
			logs.Println(req)
			d := defines.CompletionItemKindText
			return &[]defines.CompletionItem{defines.CompletionItem{
				Label:      "code",
				Kind:       &d,
				InsertText: &hello,
			}}, nil
		})

		//server.OnDidChangeTextDocument(func(ctx context.Context, req *defines.DidChangeTextDocumentParams) (err error) {
		//
		//	log.Println(req.ContentChanges)
		//
		//	return errors.New("some unhjappy result")
		//})

		server.OnDocumentColor(func(ctx context.Context, req *defines.DocumentColorParams) (result *[]defines.ColorInformation, err error) {
			if result == nil {
				result = new([]defines.ColorInformation)
			}
			*result = append(*result, defines.ColorInformation{
				Range: defines.Range{
					Start: defines.Position{
						Line:      0,
						Character: 0,
					},
					End: defines.Position{
						Line:      0,
						Character: 10,
					},
				},
				Color: defines.Color{
					Red:   20,
					Green: 200,
					Blue:  20,
					Alpha: 90,
				},
			})

			return result, err
		})

		server.Run()

		return nil
	},
}

type stdrwc struct{}

func (stdrwc) Read(p []byte) (int, error) {
	return os.Stdin.Read(p)
}

func (c stdrwc) Write(p []byte) (int, error) {
	return os.Stdout.Write(p)
}

func (c stdrwc) Close() error {
	if err := os.Stdin.Close(); err != nil {
		return err
	}
	return os.Stdout.Close()
}

func init() {
	rootCmd.AddCommand(lspCmd)
}
