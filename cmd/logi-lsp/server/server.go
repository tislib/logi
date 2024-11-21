package server

import (
	"github.com/tislib/logi/pkg/lsp/common"
	"time"
)

var DefaultTimeout = time.Minute

//
// Server
//

type Server struct {
	Handler common.Handler
	Debug   bool

	Timeout          time.Duration
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	StreamTimeout    time.Duration
	WebSocketTimeout time.Duration
}

func NewServer(handler common.Handler, debug bool) *Server {
	return &Server{
		Handler:          handler,
		Debug:            debug,
		Timeout:          DefaultTimeout,
		ReadTimeout:      DefaultTimeout,
		WriteTimeout:     DefaultTimeout,
		StreamTimeout:    DefaultTimeout,
		WebSocketTimeout: DefaultTimeout,
	}
}
