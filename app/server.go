package main

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

type RespCommand struct {
	command string
	args    []string
}

type Store struct {
	data map[string]string
}

func (s Store) SET(key string, val string) {
	s.data[key] = val
}

func (s Store) GET(key string) (string, bool) {
	val, ok := s.data[key]
	return val, ok
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

	store := Store{
		data: make(map[string]string),
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleClient(conn, store)
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

func handleClient(conn net.Conn, store Store) {
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
		fmt.Println(respCommand.args)
		switch strings.ToLower(respCommand.command) {
		case "ping":
			conn.Write([]byte("+PONG\r\n"))
		case "echo":
			conn.Write([]byte(createBulkString(respCommand.args[0])))
		case "set":
			key, val := respCommand.args[0], respCommand.args[1]
			store.SET(key, val)
			if len(respCommand.args) > 2 {
				fmt.Println("3rd argument", respCommand.args[2])
				if strings.ToLower(respCommand.args[2]) == "px" {
					if len(respCommand.args) < 4 {
						fmt.Println("Insufficient arguments")
					}
					if len(respCommand.args) == 4 {
						expiryTime, err := strconv.ParseInt(respCommand.args[3], 10, 64)
						if err != nil {
							fmt.Println("Error:", err)
							return
						}
						go expire(store, key, expiryTime)
					}
				}
			}

			conn.Write([]byte("+OK\r\n"))
		case "get":
			key := respCommand.args[0]
			val, ok := store.GET(key)
			if !ok {
				fmt.Println("Key not found")
				conn.Write([]byte("$-1\r\n"))
			} else {
				conn.Write([]byte(createBulkString(val)))
			}
		default:
			fmt.Printf("unhandled command: %s\n", respCommand.command)

		}
	}
}
