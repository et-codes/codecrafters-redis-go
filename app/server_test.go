package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestClient struct {
	*bytes.Buffer
}

func NewTestClient(buffer []byte) *TestClient {
	return &TestClient{bytes.NewBuffer(buffer)}
}

func (tc *TestClient) Close() error { return nil }

func TestPingResponse(t *testing.T) {
	// Create client.
	message := make([]byte, len(pingCommand))
	tc := NewTestClient(message)
	c := NewClient(tc)

	// Send ping command to buffer.
	_, err := c.Conn.Write([]byte(pingCommand))
	assert.NoError(t, err)

	// Run the handler.
	c.Handle()

	// Get the response.
	response := make([]byte, len(pingResponse))
	n, err := c.Conn.Read(response)
	assert.NoError(t, err)
	assert.Equal(t, pingResponse, string(response[:n]))
}
