package lftp

import (
	"io"
	"net"

	"github.com/Lunarisnia/lftp/internal/lftparser"
)

type Server struct{}

func (s *Server) Listen(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			return err
		}

		s.handleConnection(connection)
	}
}

func (s *Server) handleConnection(c net.Conn) {
	defer c.Close()
	rawContent, err := io.ReadAll(c)
	if err != nil {
		panic(err)
	}
	lftparser.ParseHeader(rawContent)
}
