package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/s0okjug/gonetworker/controller"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Function to make a request
func makeRequest(work controller.Work, task controller.Task, wg *sync.WaitGroup, sleepRange int) {
	defer wg.Done()

	// Create the full URL
	url := fmt.Sprintf(strconv.Itoa(work.Port), task.Path)

	// Sleep for a random time within the range
	sleepDuration := time.Duration(rand.Intn(sleepRange)) * time.Second
	time.Sleep(sleepDuration)

	// Perform the request based on the HTTP method
	var req *http.Request
	var err error

	switch strings.ToUpper(task.Method) {
	case "GET":
		req, err = http.NewRequest("GET", url, nil)
	case "POST":
		// In case of POST, we assume some random JSON data as a payload
		payload := []byte(task.Body)
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
	default:
		fmt.Printf("Unsupported HTTP method: %s\n", task.Method)
		return
	}

	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	// Print the response status and body
	fmt.Printf("Response from %s: %s\n", url, resp.Status)
	fmt.Printf("Response body: %s\n", string(body))
}
func main() {
	// Open the JSON file
	jsonFile, err := os.Open("data.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)

	// Read the file's content
	// Parse the JSON input
	var config controller.Config
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Create a wait group for goroutines
	var wg sync.WaitGroup

	// Loop through works and tasks
	for _, work := range config.Works {
		for _, task := range work.Tasks {
			wg.Add(1)
			go makeRequest(work, task, &wg, config.Settings.SleepRange)
		}
	}

	// Wait for all goroutines to complete
	wg.Wait()
	fmt.Println("All requests completed.")
}
