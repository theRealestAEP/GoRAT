package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

const (
	multicastAddr = "224.0.0.1"
	multicastPort = 9999
)

func main() {
	serverIP, serverPort := findServer()

	fmt.Printf("Found server at %s:%d\n", serverIP, serverPort)

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, serverPort))
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	done := make(chan struct{})

	go func() {
		_, err := io.Copy(os.Stdout, conn)
		if err != nil {
			fmt.Println("Error reading from server:", err)
		}
		close(done)
	}()

	go func() {
		_, err := io.Copy(conn, os.Stdin)
		if err != nil {
			fmt.Println("Error sending input to server:", err)
		}
		close(done)
	}()

	<-done
}

func findServer() (string, int) {
	maxAttempts := 5
	attempts := 0

	for attempts < maxAttempts {
		addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", multicastAddr, multicastPort))
		if err != nil {
			fmt.Println("Error resolving multicast address:", err)
			os.Exit(1)
		}

		conn, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			fmt.Println("Error setting up UDP multicast:", err)
			os.Exit(1)
		}
		defer conn.Close()

		_, err = conn.Write([]byte("CLIENT_SEARCHING"))
		if err != nil {
			fmt.Println("Error sending multicast:", err)
			os.Exit(1)
		}

		buf := make([]byte, 1024)
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		n, srcAddr, err := conn.ReadFromUDP(buf)
		if err == nil {
			message := string(buf[:n])
			parts := strings.Split(message, ":")
			if len(parts) == 2 && parts[0] == "SERVER_READY" {
				serverIP := srcAddr.IP.String()
				serverPort := 0
				fmt.Sscanf(parts[1], "%d", &serverPort)

				return serverIP, serverPort
			}
		} else if err, ok := err.(net.Error); ok && err.Timeout() {
			fmt.Printf("Attempt %d: Server not found. Retrying...\n", attempts+1)
		} else {
			fmt.Println("Error reading server response:", err)
			os.Exit(1)
		}

		attempts++
	}

	fmt.Println("Server not found after multiple attempts.")
	os.Exit(1)
	return "", 0
}
