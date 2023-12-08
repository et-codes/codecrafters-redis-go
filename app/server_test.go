package main

import (
	"bytes"
	"sync"
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

func TestPingResponse(t *testing.T) {
	// Create client.
	tc := NewTestClient()
	c := NewClient(tc)

	// Send ping command to buffer.
	_, err := c.Conn.Write([]byte(pingCommand))
	assert.NoError(t, err)

	// Run the handler.
	var wg sync.WaitGroup
	wg.Add(1)
	c.Handle(&wg)
	defer wg.Wait()

	t.Run("responds to ping command", func(t *testing.T) {
		// Get the response.
		response := make([]byte, len(pingResponse))
		n, err := c.Conn.Read(response)
		assert.NoError(t, err)
		assert.Equal(t, pingResponse, string(response[:n]))
	})

	t.Run("responds to second ping command", func(t *testing.T) {
		// Send second ping.
		_, err = c.Conn.Write([]byte(pingCommand))
		assert.NoError(t, err)

		// Get the response.
		response := make([]byte, len(pingResponse))
		n, err := c.Conn.Read(response)
		assert.NoError(t, err)
		assert.Equal(t, pingResponse, string(response[:n]))
	})
}
