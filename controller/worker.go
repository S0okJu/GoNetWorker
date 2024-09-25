package controller

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

// worker is a struct that contains a channel to signal when the worker's task is done.
type worker struct {
	channel chan bool
}

func newWorker() worker {
	return worker{channel: make(chan bool)}
}

// Workers is an array of workers.
type Workers []worker

// NewWorkers creates a list of workers (CCU).
func NewWorkers(cnt int) (*Workers, error) {
	var workers Workers
	for i := 0; i < cnt; i++ {
		workers = append(workers, newWorker())
	}
	return &workers, nil
}

// Assign tasks to workers and handle cancellation context.
func (ws *Workers) assign(ctx context.Context, task func()) {
	var wg sync.WaitGroup

	for _, w := range *ws {
		wg.Add(1)
		go func(wk worker) {
			defer wg.Done()

			select {
			case <-ctx.Done(): // Handle context cancellation
				fmt.Println("Task cancelled")
				return
			default:
				task() // Execute the task function if no cancellation
				wk.channel <- true
			}
		}(w)
	}

	wg.Wait()
}

// Updated makeRequest function with context for cancellation
func makeRequest(ctx context.Context, work Work, task Task, client *http.Client, sleepRange int) error {
	// Create the full URL
	url := fmt.Sprintf("%s:%d%s", work.Uri, work.Port, task.Path)

	// Sleep for a random time within the range
	sleepDuration := time.Duration(rand.Intn(sleepRange)) * time.Second
	time.Sleep(sleepDuration)

	// Perform the request based on the HTTP method
	var req *http.Request
	var err error

	switch strings.ToUpper(task.Method) {
	case "GET":
		req, err = http.NewRequestWithContext(ctx, "GET", url, nil)
		fmt.Println("GET request to", url)
	case "POST":
		req, err = http.NewRequestWithContext(ctx, "POST", url, nil) // Add payload handling if needed
		req.Header.Set("Content-Type", "application/json")
		fmt.Println("POST request to", url)
	default:
		return fmt.Errorf("unsupported HTTP method: %s", task.Method)
	}

	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	fmt.Printf("Response from %s: %s\n", url, resp.Status)
	return nil
}

// Start begins the worker tasks and listens for cancellation signals.
func (ws *Workers) Start(ctx context.Context, cfg *Config) {

	// Initialize HTTP client
	client := &http.Client{}

	// Run tasks in workers
	go func() {
		ws.assign(ctx, func() {
			// Randomly select a work and a task
			work := cfg.Works[rand.Intn(len(cfg.Works))]
			task := work.Tasks[rand.Intn(len(work.Tasks))]

			// Perform the request
			err := makeRequest(ctx, work, task, client, cfg.GetSleepRange())
			if err != nil {
				fmt.Printf("Error in request: %v\n", err)
			}
		})
	}()
}

// Done waits for all workers to finish.
func (ws *Workers) Done() {
	for _, worker := range *ws {
		<-worker.channel // Wait for each worker to complete
	}
	fmt.Println("All workers are done.")
}
