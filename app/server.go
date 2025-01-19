package main

import (
	"fmt"
	"net"
	"os"
)

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
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			return
		}
		fmt.Println(string(buf))
		conn.Write([]byte("+PONG\r\n"))
	}
}
