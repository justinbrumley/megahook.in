package main

import (
	"net"
	"fmt"
	"bufio"
	"testing"
	"github.com/gorilla/websocket"
)

func TestServer(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		return
	}

	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	status, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading from server: %v\n", err)
		return
	}

	fmt.Printf("Server Status: %v\n", status)
}