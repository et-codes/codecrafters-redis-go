package main

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestClient struct {
	*bytes.Buffer
}

func NewTestClient() *TestClient {
	return &TestClient{bytes.NewBuffer([]byte{})}
}

func (tc *TestClient) Close() error { return nil }

func TestSendMessage(t *testing.T) {
	// Create client.
	tc := NewTestClient()
	c := NewClientHandler(tc)

	// Send message.
	c.sendMessage(pingResponse)

	// Read message back.
	scanner := bufio.NewScanner(c.Conn)
	var msg string
	for scanner.Scan() {
		msg = scanner.Text()
	}
	assert.Equal(t, msg, "+PONG")
}
