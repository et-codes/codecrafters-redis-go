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
		"decodes a simple string":       {"+PONG", "pong"},
		"decodes a bulk string":         {"$5\r\nhello\r\n", "hello"},
		"decodes an empty array":        {"*0\r\n", []any{}},
		"decodes array of bulk strings": {"*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n", []any{"hello", "world"}},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := DecodeRESP(test.input)
			assert.Equal(t, test.want, got)
		})
	}
}
