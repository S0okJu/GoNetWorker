package core

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func request(job *Job, wg *sync.WaitGroup) error {
	defer wg.Done()
	var req *http.Request
	var err error

	switch job.Method {
	case "GET":
		req, err = http.NewRequest("GET", job.Url, nil)
		fmt.Println("GET request to", job.Url)
	case "POST":
		fmt.Println("POST request to", job.Url)
	default:
		fmt.Printf("Unsupported HTTP method: %s\n", job.Method)
		return fmt.Errorf("Unsupported HTTP method: %s", job.Method)
	}

	if err != nil {
		return err
	}

	var client http.Client
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}

func Execute(cfg Config) error {
	parser := NewParser(cfg)
	jobs, err := parser.Parse()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	// Create a channel to listen for interrupt signals (Ctrl+C)
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-stopChan:
				fmt.Println("Interrupt signal received")
				return
			default:
				wg.Add(1)

				// Randomly select a job
				job := jobs[rand.Intn(len(jobs))]

				go func() {
					err := request(&job, &wg)
					if err != nil {
						fmt.Println("Error making request:", err)
					}
				}()
				// Sleep for a random duration before the next request
				randomSleep := time.Duration(rand.Intn(cfg.Settings.SleepRange)) * time.Second
				time.Sleep(randomSleep)
			}
		}
	}()
	<-stopChan
	fmt.Println("Shutting down...")
	return nil
}
