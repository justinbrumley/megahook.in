package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"testing"
)

func TestServer(t *testing.T) {
	origin := "http://localhost/"
	url := "ws://localhost:8080/ws"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		t.Errorf("Error connecting to Server: %v\n", err)
		return
	}

	// Start a test connection
	if _, err := ws.Write([]byte("Test Client Here")); err != nil {
		t.Errorf("Error writing to Server: %v\n", err)
		return
	}

	var msg = make([]byte, 255)
	n, err := ws.Read(msg)
	if err != nil {
		if err == io.EOF {
			return
		}

		t.Errorf("Error reading from Server: %v\n", err)
		return
	}

	fmt.Printf("Message received from Server: %v %v\n", n, string(msg[:n]))
}
