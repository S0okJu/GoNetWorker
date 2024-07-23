package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

// SendConfig
// Method Post
func SendConfigHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	// Read body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error Reading request body", http.StatusInternalServerError)
		return
	}
	r.Body.Close()

	// Restore to file
	var jsonData map[string]interface{}
	if jerr := json.Unmarshal(body, &jsonData); jerr != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	resJson, err := json.MarshalIndent(jsonData, "", " ")
	if err != nil {
		http.Error(w, "Error creating response JSON", http.StatusInternalServerError)
		return
	}

	// write file
	if werr := os.WriteFile("../config.json", resJson, 0644); werr != nil {
		http.Error(w, "Error writing to file", http.StatusInternalServerError)
		return
	}

	// Success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, werr := w.Write(resJson)
	if werr != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		return
	}

	return
}
