package server

import (
	"io"
)

// See: https://github.com/sourcegraph/go-langserver/blob/master/main.go#L179

func (self *Server) ServeStream(stream io.ReadWriteCloser) {
	<-self.newStreamConnection(stream).DisconnectNotify()
}
