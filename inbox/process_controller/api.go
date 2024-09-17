package processcontroller

import (
	"encoding/json"
	"github.com/s0okjug/gonetworker/core"
	"net/http"
)

func StartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
		return
	}

	var sender core.Sender
	err := json.NewDecoder(r.Body).Decode(&sender)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Create new process
	processes := Convert(sender.Jobs)
	resultChan := make(chan Process)

	for _, process := range *processes {
		go process.Start(resultChan)
	}

	return
}
