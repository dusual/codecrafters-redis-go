package main

import (
	"fmt"
	"io"
	"log"
	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	// Ensure we close the connection after we're done
	defer conn.Close()

	// Read data

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		command := string(buf[:n])

		if err == io.EOF {
			fmt.Println("Connection closed")
			break
		}

		if err != nil {
			return
			break
		}

		log.Println("Received data", buf[:n])
		log.Println("command", command)

		// Compile the regex pattern

		// Check if the multiline string matches the regex pattern
		if command == "*1\r\n$4\r\nping\r\n" {
			conn.Write([]byte("+PONG\r\n"))
		} else {
			fmt.Println("Multiline string does not match the pattern")
		}
	}
}
