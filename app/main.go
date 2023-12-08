package main

import (
	"github.com/et-codes/codecrafters-redis-go/logging"
)

const (
	pingCommand  = "*1\r\n$4\r\nping\r\n" // 14 bytes
	pingResponse = "+PONG\r\n"            // 7 bytes
)

var logger = logging.New(logging.LevelDebug)

func main() {
	// Initiate server.
	s := NewServer("localhost", 6379)
	if err := s.Listen(); err != nil {
		logger.Fatal("Failed to bind to port %d", s.Port)
	}
	logger.Info("Listening on port %d...", s.Port)

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
}
