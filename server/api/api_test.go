package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestSendConfigHandler(t *testing.T) {
	_ = os.Remove("../config.json")

	reqBody := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
		"key3": map[string]interface{}{
			"nestedKey": "nestedValue",
		},
	}
	reqBodyJSON, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", "/post", bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SendConfigHandler)

	handler.ServeHTTP(rr, req)

	// Check 200 OK status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Read file and compare with
	fileData, err := os.ReadFile("../config.json")
	if err != nil {
		t.Fatalf("Could not read config.json: %v", err)
	}

	// Compare file data to request body
	reqJson, _ := json.MarshalIndent(reqBody, "", " ")
	if string(fileData) != string(reqJson) {
		t.Errorf("Config file contains unexpected data: got %v want %v", string(fileData), string(reqJson))
	}
}
