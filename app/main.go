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
	if err := s.Run(); err != nil {
		logger.Fatal("Server error: ", err)
	}
	logger.Info("Listening on port %d...", s.Port)
}
