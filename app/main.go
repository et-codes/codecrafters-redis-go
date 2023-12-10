package main

import (
	"context"

	"github.com/et-codes/codecrafters-redis-go/logging"
)

var logger = logging.New(logging.LevelDebug)

func main() {
	ctx := context.Background()

	// Initiate server.
	c := ServerConfig{
		Host: "localhost",
		Port: 6379,
	}
	s := NewServer(ctx, c)
	if err := s.Run(); err != nil {
		logger.Fatal("Server error: ", err)
	}

	logger.Info("Server shutdown complete.")
}
