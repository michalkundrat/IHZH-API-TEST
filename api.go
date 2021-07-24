package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

var Messages []Message

// Boring Stuff

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
	http.HandleFunc("/messages", showAllMessages)
	http.HandleFunc("/send", sendMessage)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	Messages = []Message{
		Message{Username: "Test1", Content: "Hello World #1"},
	}
	handleRequests()
}
