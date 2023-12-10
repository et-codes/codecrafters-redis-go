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
	Context context.Context
	Config  ServerConfig
}

type ServerConfig struct {
	Host       string
	Port       int
	DBDir      string
	DBFilename string
}

func NewServer(ctx context.Context, config ServerConfig) *Server {
	return &Server{Context: ctx, Config: config}
}

func (s *Server) Run() error {
	// Listen on Context for SIGINT or SIGTERM to shutdown gracefully.
	ctx, stop := signal.NotifyContext(s.Context, syscall.SIGINT, syscall.SIGTERM)
	defer stop() // stop will trigger ctx.Done() signal

	var wg sync.WaitGroup

	// Start TCP listener.
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port))
	if err != nil {
		return fmt.Errorf("failed to bind to port %d: %v", s.Config.Port, err)
	}
	logger.Info("Listening on port %d...", s.Config.Port)

	// Start goroutine that stops listener when signal is received.
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		logger.Debug("Closing listener...")
		listener.Close()
	}()

	// Listen for client connections and send to handler.
	for {
		conn, err := listener.Accept()
		if err != nil {
			// Accept will return error when the listener is closed.  Use
			// this to implement graceful shutdown.
			logger.Info("Shutting down server...")

			// Send ctx.Done() to handlers. TODO: this doesn't work!
			stop()
			logger.Debug("Stop message sent, waiting for goroutines to complete...")

			// Wait for goroutines to complete.
			wg.Wait()
			logger.Debug("Waitgroup clear.")

			switch err.(type) {
			case *net.OpError:
				// net.OpError is received when listener is closed. No
				// need to return this error, it is expected on shutdown.
				return nil
			default:
				// Any other error type gets sent back to caller.
				return err
			}
		}

		// Create ClientHandler and start goroutine.
		client := NewClientHandler(ctx, conn)
		wg.Add(1)
		go client.Handle(&wg)

		logger.Info("Client connected.")
	}
}
