package main

import (
	"regexp"
	"strconv"
	"strings"
)

const sep = "\r\n"

func DecodeRESP(msg string) any {
	switch msg[0] {
	case '+':
		return decodeSimpleString(msg)
	case '$':
		return decodeBulkString(msg)
	case '*':
		return decodeArray(msg)
	default:
		return ""
	}
}

func decodeSimpleString(msg string) string {
	return strings.ToLower(msg[1:])
}

func decodeBulkString(msg string) string {
	parts := strings.Split(msg[1:], sep)
	return strings.ToLower(parts[1])
}

func decodeArray(msg string) []any {
	result := []any{}

	lengthStr, remainder, _ := strings.Cut(msg[1:], sep)
	length, _ := strconv.Atoi(lengthStr)
	logger.Debug("Array with length %d: %s", length, remainder)

	if length == 0 {
		return result
	}

	re := regexp.MustCompile(".*\r\n.*\r\n")
	matches := re.FindAllString(remainder, -1)
	for _, part := range matches {
		result = append(result, DecodeRESP(part))
	}
	return result
}
