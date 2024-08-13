package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Requests struct {
	requests []Request
}

func (r *Requests) Send(sleepRange int) {
	cnt := len(r.requests)
	client := &http.Client{Timeout: 10 * time.Second}

	for {
		num := rand.Intn(cnt)
		req := r.requests[num]

		httpReq, err := r.createHTTPRequest(req)
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			time.Sleep(1 * time.Second)
			continue
		}

		res, err := client.Do(httpReq)
		log.Println("Request sent to", req.Url, res.Status)
		if err != nil {
			fmt.Printf("Error sending request: %v\n", err)
			time.Sleep(1 * time.Second)
			continue
		}

		// Generate random sleep duration between 1 and 10 seconds
		sleepDuration := time.Duration(rand.Intn(sleepRange)) * time.Second
		fmt.Printf("Sleeping for %v before next request\n", sleepDuration)
		time.Sleep(sleepDuration)
	}
}

func (r *Requests) createHTTPRequest(req Request) (*http.Request, error) {
	var httpReq *http.Request
	var err error

	switch req.Method {
	case "GET":
		httpReq, err = http.NewRequest(req.Method, req.Url, nil)
	case "POST":
		//formData := url.Values{}
		//for key, dataType := range req.Body {
		//	formData.Add(key, generateRandomValue(dataType))
		//}
		//fmt.Println(formData.Encode())
		// Create a map to hold our JSON data
		jsonData := make(map[string]interface{})

		for key, dataType := range req.Body {
			jsonData[key] = generateRandomValue(dataType)
		}

		// Marshal the map into JSON
		jsonBody, jerr := json.Marshal(jsonData)
		if jerr != nil {
			return nil, fmt.Errorf("error marshaling JSON: %v", err)
		}

		fmt.Println("JSON body:", string(jsonBody))

		httpReq, err = http.NewRequest(req.Method, req.Url, bytes.NewBuffer(jsonBody))
		//httpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	default:
		return nil, fmt.Errorf("unsupported method: %s", req.Method)
	}

	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	return httpReq, nil
}

type Request struct {
	Url    string
	Method string
	Body   map[string]string
}

func (r *Requests) Add(req Request) {
	r.requests = append(r.requests, req)
	fmt.Println("Done")
}

func Send(path string) {
	// Parse Json

	p := newParser(path)
	p.parseJson()

	res := p.cfg
	var reqs Requests
	// Convert to requests
	for _, work := range res.Works {
		baseURI := Url(work.Uri, work.Port, "")
		for _, req := range work.Info {
			fullURI := Url(baseURI, 0, req.Path)
			reqs.Add(Request{
				Url:    fullURI,
				Method: req.Method,
				Body:   req.Param,
			})

		}
	}

	// send
	reqs.Send(p.cfg.Settings.SleepRange)
}
