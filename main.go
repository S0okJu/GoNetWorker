package main

import (
	"fmt"
	"github.com/s0okjug/gonetworker/controller"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

// Updated makeRequest function to accept http.Client as an argument
func makeRequest(work controller.Work, task controller.Task, client *http.Client, wg *sync.WaitGroup, sleepRange int) {
	defer wg.Done()

	// Create the full URL
	url := fmt.Sprintf("http://localhost:%d%s", work.Port, task.Path)

	// Sleep for a random time within the range
	sleepDuration := time.Duration(rand.Intn(sleepRange)) * time.Second
	time.Sleep(sleepDuration)

	// Perform the request based on the HTTP method
	var req *http.Request
	var err error

	switch strings.ToUpper(task.Method) {
	case "GET":
		req, err = http.NewRequest("GET", url, nil)
		fmt.Println("GET request to", url)
	case "POST":
		fmt.Println("POST request to", url)
		// Uncomment and handle payload if necessary
		// payload := []byte()
		// req, err = http.NewRequest("POST", url, bytes.NewBuffer(payload))
		// req.Header.Set("Content-Type", "application/json")
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

	fmt.Printf("Response from %s: %s\n", url, resp.Status)
}

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	reader := controller.NewReader("./test.json")
	config, err := reader.GetConfig()
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		return
	}

	// Create a wait group for goroutines
	var wg sync.WaitGroup

	// Create a channel to listen for interrupt signals (Ctrl+C)
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// Infinite loop to send random requests continuously
	go func() {
		for {
			select {
			case <-stopChan:
				fmt.Println("Received interrupt signal. Shutting down...")
				return
			default:
				// Randomly select a work and a task
				work := config.Works[rand.Intn(len(config.Works))]
				task := work.Tasks[rand.Intn(len(work.Tasks))]

				// Start the request in a new goroutine
				wg.Add(1)
				client := &http.Client{}
				go makeRequest(work, task, client, &wg, config.Settings.SleepRange)

				// Sleep for a random duration before the next request
				randomSleep := time.Duration(rand.Intn(config.Settings.SleepRange)) * time.Second
				time.Sleep(randomSleep)
			}
		}
	}()

	// Block main goroutine until interrupt signal is received
	<-stopChan
	fmt.Println("Program terminated.")
}
