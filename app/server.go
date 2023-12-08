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

func (s *Server) Listen() error {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Host, s.Port))
	if err != nil {
		return fmt.Errorf("failed to bind to port %d: %v", s.Port, err)
	}
	s.Listener = l
	return nil
}