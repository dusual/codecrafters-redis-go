package main

import (
	"fmt"
	"io"
	"log"
	"strings"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

type RespCommand struct {
	command string
	args    []string
}

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
		go handleClient(conn)
	}
}

func parseRespCommand(commandString string) (RespCommand, error) {
	lines := strings.Split(commandString, "\r\n")
	respCommand := RespCommand{}
	respCommand.command = lines[2]
	if len(lines) > 3 {
		for i, argLine := range lines[3:] {
			if i%2 == 1 {
				respCommand.args = append(respCommand.args, argLine)
			}
		}
	}
	return respCommand, nil
}

func createBulkString(s string) string {
	return fmt.Sprintf("$%d\r\n%s\r\n", len(s), s)
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
		respCommand, err := parseRespCommand(command)
		if err != nil {
			log.Println("Error parsing command: ", err.Error())
			break
		}
		switch strings.ToLower(respCommand.command) {
		case "ping":
			conn.Write([]byte("+PONG\r\n"))
		case "echo":
			conn.Write([]byte(createBulkString(respCommand.args[0])))
		default:
			fmt.Printf("unhandled command: %s\n", respCommand.command)

		}
	}
}
