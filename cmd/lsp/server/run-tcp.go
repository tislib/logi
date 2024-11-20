package server

import (
	log "github.com/sirupsen/logrus"
)

func (self *Server) RunTCP(address string) error {
	listener, err := self.newNetworkListener("tcp", address)
	if err != nil {
		return err
	}

	log.Info("listening for TCP connections")

	var connectionCount uint64

	for {
		connection, err := (*listener).Accept()
		if err != nil {
			return err
		}

		connectionCount++

		go self.ServeStream(connection)
	}
}
