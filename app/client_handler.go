package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
)

type ClientHandler struct {
	Context context.Context
	Conn    io.ReadWriteCloser
}

func NewClientHandler(ctx context.Context, conn io.ReadWriteCloser) *ClientHandler {
	return &ClientHandler{Context: ctx, Conn: conn}
}

// Handle manages client communication with the server.
func (c *ClientHandler) Handle(wg *sync.WaitGroup) {
	defer c.Conn.Close()
	defer wg.Done()
	logger.Info("Connection initiated.")

	scanner := bufio.NewScanner(c.Conn)

	for {
		select {
		case <-c.Context.Done():
			logger.Info("Client handler stopping...")
			return
		default:
			for scanner.Scan() {
				msg := scanner.Text()

				switch strings.ToLower(msg) {
				case "ping":
					logger.Info("PING command received.")
					err := c.sendMessage(pingResponse)
					if err != nil {
						logger.Error(err.Error())
					}
				case "quit":
					logger.Info("QUIT command received, stopping client...")
					return
				default:
					logger.Debug("Received: %s", msg)
				}
			}
		}
	}
}

// sendMessage sends the passed message to the client.
func (c *ClientHandler) sendMessage(msg string) error {
	_, err := c.Conn.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}
	return nil
}
