package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
)

// -- Json Formats ---
type Config struct {
	Settings Settings `json:"settings,omitempty"`
	Works    []Work   `json:"works,omitempty"`
}

func (c *Config) GetSleepRange() int {
	return c.Settings.SleepRange
}

func (c *Config) GetCcuMax() int {
	return c.Settings.CcuMax
}

type Settings struct {
	SleepRange int `json:"sleep_range"`
	CcuMax     int `json:"ccu_max"`
}

type Work struct {
	Uri   string `json:"uri,omitempty"`
	Port  int    `json:"port,omitempty"`
	Tasks []Task `json:"tasks,omitempty"`
}

type Task struct {
	Path   string            `json:"path,omitempty"`
	Method string            `json:"method,omitempty"`
	Query  map[string]string `json:"query,omitempty"`
	Body   map[string]string `json:"body"`
}

// Job request를 위한 구조체
type Job struct {
	// RequestID string            `json:"request_id"`
	Url    string            `json:"url"`
	Method string            `json:"method"`
	Body   map[string]string `json:"body"`
}

// Jobs Job 배열
type Jobs []Job

// Parser 파일로부터 읽은 데이터를 파싱
type Parser struct {
	config Config
}

func NewParser(cfg Config) *Parser {
	return &Parser{
		config: cfg,
	}
}

// Parse 파일로부터 읽은 데이터를 Jobs 형태로 변환
func (p *Parser) Parse() (Jobs, error) {
	validator := NewValidator()

	var jobs Jobs
	if len(p.config.Works) == 0 {
		return nil, fmt.Errorf("no works found")
	}

	for _, work := range p.config.Works {
		validator.Port(work.Port)
		if validator.IsError() == true {
			return nil, fmt.Errorf("Invalid port number")
		}

		for _, task := range work.Tasks {
			fullUrl, err := getUrl(work.Uri, work.Port, task)
			if err != nil {
				break
			}
			jobs = append(jobs, Job{
				Url:    fullUrl,
				Method: task.Method,
				Body:   task.Body,
			})
		}
	}

	if len(jobs) == 0 {
		return nil, fmt.Errorf("no jobs found")
	}
	return jobs, nil
}

// getUrl URL 생성
func getUrl(uri string, port int, task Task) (string, error) {
	var baseURI string
	if port != 0 {
		baseURI = fmt.Sprintf("%s:%d", uri, port)
	} else {
		return "", fmt.Errorf("Invalid port number")
	}

	baseURI = baseURI + task.Path

	// Get Query Parameters
	if task.Method == "GET" {
		params := url.Values{}
		if len(task.Query) > 0 {
			for key, value := range task.Query {
				params.Add(key, value)
			}
		}

		u, err := url.Parse(baseURI)
		if err != nil {
			return "", err
		}

		u.RawQuery = params.Encode()
		finalUrl := u.String()
		return finalUrl, nil
	} else {
		return baseURI, nil
	}
}

// ConvertTo request POST Body 변환을 위한 함수
func (j *Job) ConvertTo() (io.Reader, error) {
	jsonData, err := json.Marshal(j.Body)
	if err != nil {
		return nil, fmt.Errorf("error marshaling to JSON: %v", err)
	}

	return bytes.NewBuffer(jsonData), nil
}
