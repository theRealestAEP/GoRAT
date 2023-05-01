package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"

	"github.com/creack/pty"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Server listening on port 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	shell := exec.Command("bash")
	pty, err := pty.Start(shell)
	if err != nil {
		fmt.Println("Error starting PTY:", err)
		return
	}
	defer pty.Close()

	done := make(chan struct{})
	go func() {
		io.Copy(pty, conn)
		close(done)
	}()

	go func() {
		io.Copy(conn, pty)
		close(done)
	}()

	<-done
}
