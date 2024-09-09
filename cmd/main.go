package main

import (
	"encoding/json"
	"fmt"
	"github.com/s0okjug/gonetworker/controller"
	"io"
	"net/http"
	"os"
	_ "strconv"
	"strings"
	"sync"
	"time"
)

// Updated makeRequest function to accept http.Client as an argument
func makeRequest(work controller.Work, task controller.Task, client *http.Client, wg *sync.WaitGroup, sleepRange int) {
	defer wg.Done()

	// Create the full URL
	url := fmt.Sprintf("http://localhost:%d%s", work.Port, task.Path)

	// Sleep for a random time within the range
	sleepDuration := time.Duration(sleepRange) * time.Second
	time.Sleep(sleepDuration)

	// Perform the request based on the HTTP method
	var req *http.Request
	var err error

	switch strings.ToUpper(task.Method) {
	case "GET":
		req, err = http.NewRequest("GET", url, nil)
	case "POST":
		fmt.Println("post")
		//payload := []byte()
		//req, err = http.NewRequest("POST", url, bytes.NewBuffer(payload))
		//req.Header.Set("Content-Type", "application/json")
	default:
		fmt.Printf("Unsupported HTTP method: %s\n", task.Method)
		return
	}

	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	fmt.Printf("Response from %s: %s\n", url, resp.Status)
	fmt.Printf("Response body: %s\n", string(body))
}

func main() {
	// Open the JSON file
	jsonFile, err := os.Open("data.json")
	if err != nil {
		fmt.Printf("Error opening JSON file: %v\n", err)
		return
	}
	defer jsonFile.Close()

	// Read the file's content
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Printf("Error reading JSON file: %v\n", err)
		return
	}

	// Parse the JSON input
	var config controller.Config
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

	// Create a wait group for goroutines
	var wg sync.WaitGroup

	// Loop through works and tasks
	for _, work := range config.Works {
		for _, task := range work.Tasks {
			wg.Add(1)
			// Use a default HTTP client for the requests
			client := &http.Client{}
			go makeRequest(work, task, client, &wg, config.Settings.SleepRange)
		}
	}

	// Wait for all goroutines to complete
	wg.Wait()
	fmt.Println("All requests completed.")
}
