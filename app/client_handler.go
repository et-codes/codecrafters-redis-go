package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"sync"
)

const (
	pingCommand  = "*1\r\n$4\r\nping\r\n"
	pingResponse = "+PONG\r\n"
	echoResponse = "$%d\r\n%s\r\n" // follow with length and value
	okResponse   = "+OK\r\n"
	getResponse  = "+%s\r\n" // follow with value
)

type ClientHandler struct {
	Context context.Context
	Conn    io.ReadWriteCloser
	Store   map[string]string
}

func NewClientHandler(ctx context.Context, conn io.ReadWriteCloser) *ClientHandler {
	store := make(map[string]string)
	return &ClientHandler{Context: ctx, Conn: conn, Store: store}
}

type Command struct {
	Command  string
	Args     []string
	Response string
}

func NewCommmand() *Command {
	return &Command{}
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
					// Reset for next command.
					cmdArray = []string{}
					cmdArrayLength = 0
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
		if len(cmdArray) < 2 {
			return fmt.Errorf("insufficient number of arguments for ECHO")
		}
		length := len(cmdArray[1])
		valToEcho := cmdArray[1]
		logger.Info("ECHO %q command received.", valToEcho)
		err = c.sendMessage(fmt.Sprintf(echoResponse, length, valToEcho))
	case "set":
		if len(cmdArray) < 3 {
			return fmt.Errorf("insufficient number of arguments for SET")
		}
		key := cmdArray[1]
		value := cmdArray[2]
		logger.Info("SET %s: %q command received.", key, value)
		c.Store[key] = value
		err = c.sendMessage(okResponse)
	case "get":
		if len(cmdArray) < 2 {
			return fmt.Errorf("insufficient number of arguments for GET")
		}
		key := cmdArray[1]
		val, ok := c.Store[key]
		if !ok {
			return fmt.Errorf("key %s not found", key)
		}
		err = c.sendMessage(fmt.Sprintf(getResponse, val))
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
