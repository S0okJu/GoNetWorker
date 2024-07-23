package main

import (
	"gonetworker/api"
	"log"
	"net/http"
)

func main() {
	//core.Send("./example/docs/gonetwork.config.json")

	// Api
	http.HandleFunc("/api/json", api.SendConfigHandler)
	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))

}
