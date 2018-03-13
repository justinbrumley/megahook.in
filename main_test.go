package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	dialer := &websocket.Dialer{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	header := &http.Header{}
	header.Add("Origin", "http://localhost/")
	conn, _, err := dialer.Dial("ws://localhost:8080/ws", *header)
	if err != nil {
		fmt.Printf("Error establishing connection: %v\n", err)
		return
	}

	defer conn.Close()

	// Write message to server so it knows what name to listen on
	conn.WriteMessage(websocket.TextMessage, []byte("test-client"))

	// Next message from server will be the url to use for webhooks
	_, _, err = conn.ReadMessage()
	if err != nil {
		fmt.Printf("Error reading message from server: %v\n", err)
		return
	}

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("Error reading message from server: %v\n", err)
			break
		}

		if messageType == websocket.CloseMessage {
			conn.Close()
			break
		}

		if messageType == websocket.TextMessage {
			fmt.Printf("Received message from server: %v\n", string(message))
		}
	}
}
