package main

import (
	"fmt"
	"github.com/s0okjug/gonetworker/core"
)

func main() {

	// Open the JSON file
	cfg, err := core.LoadJSON("../test.json")
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		return
	}

	err = core.Execute(*cfg)
	if err != nil {
		fmt.Println("Error executing configuration:", err)
		return
	}

}
