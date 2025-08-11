package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

func parseRespCommand(reader *bufio.Reader) ([]string, error) {
	line, err := reader.ReadString('\n')

	if err != nil {
		return nil, err
	}
	line = strings.TrimRight(line, "\r\n")
	switch {
	case strings.HasPrefix(line, "*"):
		numElems, _ := strconv.Atoi(line[1:])
		parts := make([]string, 0, numElems)
		raiseError := false
		errorMessage := ""
		for i := 0; i < numElems; i++ {
			// read $len

			lenLine, _ := reader.ReadString('\n')
			lenLine = strings.TrimRight(lenLine, "\r\n")

			if !strings.HasPrefix(lenLine, "$") {
				fmt.Println("unexpected bulk string header:", lenLine)
				raiseError = true
				errorMessage = "unexpected bulk string header: " + lenLine
				break

			}
			strLen, _ := strconv.Atoi(lenLine[1:])

			// read actual string
			buf := make([]byte, strLen+2)
			_, err := io.ReadFull(reader, buf)
			if err != nil {
				raiseError = true
				errorMessage = err.Error()
				break
			}
			parts = append(parts, string(buf[:strLen]))
		}
		if raiseError {
			return nil, errors.New(errorMessage)
		}
		return parts, nil
	default:
		return nil, errors.New("Could not parse bulk String")
	}
}

func writeBulkString(w io.Writer, s string) {
	fmt.Fprintf(w, "$%d\r\n%s\r\n", len(s), s)
}

func acceptConns(l net.Listener, acceptChan chan net.Conn) {
	// Need a for loop to handle continuously handle connections
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection ", err.Error())
			// os.Exit(1)
			continue
		}
		acceptChan <- conn
	}
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
	connChannel := make(chan net.Conn)
	go acceptConns(l, connChannel)

	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	for {
		conn := <-connChannel
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	// Ensure we close the connection after we're done
	// Read data
	defer conn.Close()
	for {
		reader := bufio.NewReader(conn)
		//cmd, err := parseRespCommand(reader)
		cmdParts, err := parseRespCommand(reader)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			return
		}
		switch strings.ToUpper(cmdParts[0]) {
		case "PING":
			conn.Write([]byte("+PONG\r\n"))
		case "ECHO":
			if len(cmdParts) < 2 {
				conn.Write([]byte("-ERR wrong number of arguments for 'echo' command\r\n"))
				continue
			}
			writeBulkString(conn, cmdParts[1])
		default:
			conn.Write([]byte("-ERR unknown command\r\n"))
		}
	}
	fmt.Println("Connection closed")
}
