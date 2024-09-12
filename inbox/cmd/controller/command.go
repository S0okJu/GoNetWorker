package controller

import (
	"fmt"
	"github.com/d7mekz/gonetworker/controller"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func Start() {
	mux := http.NewServeMux()

	mux.HandleFunc("/data", controller.GetConfigHandler)

	// middleware
	handler := cors.Default().Handler(mux)

	fmt.Println("Test server is listening on port 8000...")
	if err := http.ListenAndServe(":8000", handler); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
