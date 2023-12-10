package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	tests := map[string]struct {
		input string
		want  any
	}{
		"decodes a simple string":       {"+PONG", Decoded{"pong"}},
		"decodes an integer":            {":1000\r\n", Decoded{1000}},
		"decodes a bulk string":         {"$5\r\nhello\r\n", Decoded{"hello"}},
		"decodes an empty array":        {"*0\r\n", Decoded{}},
		"decodes array of bulk strings": {"*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n", Decoded{"hello", "world"}},
		"decodes ECHO message":          {"*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n", Decoded{"echo", "hey"}},
		"decoded PING command":          {pingCommand, Decoded{"ping"}},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := DecodeRESP(test.input)
			assert.Equal(t, test.want, got)
		})
	}
}
