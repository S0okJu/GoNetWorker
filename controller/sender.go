package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Signal struct {
	RequestId string `json:"request_id"`
	Status    string `json:"status"`
	Jobs      Jobs   `json:"jobs"`
}

func NewSignal(rid string, status string, jobs Jobs) *Signal {
	return &Signal{
		RequestId: rid,
		Status:    status,
		Jobs:      jobs,
	}
}

type Sender struct {
	jobs Jobs
}

func NewSender(jobs Jobs) *Sender {
	return &Sender{
		jobs: jobs,
	}
}

func (s *Sender) Send() error {
	// create requestid
	g := NewGenerator()
	rid := g.generate()

	// Signal
	signal := NewSignal(rid, "start", s.jobs)
	jsonData, err := json.Marshal(signal)
	if err != nil {
		return err
	}

	// Send to websocket server
	req, err := http.NewRequest("POST", "http://localhost:8080", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing response body")
		}
	}(resp.Body)

	// Debugging
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	fmt.Println("Response Body:", result)

	return nil
}
