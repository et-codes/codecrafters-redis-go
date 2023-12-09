package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
)

type ClientHandler struct {
	Context context.Context
	Conn    io.ReadWriteCloser
}

func NewClientHandler(ctx context.Context, conn io.ReadWriteCloser) *ClientHandler {
	return &ClientHandler{Conn: conn}
}

// Handle manages client communication with the server.
func (c *ClientHandler) Handle() {
	defer c.Conn.Close()
	logger.Info("Connection initiated.")

	scanner := bufio.NewScanner(c.Conn)

	for scanner.Scan() {
		msg := scanner.Text()

		switch strings.ToLower(msg) {
		case "ping":
			logger.Info("Ping command received.")
			err := c.sendMessage(pingResponse)
			if err != nil {
				logger.Error(err.Error())
			}
		case "quit":
			logger.Info("Quit command received, exiting client handler...")
			return
		default:
			logger.Debug("Received: %s", msg)
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
