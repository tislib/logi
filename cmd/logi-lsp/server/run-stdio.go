package server

import (
	"errors"
	"os"
)

func (self *Server) RunStdio() error {
	self.ServeStream(Stdio{})
	return nil
}

type Stdio struct{}

// ([io.Reader] interface)
func (Stdio) Read(p []byte) (int, error) {
	return os.Stdin.Read(p)
}

// ([io.Writer] interface)
func (Stdio) Write(p []byte) (int, error) {
	return os.Stdout.Write(p)
}

// ([io.Closer] interface)
func (Stdio) Close() error {
	return errors.Join(os.Stdin.Close(), os.Stdout.Close())
}
