package main

import (
	"context"
	"fmt"
	"net"
	"os/signal"
	"sync"
	"syscall"
)

type Server struct {
	Context  context.Context
	Host     string
	Port     int
	Listener net.Listener
}

func NewServer(ctx context.Context, host string, port int) *Server {
	return &Server{Context: ctx, Host: host, Port: port}
}

func (s *Server) Run() error {
	ctx, stop := signal.NotifyContext(s.Context, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	// Start listener.
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Host, s.Port))
	if err != nil {
		return fmt.Errorf("failed to bind to port %d: %v", s.Port, err)
	}
	s.Listener = listener
	logger.Info("Listening on port %d...", s.Port)

	// Stop listener when required.
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		logger.Debug("Closing listener...")
		listener.Close()
	}()

	// Listen for client connections and send to handler.
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			logger.Info("Shutting down...")
			stop()
			// Wait for workers to complete.
			wg.Wait()

			switch err.(type) {
			case *net.OpError:
				return nil
			default:
				return err
			}
		}

		client := NewClientHandler(ctx, conn)
		client.Handle()
		logger.Info("Client connected.")
	}
}
