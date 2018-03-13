package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
	"strings"
)

var clients = make(map[string]chan *http.Request)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Error upgrading connection: %v\n", err)
		return
	}
	defer conn.Close()

	messageType, message, err := conn.ReadMessage()
	if err != nil {
		fmt.Printf("Error reading from Client: %v\n", err)
		return
	}

	if messageType != websocket.TextMessage || strings.ToLower(string(message)) == "ws" {
		fmt.Printf("Invalid message received from Client: %v\n", messageType)
		return
	}

	out := string(message)
	if len(out) == 0 {
		// Generate a random name for the user.
		out = uuid.Must(uuid.NewV4()).String()
	} else {
		// Format string appropriately
		out = strings.ToLower(strings.Replace(strings.Trim(out, " "), " ", "-", -1))
	}

	clients[out] = make(chan *http.Request)

	// TODO: Check if name exists in redis already and generate a new one if it does.
	url := "http://localhost/" + out
	err = conn.WriteMessage(websocket.TextMessage, []byte(url))

	for {
		select {
		case r := <-clients[out]:
			fmt.Printf("Received a message on endpoint: %v %v\n", out, r)
			// TODO: Pass request data to connected client so they can handle it however they like
			break
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	id := v["id"]
	clients[id] <- r
}

func main() {
	r := mux.NewRouter()
	r.StrictSlash(true)

	r.HandleFunc("/ws", websocketHandler).
		Methods("GET")

	r.HandleFunc("/{id}", handler)

	log.Fatal(http.ListenAndServe(":8080", r))
}
