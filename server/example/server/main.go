package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Message struct {
	Text string `json:"text"`
}

func main() {
	http.HandleFunc("/", logMiddleware(homeHandler))
	//http.HandleFunc("/hello", logMiddleware(helloHandler))
	http.HandleFunc("/echo", logMiddleware(echoHandler))

	port := 8080
	log.Printf("Server is starting on port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func logMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		log.Printf("INFO: Request received: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("INFO: Request completed: %s %s, duration: %v", r.Method, r.URL.Path, time.Since(startTime))
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO: Serving home page")
	fmt.Fprintf(w, "Welcome to the home page!")
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Printf("WARN: Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}
	log.Printf("INFO: Saying hello to %s", name)
	fmt.Fprintf(w, "Hello, %s!", name)
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("WARN: Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERROR: Failed to read request body: %v", err)
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	// Don't forget to close the original body
	r.Body.Close()

	// Log the body content
	log.Printf("INFO: Request body: %s", string(body))

	// Create a new reader with the body for JSON decoding
	bodyReader := bytes.NewReader(body)

	var msg Message
	err = json.NewDecoder(bodyReader).Decode(&msg)
	if err != nil {
		log.Printf("ERROR: Error decoding JSON: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("INFO: Echoing message: %s", msg.Text)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}
