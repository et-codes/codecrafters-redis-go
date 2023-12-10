package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	tests := map[string]struct {
		got  any
		want any
	}{
		"decodes length of array": {decodeArrayLength("*3"), 3},
		"encodes simple string": {
			encodeSimpleString("hello"),
			"+hello\r\n",
		},
		"encodes bulk string": {
			encodeBulkString("hello"),
			"$5\r\nhello\r\n",
		},
		"encodes array of bulk strings": {
			encodeBulkStringArray(2, "dir", "/tmp/redis-files"),
			"*2\r\n$3\r\ndir\r\n$16\r\n/tmp/redis-files\r\n",
		},
		"encodes ping command": {
			encodeBulkStringArray(1, "ping"),
			"*1\r\n$4\r\nping\r\n",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.want, test.got)
		})
	}
}
