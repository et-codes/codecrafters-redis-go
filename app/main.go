package main

import (
	"context"

	"github.com/et-codes/codecrafters-redis-go/logging"
)

const (
	pingCommand  = "*1\r\n$4\r\nping\r\n" // 14 bytes
	pingResponse = "+PONG\r\n"            // 7 bytes
)

var logger = logging.New(logging.LevelDebug)

func main() {
	ctx := context.Background()

	// Initiate server.
	s := NewServer(ctx, "localhost", 6379)
	if err := s.Run(); err != nil {
		logger.Fatal("Server error: ", err)
	}

	logger.Info("Server shutdown complete.")
}
