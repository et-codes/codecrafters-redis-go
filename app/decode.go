package main

import (
	"strconv"
)

func decodeArrayLength(msg string) int {
	length, err := strconv.Atoi(msg[1:])
	if err != nil {
		logger.Error("Error parsing array length: %v", err)
	}
	return length
}
