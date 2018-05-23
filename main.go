package main

import (
	"bufio"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Request struct {
	Headers map[string][]string `json:"headers,omitempty"`
	Method  string              `json:"method,omitempty"`
	Body    string              `json:"body,omitempty"`
	Query   url.Values          `json:"query,omitempty"`
}

var clients = make(map[string]chan *Request)

func checkOrigin(r *http.Request) bool {
	return true
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
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
	if _, ok := clients[out]; ok || len(out) == 0 {
		// Generate a random name for the user if it's taken or they didn't provide one
		out = uuid.Must(uuid.NewV4()).String()
	} else {
		// Format string appropriately
		out = strings.ToLower(strings.Replace(strings.Trim(out, " "), " ", "-", -1))
	}

	clients[out] = make(chan *Request)
	defer delete(clients, out)

	// TODO: Check if name exists in redis already and generate a new one if it does.
	url := "https://megahook.in/m/" + out
	err = conn.WriteMessage(websocket.TextMessage, []byte(url))

	fmt.Printf("Listening for request at %v\n", url)
	for {
		select {
		case r := <-clients[out]:
			fmt.Printf("Received a message on endpoint: %v\n", out)
			conn.WriteJSON(r)
			break
		}
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Home page view!")
	http.ServeFile(w, r, "index.html")
}

func handler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	id := v["id"]

	reader := bufio.NewReader(r.Body)
	body := ""

	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				body += s
				break
			}

			fmt.Printf("Error reading from body: %v\n", err)
			break
		}

		body += s
	}

	req := &Request{
		Headers: map[string][]string(r.Header),
		Method:  r.Method,
		Body:    body,
		Query:   r.URL.Query(),
	}

	clients[id] <- req
}

func main() {
	r := mux.NewRouter()
	r.StrictSlash(true)

	r.HandleFunc("/ws", websocketHandler).
		Methods("GET")

	r.HandleFunc("/", homeHandler).
		Methods("GET")

	r.HandleFunc("/m/{id}", handler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	// Start HTTPS server on different Goroutine
	go func() {
		log.Fatal(http.ListenAndServeTLS(":443", "/etc/letsencrypt/live/megahook.in/fullchain.pem", "/etc/letsencrypt/live/megahook.in/privkey.pem", r))
	}()

	log.Fatal(http.ListenAndServe(":80", r))
}
