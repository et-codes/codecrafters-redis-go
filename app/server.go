package main

import (
	"fmt"
	"net"
	"os"
)

const (
	redisHost = "localhost"
	redisPort = 6379
)

func main() {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", redisHost, redisPort))
	if err != nil {
		fmt.Printf("Failed to bind to port %d\n", redisPort)
		os.Exit(1)
	}
	_, err = l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
}
