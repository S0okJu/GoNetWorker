package main

import (
	"context"
	"fmt"
	"github.com/s0okju/gonetworker/core"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	reader := core.NewReader("./test.json")
	config, err := reader.GetConfig()
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		return
	}
	ws, err := core.NewWorkers(config.GetCcuMax())
	if err != nil {
		fmt.Printf("Error creating workers: %v\n", err)
		return
	}
	// Create a context that can be canceled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up channel to listen for OS signals (e.g., Ctrl+C)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signalChan
		fmt.Printf("Received signal: %v. Shutting down...\n", sig)
		cancel()
	}()

	// Start the worker pool
	err = ws.Start(ctx, config)
	if err != nil {
		fmt.Printf("Error starting workers: %v\n", err)
		return
	}

	fmt.Println("All tasks completed. Exiting...")

}
