package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func main() {
	// Set up the handler for "/users/"
	http.HandleFunc("/users/", getUserID)

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func getUserID(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s %s", r.Method, r.RequestURI)

	userID := strings.TrimPrefix(r.URL.Path, "/users/")
	if userID == "" {
		http.Error(w, "User ID not provided", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "User ID: %s", userID)
}
