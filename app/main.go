package main

import (
	"context"
	"os"

	"github.com/et-codes/codecrafters-redis-go/logging"
)

var logger = logging.New(logging.LevelDebug)

func main() {
	ctx := context.Background()

	cfg := ServerConfig{
		Host: "localhost",
		Port: 6379,
	}

	// Get command-line arguments.
	args := os.Args
	for i := 1; i < len(args); i += 2 {
		if args[i] == "--dir" && len(args) >= i {
			cfg.DBDir = args[i+1]
		} else if args[i] == "--dbfilename" && len(args) >= i {
			cfg.DBFilename = args[i+1]
		}
	}
	logger.Debug("Server config: %+v", cfg)

	// Initiate server.
	s := NewServer(ctx, cfg)
	if err := s.Run(); err != nil {
		logger.Fatal("Server error: ", err)
	}

	logger.Info("Server shutdown complete.")
}
