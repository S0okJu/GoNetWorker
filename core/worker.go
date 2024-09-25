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

// Assign tasks to workers by sending tasks through channels.
func (ws *Workers) assign(ctx context.Context, cfg *Config, wg *sync.WaitGroup, client *http.Client) {
	for _, w := range *ws {
		wg.Add(1)
		go func(wk worker) {
			defer wg.Done()

			for {
				select {
				case job := <-wk.channel:
					// Execute the task received via channel
					err := makeRequest(ctx, job, client, cfg.GetSleepRange())
					if err != nil {
						fmt.Printf("Error in request: %v\n", err)
					}

					// Sleep before next task to avoid a tight loop
					time.Sleep(1 * time.Second)

				case <-ctx.Done():
					// Stop the worker when context is canceled
					fmt.Println("Worker stopped due to context cancellation")
					return
				}
			}
		}(w)
	}
}

// Updated makeRequest function with context for cancellation
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

// Start begins the worker tasks and listens for cancellation signals.
func (ws *Workers) Start(ctx context.Context, cfg *Config) {
	// Initialize HTTP client
	client := &http.Client{}
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Stopping task assignment due to context cancellation.")
				return
			default:
				// Randomly select a work and a task
				work := cfg.Works[rand.Intn(len(cfg.Works))]
				task := work.Tasks[rand.Intn(len(work.Tasks))]

				// Randomly send the task to one of the workers' channels
				worker := (*ws)[rand.Intn(len(*ws))]
				worker.channel <-
			}
		}
	}()
}

// Done waits for all workers to finish.
func (ws *Workers) Done() {
	for _, worker := range *ws {
		close(worker.channel)
	}
	fmt.Println("All workers are done.")
}
