package main

import (
	"context"
	"os"

	"github.com/et-codes/codecrafters-redis-go/logging"
)

var logger = logging.New(logging.LevelDebug)

func main() {
	ctx := context.Background()
	cfg := NewStore()

	if err := cfg.Add("host", "localhost"); err != nil {
		logger.Fatal(err.Error())
	}

	if err := cfg.Add("port", "6379"); err != nil {
		logger.Fatal(err.Error())
	}

	// Handle command-line arguments.
	args := os.Args
	for i := 1; i < len(args); i += 2 {
		if args[i] == "--dir" && len(args) >= i {
			if err := cfg.Add("dir", args[i+1]); err != nil {
				logger.Fatal(err.Error())
			}
		} else if args[i] == "--dbfilename" && len(args) >= i {
			if err := cfg.Add("dbfilename", args[i+1]); err != nil {
				logger.Fatal(err.Error())
			}
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
