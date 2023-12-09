package main

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestClient implements io.ReadWriteCloser and can be used in place of
// a net.Conn connection.
type TestClient struct {
	*bytes.Buffer
}

func NewTestClient() *TestClient {
	return &TestClient{bytes.NewBuffer([]byte{})}
}

// Close doesn't do anything except make TestClient a Closer.
func (tc *TestClient) Close() error { return nil }

func TestClientHandler(t *testing.T) {
	// Create client.
	tc := NewTestClient()
	c := NewClientHandler(tc)

	t.Run("send message", func(t *testing.T) {
		// Send message.
		err := c.sendMessage(pingResponse)
		assert.NoError(t, err)

		// Read message back.
		scanner := bufio.NewScanner(c.Conn)
		var msg string
		for scanner.Scan() {
			msg = scanner.Text()
		}
		assert.Equal(t, msg, "+PONG")
	})
}
