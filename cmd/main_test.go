package main

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/s0okjug/gonetworker/core"
)

// Mock function to simulate GET request test
func TestMakeRequest_GET(t *testing.T) {
	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the method is GET
		if r.Method != http.MethodGet {
			t.Errorf("Expected 'GET' method, got '%s'", r.Method)
		}
		// Simulate a successful response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "GET request successful"}`))
	}))
	defer mockServer.Close()

	// Define a mock work and task
	work := core.Work{
		Port: 8080, // Port is not relevant here since we're using mockServer.URL
		Tasks: []core.Task{
			{Method: "GET", Path: mockServer.URL},
		},
	}

	// Create a wait group
	var wg sync.WaitGroup
	wg.Add(1)

	// Create a mock HTTP client
	client := &http.Client{}

	// Call the makeRequest function with the mock client
	go makeRequest(work, work.Tasks[0], client, &wg, 0)

	// Wait for the goroutine to finish
	wg.Wait()
}
