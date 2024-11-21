package server

import (
	"log"
)

type JSONRPCLogger struct {
}

// ([jsonrpc2.Logger] interface)
func (self *JSONRPCLogger) Printf(format string, v ...any) {
	log.Printf(format, v...)
}
