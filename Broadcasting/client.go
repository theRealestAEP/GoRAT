package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

const (
	clientPort = 8080
)

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", clientPort))
	if err != nil {
		fmt.Println("Error starting client:", err)
		os.Exit(1)
	}

	fmt.Printf("Client listening on port %d...\n", clientPort)

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection:", err)
		os.Exit(1)
	}
	defer conn.Close()

	go func() {
		_, _ = io.Copy(conn, os.Stdin)
	}()

	_, _ = io.Copy(os.Stdout, conn)
}
