package main

import (
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/et-codes/codecrafters-redis-go/logging"
)

const (
	host         = "localhost"
	port         = "6379"
	pingCommand  = "*1\r\n$4\r\nping\r\n"
	pingResponse = "+PONG\r\n"
)

var logger = logging.New(logging.LevelDebug)

func main() {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		logger.Fatal("Failed to bind to port %s", port)
	}
	logger.Info("Listening on port %s...", port)

	for {
		conn, err := l.Accept()
		if err != nil {
			logger.Fatal("Error accepting connection: %v", err)
		}
		go handle(conn)
	}
}

func handle(conn io.ReadWriteCloser) {
	logger.Info("Connection initiated.")
	defer conn.Close()

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		logger.Error("Error reading connection: %v", err)
		return
	}
	logger.Debug("Received %d bytes: %v", n, buffer[:n])

	if strings.Contains(strings.ToLower(string(buffer)), "ping") {
		logger.Info("Ping received.")
		_, err := conn.Write([]byte(pingResponse))
		if err != nil {
			logger.Error("Error reading connection: %v", err)
			return
		}
	}
	logger.Info("Closing client connection.")
}
