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

	// Start the server on port 8080, and log the error if it fails
	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// Handler function for the "/users/{user_id}" route
func getUserID(w http.ResponseWriter, r *http.Request) {
	// Log the incoming request
	log.Printf("Received request: %s %s", r.Method, r.RequestURI)

	// Extract the user_id from the URL path
	userID := strings.TrimPrefix(r.URL.Path, "/users/")
	if userID == "" {
		http.Error(w, "User ID not provided", http.StatusBadRequest)
		return
	}

	// Log the extracted user_id
	log.Printf("Extracted user ID: %s", userID)

	// Respond with the user_id
	fmt.Fprintf(w, "User ID: %s", userID)
}
