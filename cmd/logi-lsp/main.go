package main

import (
	"github.com/tislib/logi/cmd/logi-lsp/server"
	"github.com/tislib/logi/pkg/lsp"
)

func main() {
	hand := lsp.NewHandler()

	srv := server.NewServer(hand, false)

	//srv.RunTCP(":7998")
	srv.RunStdio()
}
