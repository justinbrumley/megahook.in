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
	"time"
)

// Request sent to webhook from outside source.
type Request struct {
	Headers  map[string][]string `json:"headers,omitempty"`
	Method   string              `json:"method,omitempty"`
	Body     string              `json:"body,omitempty"`
	Query    url.Values          `json:"query,omitempty"`
	Response chan *Response      `json:"-"`
}

// Response from the megahook client. Will be forwarded back to outside source.
type Response struct {
	Headers    map[string][]string `json:"headers,omitempty"`
	Body       string              `json:"body,omitempty"`
	StatusCode int                 `json:"status_code,omitempty"`
}

type ClientOptions struct {
	Name  string `json:"name"`
	Token string `json:"token"`
	Track bool   `json:"track"`
}

const (
	readTimeout  = time.Second * 30
	writeTimeout = time.Second * 30
	pingPeriod   = time.Second * 5

	port    = "8080"
	version = "0.1.0"
)

var namespaces = make(map[string](map[string]chan *Request))

func checkOrigin(r *http.Request) bool {
	return true
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

const STATIC_DIR = "/static/"

// Ensure that namespace client map is initialized
func ensureNamespace(n string) map[string]chan *Request {
	if _, ok := namespaces[n]; !ok {
		namespaces[n] = make(map[string]chan *Request)
	}

	return namespaces[n]
}

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

	// Expecting connection info from client
	opts := &ClientOptions{}
	err = conn.ReadJSON(&opts)
	if err != nil {
		fmt.Printf("Error reading opts from client: %v\n", err)
		return
	}

	ns := ""
	if opts.Token != "" {
		token, err := getTokenNamespace(opts.Token)
		if err != nil {
			fmt.Printf("Error getting token namespace: %v %v", opts.Token, err)
			return
		}

		if token != nil {
			ns = token.Namespace
		}
	}

	clients := ensureNamespace(ns)

	out := opts.Name
	if _, ok := clients[out]; ok || len(out) == 0 {
		// Generate a random name for the user if it's taken or they didn't provide one
		out = uuid.Must(uuid.NewV4(), nil).String()
	} else {
		// Format string appropriately
		out = strings.ToLower(strings.Replace(strings.Trim(out, " "), " ", "-", -1))
	}

	clients[out] = make(chan *Request)
	defer (func() {
		fmt.Println("Closing connection...")
		close(clients[out])
		delete(clients, out)

		if len(clients) == 0 {
			fmt.Printf("Removing namespace clients map: %v\n", ns)
			delete(namespaces, ns)
		}
	})()

	// Write the public URL back to client
	url := "https://api.megahook.in/m/" + out
	if ns != "" {
		url = fmt.Sprintf("https://%v.api.megahook.in/m/%v", ns, out)
	}

	conn.SetWriteDeadline(time.Now().Add(writeTimeout))
	if err = conn.WriteMessage(websocket.TextMessage, []byte(url)); err != nil {
		return
	}

	close := make(chan bool)
	responseChan := make(chan *Response)
	go (func() {
		for {
			// Read messages from client
			r := &Response{}
			if err := conn.ReadJSON(&r); err != nil {
				// Connection closed or errored out
				switch err.(type) {
				case *websocket.CloseError:
					c := err.(*websocket.CloseError).Code
					if c != 1006 && c != 1000 { // Normal and Abnormal Closures are okay
						fmt.Printf("Failed to read message: %v\n", err.(*websocket.CloseError).Code)
					}
				}

				close <- true
				return
			}

			responseChan <- r
		}
	})()

	ticker := time.NewTicker(pingPeriod)

	fmt.Printf("Listening for requests at %v\n", url)

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

			// Wait for response
			respTicker := time.NewTicker(readTimeout)
			select {
			case <-respTicker.C:
				fmt.Println("Response took too long. Not waiting anymore")
				continue

			case resp := <-responseChan:
				r.Response <- resp
			}
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	id := v["id"]

	ns := strings.Split(r.Host, ".")[0]
	if ns == "api" || ns == "megahook" {
		ns = ""
	}

	clients := ensureNamespace(ns)

	if _, ok := clients[id]; !ok {
		w.WriteHeader(404)
		fmt.Fprint(w, "Not Found")
		return
	}

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
		Headers:  map[string][]string(r.Header),
		Method:   r.Method,
		Body:     body,
		Query:    r.URL.Query(),
		Response: make(chan *Response),
	}

	rec := &Record{
		Request:   req,
		Timestamp: time.Now().Unix(),
	}

	clients[id] <- req

	// Wait for response from client
	ticker := time.NewTicker(readTimeout)
	select {
	case <-ticker.C:
		break

	case response := <-req.Response:
		for key, headers := range response.Headers {
			for _, value := range headers {
				w.Header().Set(key, value)
			}
		}

		rec.Response = response

		w.WriteHeader(response.StatusCode)
		fmt.Fprint(w, response.Body)
	}
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	fmt.Fprintf(w, "%v", version)
}

func main() {
	initRedis()
	err := initDB()
	if err != nil {
		fmt.Printf("Error initializing DB connection: %v\n", err)
	}

	r := mux.NewRouter()
	r.StrictSlash(true)

	r.HandleFunc("/ws", websocketHandler).
		Methods("GET")

	r.HandleFunc("/hooks/{name}/history", historyHandler).
		Methods("GET")

	r.HandleFunc("/", indexHandler)

	r.HandleFunc("/register", registerHandler)

	r.HandleFunc("/m/{id}", handler)

	r.HandleFunc("/version", versionHandler).
		Methods("GET")

	fmt.Printf("Starting server on port %v\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
