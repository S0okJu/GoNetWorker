package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func RequestTo() {
	// Read the config file
	data, err := ioutil.ReadFile("docs/config.json")
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		return
	}

	// Parse the JSON
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Process each work item
	for _, work := range config.Works {
		baseURI := work.Uri
		if work.Port != 0 {
			baseURI = fmt.Sprintf("%s:%d", baseURI, work.Port)
		}

		// Process each request
		for _, req := range work.Request {
			fullURI := baseURI + req.Path
			fmt.Printf("Making %s request to %s\n", req.Method, fullURI)

			// Prepare the request
			var httpReq *http.Request
			var err error

			switch req.Method {
			case "GET":
				params := url.Values{}
				for key, dataType := range req.Param {
					params.Add(key, generateRandomValue(dataType))
				}
				httpReq, err = http.NewRequest(req.Method, fullURI+"?"+params.Encode(), nil)
			case "POST":
				formData := url.Values{}
				for key, dataType := range req.Param {
					formData.Add(key, generateRandomValue(dataType))
				}
				httpReq, err = http.NewRequest(req.Method, fullURI, strings.NewReader(formData.Encode()))
				httpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			default:
				fmt.Printf("Unsupported method: %s\n", req.Method)
				continue
			}

			if err != nil {
				fmt.Printf("Error creating request: %v\n", err)
				continue
			}

			// Send the request
			client := &http.Client{}
			resp, err := client.Do(httpReq)
			if err != nil {
				fmt.Printf("Error sending request: %v\n", err)
				continue
			}
			defer resp.Body.Close()

			// Print the response
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Printf("Response status: %s\n", resp.Status)
			fmt.Printf("Response body: %s\n\n", string(body))
		}
	}
}

func generateRandomValue(dataType string) string {
	switch dataType {
	case "string":
		return randomString(10)
	case "int":
		return fmt.Sprintf("%d", rand.Intn(1000))
	case "float":
		return fmt.Sprintf("%.2f", rand.Float64()*100)
	case "bool":
		return fmt.Sprintf("%t", rand.Intn(2) == 1)
	default:
		return "unknown_type"
	}
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
