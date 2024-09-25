package core

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
	channel chan Job
}

func newWorker() worker {
	return worker{channel: make(chan Job)}
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

// makeRequest simulates an HTTP request
func makeRequest(ctx context.Context, job Job, client *http.Client, sleepRange int) error {
	// Create the full URL
	url := job.Url

	// Sleep for a random time within the range
	sleepDuration := time.Duration(rand.Intn(sleepRange)) * time.Second
	time.Sleep(sleepDuration)

	// Perform the request based on the HTTP method
	var req *http.Request
	var err error

	switch strings.ToUpper(job.Method) {
	case "GET":
		req, err = http.NewRequestWithContext(ctx, "GET", url, nil)
		fmt.Println("GET request to", url)
	case "POST":
		body, _ := job.ConvertTo()
		req, err = http.NewRequestWithContext(ctx, "POST", url, body) // Add payload handling if needed
		req.Header.Set("Content-Type", "application/json")
		fmt.Println("POST request to", url)
	default:
		return fmt.Errorf("unsupported HTTP method: %s", job.Method)
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

// worker listens on a channel for tasks and processes them.
func work(ctx context.Context, id int, jobs <-chan Job, wg *sync.WaitGroup, client *http.Client, sleepRange int) {
	defer wg.Done()
	wg.Add(1)

	for {
		select {
		case j, ok := <-jobs:
			if !ok {
				// If the channel is closed, the worker stops.
				fmt.Printf("Worker %d shutting down.\n", id)
				return
			}
			// Execute the request
			fmt.Printf("Worker %d processing task\n", id)
			err := makeRequest(ctx, j, client, sleepRange)
			if err != nil {
				fmt.Printf("Worker %d encountered error: %v\n", id, err)
			}

		case <-ctx.Done():
			// Context was canceled, exit the worker
			fmt.Printf("Worker %d stopped due to context cancellation.\n", id)
			return
		}
	}
}

func (ws *Workers) addJob(jobs Jobs) {
	for {
		j := jobs[rand.Intn(len(jobs))]
		w := (*ws)[rand.Intn(len(*ws))]
		w.channel <- j
	}
}

// Start begins the worker tasks and listens for cancellation signals.
func (ws *Workers) Start(ctx context.Context, cfg *Config, wg *sync.WaitGroup) error {
	// Initialize HTTP client
	client := &http.Client{}

	// Parse the jobs
	parser := NewParser(*cfg)
	jobs, err := parser.Parse()
	if err != nil {
		return err
	}

	go ws.addJob(jobs)

	for i, w := range *ws {
		go work(ctx, i, w.channel, wg, client, cfg.GetSleepRange())
	}

	return nil
}

// Done waits for all workers to finish.
func (ws *Workers) Done() {
	for _, worker := range *ws {
		close(worker.channel)
	}
	fmt.Println("All workers are done.")
}
