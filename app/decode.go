package main

import (
	"regexp"
	"strconv"
	"strings"
)

const sep = "\r\n"

type Decoded []any

func DecodeRESP(msg string) Decoded {
	switch msg[0] {
	case '+':
		return decodeSimpleString(msg)
	case ':':
		return decodeInteger(msg)
	case '$':
		return decodeBulkString(msg)
	case '*':
		return decodeArray(msg)
	default:
		return Decoded{}
	}
}

func decodeSimpleString(msg string) Decoded {
	return Decoded{strings.ToLower(msg[1:])}
}

func decodeInteger(msg string) Decoded {
	numstr := strings.TrimSuffix(msg[1:], sep)
	num, _ := strconv.Atoi(numstr)
	return Decoded{num}
}

func decodeBulkString(msg string) Decoded {
	parts := strings.Split(msg[1:], sep)
	return Decoded{strings.ToLower(parts[1])}
}

func decodeArray(msg string) Decoded {
	result := Decoded{}

	lengthStr, remainder, _ := strings.Cut(msg[1:], sep)
	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		logger.Error("Error decoding array %s: %v", msg, err)
		return result
	}

	// Empty array.
	if length == 0 {
		return result
	}

	re := regexp.MustCompile(".*\r\n.*\r\n")
	matches := re.FindAllString(remainder, -1)
	for _, part := range matches {
		result = append(result, DecodeRESP(part)...)
	}
	return result
}
