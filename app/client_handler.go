package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"reflect"
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
		remainingElements int
		cmdArray          []any
	)

	for {
		select {
		case <-c.Context.Done():
			logger.Info("Client handler stopping...")
			return
		default:
			for scanner.Scan() {
				msg := scanner.Text()
				if err := scanner.Err(); err != nil {
					logger.Error("Error scanning client: %v", err)
					return
				}

				decoded := DecodeRESP(msg)
				logger.Debug("Received message %s, decoded to %v", msg, decoded)
				if decoded == nil {
					continue
				}

				if remainingElements > 0 {
					cmdArray = append(cmdArray, decoded)
					remainingElements--
					logger.Debug("adding %v to array", decoded)
				}
				if remainingElements == 0 && len(cmdArray) > 0 {
					logger.Debug("Command array: %v", cmdArray)
					switch cmdArray[0] {
					case "ping":
						logger.Info("PING command received.")
						err := c.sendMessage(pingResponse)
						if err != nil {
							logger.Error(err.Error())
						}
					case "echo":
						logger.Info("ECHO %q command received.", cmdArray[1])
						err := c.sendMessage(fmt.Sprintf("$%d\r\n%s\r\n", len(cmdArray[1].(string)), cmdArray[1]))
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

				// If array is being built, intercept the elements.
				if reflect.TypeOf(decoded).Kind() == reflect.Int {
					// Start building array.
					logger.Debug("Building command array...")
					remainingElements = decoded.(int)
					cmdArray = []any{}
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
