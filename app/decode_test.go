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
		"decodes a simple string": {"+PONG", "pong"},
		"decodes an integer":      {":1000\r\n", 1000},
		"decodes a bulk string":   {"$5\r\nhello\r\n", "hello"},
		"decodes length of array": {"*3\r\n", 3},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := DecodeRESP(test.input)
			assert.Equal(t, test.want, got)
		})
	}
}
