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
	"os"
	"strings"
	"time"
)

type Request struct {
	Headers map[string][]string `json:"headers,omitempty"`
	Method  string              `json:"method,omitempty"`
	Body    string              `json:"body,omitempty"`
	Query   url.Values          `json:"query,omitempty"`
}

const (
	readTimeout  = time.Second * 30
	writeTimeout = time.Second * 30
	pingPeriod   = time.Second * 5
)

var clients = make(map[string]chan *Request)

func checkOrigin(r *http.Request) bool {
	return true
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

const STATIC_DIR = "/static/"

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Error upgrading connection: %v\n", err)
		return
	}
	defer conn.Close()

	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(readTimeout))
		return nil
	})

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
	defer (func() {
		fmt.Println("Closing connection...")
		close(clients[out])
		delete(clients, out)
	})()

	url := "https://megahook.in/m/" + out
	conn.SetWriteDeadline(time.Now().Add(writeTimeout))
	if err = conn.WriteMessage(websocket.TextMessage, []byte(url)); err != nil {
		return
	}

	close := make(chan bool)
	go (func() {
		// Read Messages
		resp := http.Response{}
		if err := conn.ReadJSON(&resp); err != nil {
			close <- true
			return
		}

		fmt.Printf("Response: %v\n", resp)
	})()

	ticker := time.NewTicker(pingPeriod)
	fmt.Printf("Listening for request at %v\n", url)
	for {
		select {
		case <-close:
			return
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				fmt.Printf("Error sending ping: %v\n", err)
				return
			}
			break
		case r, ok := <-clients[out]:
			if !ok {
				fmt.Println("Client not ok")
				return
			}

			conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if err := conn.WriteJSON(r); err != nil {
				fmt.Printf("Failed to write message: %v\n", err)
				return
			}
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
	environment := os.Getenv("ENV")

	r := mux.NewRouter()
	r.StrictSlash(true)

	r.HandleFunc("/ws", websocketHandler).
		Methods("GET")

	r.HandleFunc("/", homeHandler).
		Methods("GET")

	r.HandleFunc("/m/{id}/inspect", homeHandler).
		Methods("GET")

	r.HandleFunc("/m/{id}", handler)

	r.PathPrefix(STATIC_DIR).Handler(http.StripPrefix(STATIC_DIR, http.FileServer(http.Dir("." + STATIC_DIR))))
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./dist/"))))

	port := "8080"

	// Start HTTPS server on different Goroutine
	if environment != "development" {
		port = "80"

		go func() {
			log.Fatal(http.ListenAndServeTLS(":443", "/etc/letsencrypt/live/megahook.in/fullchain.pem", "/etc/letsencrypt/live/megahook.in/privkey.pem", r))
		}()
	}

	fmt.Printf("Starting server on port %v\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
