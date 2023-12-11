package main

import (
	"context"
	"os"

	"github.com/et-codes/codecrafters-redis-go/logging"
)

const (
	flagDBDir      = "--dir"        // command-line flag for db directory
	flagDBFilename = "--dbfilename" // command-line flag for db filename
	keyDBDir       = "dir"          // store key for db directory
	keyDBFilename  = "dbfilename"   // store key for db filename
	keyHost        = "host"         // store key for server host URL
	keyPort        = "port"         // store key for server port
)

var logger = logging.New(logging.LevelDebug)

func main() {
	ctx := context.Background()
	cfg := NewStore()

	if err := cfg.Add(keyHost, "localhost"); err != nil {
		logger.Fatal(err.Error())
	}

	if err := cfg.Add(keyPort, "6379"); err != nil {
		logger.Fatal(err.Error())
	}

	// Handle command-line arguments.
	args := os.Args
	for i := 1; i < len(args); i += 2 {
		if args[i] == flagDBDir && len(args) >= i {
			if err := cfg.Add(keyDBDir, args[i+1]); err != nil {
				logger.Fatal(err.Error())
			}
		} else if args[i] == flagDBFilename && len(args) >= i {
			if err := cfg.Add(keyDBFilename, args[i+1]); err != nil {
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
