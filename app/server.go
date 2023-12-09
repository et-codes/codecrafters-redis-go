package main

import (
	"fmt"
	"net"
)

type Server struct {
	Host     string
	Port     int
	Listener net.Listener
}

func NewServer(host string, port int) *Server {
	return &Server{Host: host, Port: port}
}

func (s *Server) Run() error {
	// Start listener.
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Host, s.Port))
	if err != nil {
		return fmt.Errorf("failed to bind to port %d: %v", s.Port, err)
	}
	s.Listener = l

	// Listen for client connections and send to handler.
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			logger.Error("Error accepting connection: %v", err)
			break
		}

		client := NewClientHandler(conn)
		client.Handle()
	}

	return nil
}
