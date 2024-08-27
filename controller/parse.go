package controller

import "fmt"

// -- Json Formats ---
type Config struct {
	Settings Settings `json:"settings,omitempty"`
	Works    []Work   `json:"works,omitempty"`
}

type Settings struct {
	SleepRange int `json:"sleep_range"`
}

type Work struct {
	Uri   string `json:"uri,omitempty"`
	Port  int    `json:"port,omitempty"`
	Tasks []Task `json:"tasks,omitempty"`
}

type Task struct {
	Path   string            `json:"path,omitempty"`
	Method string            `json:"method,omitempty"`
	Body   map[string]string `json:"body"`
}

// Job is the task to be executed
type Job struct {
	RequestID string            `json:"request_id"`
	Url       string            `json:"url"`
	Method    string            `json:"method"`
	Body      map[string]string `json:"body"`
}

type Jobs []Job

type Parser struct {
	config Config
}

func NewParser(cfg Config) *Parser {
	return &Parser{
		config: cfg,
	}
}

// Parse parses the configuration and returns a list of tasks
func (p *Parser) Parse() (Jobs, error) {
	validator := NewValidator()

	var jobs Jobs
	if len(p.config.Works) == 0 {
		return nil, fmt.Errorf("no works found")
	}

	for _, task := range p.config.Works {
		// Validate the Port
		validator.Port(task.Port)
		if validator.IsError() == true {
			return nil, fmt.Errorf("Invalid port number")
		}

		for _, t := range task.Tasks {
			fullUrl := url(task.Uri, task.Port, t.Path)
			jobs = append(jobs, Job{
				Url:    fullUrl,
				Method: t.Method,
				Body:   t.Body,
			})
		}
	}

	if len(jobs) == 0 {
		return nil, fmt.Errorf("no jobs found")
	}
	return jobs, nil
}

// url returns the full URL for the task
func url(uri string, port int, path string) string {
	baseURI := uri
	if port != 0 {
		baseURI = fmt.Sprintf("%s:%d", baseURI, port)
	}
	return baseURI + path
}
