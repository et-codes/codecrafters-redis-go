package main

import (
	"fmt"
	"strconv"
)

const sep = "\r\n"

func decodeArrayLength(msg string) int {
	length, err := strconv.Atoi(msg[1:])
	if err != nil {
		logger.Error("Error parsing array length: %v", err)
	}
	return length
}

func encodeBulkStringArray(length int, bulkStrings ...string) string {
	encoded := fmt.Sprintf("*%d%s", length, sep)
	for _, str := range bulkStrings {
		encoded += fmt.Sprintf("$%d%s%s%s", len(str), sep, str, sep)
	}
	return encoded
}
