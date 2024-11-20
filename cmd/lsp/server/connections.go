package server

import (
	contextpkg "context"
	"io"

	"github.com/gorilla/websocket"
	"github.com/sourcegraph/jsonrpc2"
	wsjsonrpc2 "github.com/sourcegraph/jsonrpc2/websocket"
)

func (self *Server) newStreamConnection(stream io.ReadWriteCloser) *jsonrpc2.Conn {
	handler := self.newHandler()
	connectionOptions := self.newConnectionOptions()

	context, cancel := contextpkg.WithTimeout(contextpkg.Background(), self.StreamTimeout)
	defer cancel()

	return jsonrpc2.NewConn(context, jsonrpc2.NewBufferedStream(stream, jsonrpc2.VSCodeObjectCodec{}), handler, connectionOptions...)
}

func (self *Server) newWebSocketConnection(socket *websocket.Conn) *jsonrpc2.Conn {
	handler := self.newHandler()
	connectionOptions := self.newConnectionOptions()

	context, cancel := contextpkg.WithTimeout(contextpkg.Background(), self.WebSocketTimeout)
	defer cancel()

	return jsonrpc2.NewConn(context, wsjsonrpc2.NewObjectStream(socket), handler, connectionOptions...)
}

func (self *Server) newConnectionOptions() []jsonrpc2.ConnOpt {
	if self.Debug {
		return []jsonrpc2.ConnOpt{jsonrpc2.LogMessages(&JSONRPCLogger{})}
	} else {
		return nil
	}
}
