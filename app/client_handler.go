package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	pingResponse = "+PONG\r\n"
	okResponse   = "+OK\r\n"
	nullResponse = "$-1\r\n"
)

type ClientHandler struct {
	Context context.Context
	Conn    io.ReadWriteCloser
	Server  *Server
	Store   map[string]string
}

func NewClientHandler(ctx context.Context, conn io.ReadWriteCloser, server *Server) *ClientHandler {
	return &ClientHandler{
		Context: ctx,
		Conn:    conn,
		Server:  server,
		Store:   make(map[string]string),
	}
}

type Command struct {
	Command string
	Args    []string
}

// Handle manages client communication with the server.
func (c *ClientHandler) Handle(wg *sync.WaitGroup) {
	defer c.Conn.Close()
	defer wg.Done()
	logger.Info("Connection initiated.")
	scanner := bufio.NewScanner(c.Conn)

	var (
		cmd         Command // Command object
		cmdLength   int     // Number of words in command + args
		cmdReceived int     // Number of words received from client
	)

	for {
		select {
		case <-c.Context.Done():
			logger.Info("Client handler stopping...")
			return
		default:
			for scanner.Scan() {
				msg := scanner.Text()

				switch msg[0] {
				case '*':
					length := decodeArrayLength(msg)
					cmdLength = length
					cmdReceived = 0
				case '$':
					// Don't need to do anything with this token for now.
				default:
					if cmdReceived < cmdLength {
						if cmdReceived == 0 {
							cmd.Command = msg
						} else {
							cmd.Args = append(cmd.Args, msg)
						}
						cmdReceived++
					}
				}

				// Execute command if the array is complete.
				if cmdReceived == cmdLength {
					if err := c.executeCommand(cmd); err != nil {
						logger.Error("Error executing command %v: %v", cmd, err)
					}
					// Reset for next command.
					cmd = Command{}
					cmdLength = 0
				}
			}
		}
	}
}

// executeCommand executes the command in the command array.
func (c *ClientHandler) executeCommand(cmd Command) error {
	switch cmd.Command {
	case "ping":
		return c.handlePing()
	case "echo":
		return c.handleEcho(cmd.Args)
	case "set":
		return c.handleSet(cmd.Args)
	case "get":
		return c.handleGet(cmd.Args)
	case "config":
		return c.handleConfig(cmd.Args)
	default:
		return fmt.Errorf("unrecognized command %q", cmd.Command)
	}
}

// handlePing handles PING commands.
func (c *ClientHandler) handlePing() error {
	logger.Info("PING command received.")
	return c.send(pingResponse)
}

// handleEcho handles ECHO commands.
func (c *ClientHandler) handleEcho(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("insufficient number of arguments for ECHO")
	}
	valToEcho := args[0]
	logger.Info("ECHO %q command received.", valToEcho)
	return c.send(encodeBulkString(valToEcho))
}

// handleSet handles SET commands.
func (c *ClientHandler) handleSet(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("insufficient number of arguments for SET")
	}
	key := args[0]
	value := args[1]
	logger.Info("SET %s: %q command received.", key, value)
	c.Store[key] = value

	// Check for expiration arguments.
	if len(args) == 4 && args[2] == "px" {
		expiry, err := strconv.Atoi(args[3])
		if err != nil {
			return fmt.Errorf("error parsing expiration time: %v", err)
		}
		go func() {
			time.Sleep(time.Duration(expiry) * time.Millisecond)
			delete(c.Store, key)
		}()
	}
	return c.send(okResponse)
}

// handleGet handles GET commands.
func (c *ClientHandler) handleGet(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("insufficient number of arguments for GET")
	}
	key := args[0]
	logger.Info("GET %s command received.", key)
	val, ok := c.Store[key]
	if ok {
		return c.send(encodeSimpleString(val))
	}
	return c.send(nullResponse)
}

// handleConfig handles CONFIG requests.
func (c *ClientHandler) handleConfig(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("insufficient number of arguments for CONFIG")
	}

	subCmd := strings.ToLower(args[0])

	switch subCmd {
	case "get":
		key := args[1]
		logger.Info("CONFIG GET %s command received.", key)

		val, found := c.Server.Config[key]
		if !found {
			return fmt.Errorf("config key %q not found", key)
		}

		return c.send(encodeBulkStringArray(2, key, val))

	case "set":
		if len(args) < 3 {
			return fmt.Errorf("insufficient number of arguments for CONFIG SET")
		}

		key := args[1]
		val := args[2]
		logger.Info("CONFIG SET %s: %q command received.", key, val)

		c.Server.Config[key] = val

		return c.send(okResponse)
	}

	return nil
}

// send sends the message to the client.
func (c *ClientHandler) send(msg string) error {
	_, err := c.Conn.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}
	return nil
}
