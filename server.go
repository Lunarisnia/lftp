package lftp

import (
	"io"
	"net"

	"github.com/Lunarisnia/lftp/internal/dsu"
	"github.com/Lunarisnia/lftp/internal/lftparser"
)

type LFTPServer interface {
	Listen(address string) error
	Close() error
}

type RequestHandler func(header *dsu.LFTPHeader)

type Server struct {
	listener       net.Listener
	requestHandler RequestHandler
}

func NewLFTPServer(requestHandler RequestHandler) LFTPServer {
	return &Server{
		requestHandler: requestHandler,
	}
}

func (s *Server) Listen(address string) error {
	listener, err := net.Listen("tcp", address)
	s.listener = listener
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			return err
		}

		func(c net.Conn) {
			defer c.Close()
			rawContent, err := io.ReadAll(c)
			if err != nil {
				panic(err)
			}
			header, err := lftparser.ParseHeader(rawContent)
			if err != nil {
				panic(err)
			}

			s.requestHandler(header)

		}(connection)
	}
}

func (s *Server) Close() error {
	err := s.listener.Close()
	if err != nil {
		return err
	}
	return nil
}
