package main

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"sync"
)

type ClientHandler struct {
	Conn   io.ReadWriteCloser
	Reader *bufio.Reader
	Writer *bufio.Writer
}

func NewClientHandler(conn io.ReadWriteCloser) *ClientHandler {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	return &ClientHandler{Conn: conn, Reader: reader, Writer: writer}
}

// Handle manages client communication with the server.
func (c *ClientHandler) Handle(wg *sync.WaitGroup) {
	logger.Info("Connection initiated.")
	defer c.Conn.Close()
	defer wg.Done()

	msg, err := c.waitForMessage()
	if err != nil {
		logger.Error(err.Error())
	}

	if reflect.DeepEqual(msg, []byte(pingCommand)) {
		logger.Info("Ping received.")
		err := c.sendMessage(pingResponse)
		if err != nil {
			logger.Error(err.Error())
		}
	}

	logger.Info("Ending client handler.")
}

// waitForMessage waits for a client message and returns it to the caller.
func (c *ClientHandler) waitForMessage() ([]byte, error) {
	buffer := make([]byte, 14)
	n, err := io.ReadFull(c.Conn, buffer)
	if err != nil {
		if err != io.EOF {
			return nil, fmt.Errorf("error reading connection: %v", err)
		}
	}
	logger.Debug("Client message received (%d bytes): %s", n, literalMessage(buffer[:n]))

	return buffer[:n], nil
}

// sendMessage sends the passed message to the client.
func (c *ClientHandler) sendMessage(msg string) error {
	_, err := c.Conn.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}
	return nil
}

// literalMessage converts a byte array to a string with esacped literals.
// For example []byte{13, 10} would return "\r\n".
func literalMessage(msg []byte) string {
	// Convert to literal string.
	literal := strconv.Quote(string(msg))

	// Remove beginning and ending double quotes.
	return literal[1 : len(literal)-1]
}
