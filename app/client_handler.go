package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"sync"
)

const (
	pingCommand  = "*1\r\n$4\r\nping\r\n" // 14 bytes
	pingResponse = "+PONG\r\n"            // 7 bytes
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

	var (
		cmdArray         []string
		cmdArrayLength   int
		cmdArrayReceived int
	)

	for {
		select {
		case <-c.Context.Done():
			logger.Info("Client handler stopping...")
			return
		default:
			for scanner.Scan() {
				msg := scanner.Text()
				logger.Debug("Recieved: %s", msg)

				switch msg[0] {
				case '*':
					length := decodeArrayLength(msg)
					logger.Debug("Array length: %d", length)
					cmdArrayLength = length
					cmdArrayReceived = 0
				case '$':
					// Don't need to do anything with this token for now.
				default:
					if cmdArrayReceived < cmdArrayLength {
						cmdArray = append(cmdArray, msg)
						cmdArrayReceived++
					}
				}

				// Execute command if the array is complete.
				if cmdArrayReceived == cmdArrayLength {
					if err := c.executeCommand(cmdArray); err != nil {
						logger.Error("Error executing command %v: %v", cmdArray, err)
					}
				}
			}
		}
	}
}

// executeCommand executes the command in the command array.
func (c *ClientHandler) executeCommand(cmdArray []string) error {
	var err error
	switch cmdArray[0] {
	case "ping":
		logger.Info("PING command received.")
		err = c.sendMessage(pingResponse)
	case "echo":
		logger.Info("ECHO %q command received.", cmdArray[1])
		err = c.sendMessage(fmt.Sprintf("$%d\r\n%s\r\n", len(cmdArray[1]), cmdArray[1]))
	}
	return err
}

// sendMessage sends the passed message to the client.
func (c *ClientHandler) sendMessage(msg string) error {
	_, err := c.Conn.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}
	return nil
}
