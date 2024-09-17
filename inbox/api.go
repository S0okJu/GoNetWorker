package inbox

import (
	"encoding/json"
	"net/http"
)

func GetConfigHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var config Config
	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	parser := NewParser(config)
	jobs, err := parser.Parse()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sender := NewSender(jobs)
	err = sender.Send()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
