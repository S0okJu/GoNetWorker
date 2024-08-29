package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	concurrencyLevel = 2 // Number of concurrent workers
	numRequests      = 5 // Total number of requests to make
)

// worker is a function that each worker goroutine will run
func worker(id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		fmt.Printf("Worker %d processing job %d\n", id, j)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Seed random number generator

	// Create a slice of channels, one for each worker
	jobChannels := make([]chan int, concurrencyLevel)
	for i := range jobChannels {
		jobChannels[i] = make(chan int, numRequests)
	}

	var wg sync.WaitGroup

	// Start worker goroutines, each with its own channel
	for w := 0; w < concurrencyLevel; w++ {
		wg.Add(1)
		go worker(w+1, jobChannels[w], &wg)
	}

	// Set up signal channel to listen for interrupt (Ctrl + C)
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-sigChannel:
				fmt.Println("\nReceived interrupt signal, closing channels...")
				// Close all job channels to stop workers
				for _, ch := range jobChannels {
					close(ch)
				}
				return
			default:
				// Randomly select a worker and a job
				w := rand.Intn(concurrencyLevel) // Random worker
				j := rand.Intn(numRequests)      // Random job

				// Randomly wait between 0 and 30 seconds
				delay := time.Duration(rand.Intn(6)) * time.Second
				fmt.Printf("Waiting for %v before assigning job %d to worker %d\n", delay, j, w+1)
				time.Sleep(delay)

				jobChannels[w] <- j
			}
		}
	}()

	// Wait for all workers to finish
	wg.Wait()
	fmt.Println("All workers have completed their tasks.")
}
