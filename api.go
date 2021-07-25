package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}
type SSEChannel struct {
	Clients  []chan string
	Notifier chan string
}

var Messages []Message
var sseChannel SSEChannel

// SSE Stuff
func broadcaster(done <-chan interface{}) {
	fmt.Println("Broadcaster Started.")
	for {
		select {
		case <-done:
			return
		case data := <-sseChannel.Notifier:
			for _, channel := range sseChannel.Clients {
				channel <- data
			}
		}
	}
}
func logHTTPRequest(w http.ResponseWriter, r *http.Request) {
	buf := new(strings.Builder)
	if _, err := io.Copy(buf, r.Body); err != nil {
		fmt.Printf("Error: %v", err)
	}
	method := r.Method

	logMsg := fmt.Sprintf("Method: %v, Body: %v", method, buf.String())
	fmt.Println(logMsg)
	sseChannel.Notifier <- logMsg
}

// Endpoint: /sse
func endpointSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Connection doesnot support streaming", http.StatusBadRequest)
		return
	}

	sseChan := make(chan string)
	sseChannel.Clients = append(sseChannel.Clients, sseChan)

	d := make(chan interface{})
	defer close(d)
	defer fmt.Println("Closing channel.")

	for {
		select {
		case <-d:
			close(sseChan)
			return
		case data := <-sseChan:
			fmt.Printf("data: %v \n\n", data)
			fmt.Fprintf(w, "data: %v \n\n", data)
			flusher.Flush()
		}
	}

}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}
func showAllMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Request-Method", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	fmt.Println("Endpoint Hit: messages")
	json.NewEncoder(w).Encode(Messages)
}
func sendMessage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: send")
	r.ParseForm()
	var message Message
	if r.FormValue("username") == "" && r.FormValue("content") == "" {
		reqBody, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(reqBody, &message)
	} else {
		message.Username = r.FormValue("username")
		message.Content = r.FormValue("content")
		http.Redirect(w, r, "localhost:5500/index.html", http.StatusSeeOther)
	}
	Messages = append(Messages, message)
	json.NewEncoder(w).Encode(message)
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/sse", endpointSSE)
	http.HandleFunc("/messages", showAllMessages)
	http.HandleFunc("/send", sendMessage)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	sseChannel = SSEChannel{
		Clients:  make([]chan string, 0),
		Notifier: make(chan string),
	}
	Messages = []Message{
		Message{Username: "Test1", Content: "Hello World #1"},
	}
	handleRequests()
}
