package main

import (
	"sync"

	"github.com/et-codes/codecrafters-redis-go/logging"
)

const (
	pingCommand  = "*1\r\n$4\r\nping\r\n"
	pingResponse = "+PONG\r\n"
)

var logger = logging.New(logging.LevelDebug)

func main() {
	var wg sync.WaitGroup

	// Initiate server.
	s := NewServer("localhost", 6379)
	if err := s.Listen(); err != nil {
		logger.Fatal("Failed to bind to port %d", s.Port)
	}
	logger.Info("Listening on port %d...", s.Port)

	// Listen for client connections and send to handler.
	clientCount := 0
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			logger.Error("Error accepting connection: %v", err)
			break
		}

		client := NewClient(conn)
		wg.Add(1)
		go client.Handle(&wg)
		clientCount++
	}

	// Wait for goroutines to finish.
	wg.Wait()
	logger.Info("All goroutines completed.")
}
