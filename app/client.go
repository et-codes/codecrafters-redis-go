package main

import (
	"io"
	"strings"
	"sync"
)

type Client struct {
	Conn io.ReadWriteCloser
}

func NewClient(conn io.ReadWriteCloser) *Client {
	return &Client{Conn: conn}
}

// Handle manages client communication with the server.
func (c *Client) Handle(wg *sync.WaitGroup) {
	logger.Info("Connection initiated.")
	defer c.Conn.Close()
	defer wg.Done()

	buffer := make([]byte, 1024)
	n, err := c.Conn.Read(buffer)
	if err != nil {
		logger.Error("Error reading connection: %v", err)
		return
	}
	logger.Debug("Received %d bytes: %s", n, string(buffer[:n]))

	if strings.Contains(strings.ToLower(string(buffer)), "ping") {
		logger.Info("Ping received.")
		_, err := c.Conn.Write([]byte(pingResponse))
		if err != nil {
			logger.Error("Error reading connection: %v", err)
			return
		}
	}

	logger.Debug("Ending client handler.")
}
