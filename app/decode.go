package main

import (
	"strconv"
	"strings"
)

const sep = "\r\n"

func DecodeRESP(msg string) any {
	// Fix extra \ which gets added to escaped characters when passing
	// in the string as a parameter. TODO: why does this happen?
	msg = strings.ReplaceAll(msg, "\\r\\n", sep)

	switch msg[0] {
	case '+':
		return decodeSimpleString(msg)
	case ':':
		return decodeInteger(msg)
	case '$':
		return decodeBulkString(msg)
	case '*':
		return decodeArrayLength(msg)
	default:
		return nil
	}
}

// parseToken removes the \r\n separator from the incoming RESP message.
func parseToken(msg string) string {
	for i, c := range msg {
		if c == '\n' {
			return msg[:i]
		}
	}
	return ""
}

func decodeSimpleString(msg string) string {
	return strings.ToLower(msg[1:])
}

func decodeInteger(msg string) int {
	numstr := strings.TrimSuffix(msg[1:], sep)
	num, _ := strconv.Atoi(numstr)
	return num
}

func decodeBulkString(msg string) string {
	parts := strings.Split(msg[1:], sep)
	if len(parts) < 2 {
		return ""
	}
	return strings.ToLower(parts[1])
}

func decodeArrayLength(msg string) int {
	parts := strings.Split(msg[1:], sep)
	length, err := strconv.Atoi(parts[0])
	if err != nil {
		logger.Error("Error decoding array %s: %v", msg, err)
		return -1
	}
	return length
}
