package main

import (
	"context"

	"github.com/et-codes/codecrafters-redis-go/logging"
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
