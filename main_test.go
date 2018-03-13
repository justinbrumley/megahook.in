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
	_, _, err := dialer.Dial("ws://localhost:8080/ws", *header)
	if err != nil {
		fmt.Printf("Error establishing connection: %v\n", err)
		return
	}
}
