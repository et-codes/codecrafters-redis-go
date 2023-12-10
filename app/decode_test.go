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
		"encodes array of bulk strings": {
			encodeBulkStringArray(2, "dir", "/tmp/redis-files"),
			"*2\r\n$3\r\ndir\r\n$16\r\n/tmp/redis-files\r\n",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.want, test.got)
		})
	}
}
