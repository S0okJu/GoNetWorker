package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"regexp"
	"strings"
)

// Config 설정 파일 구조체
// endpoint.json 파일을 읽어들이기 위한 구조체
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
			us := NewUrlSelector(work.Uri, work.Port, task)
			fullUrl, err := us.Select()
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

// ConvertTo request POST Body 변환을 위한 함수
// value는 타입을 나타내며, 타입에 맞게 랜덤 값으로 변환
func (j *Job) ConvertTo() (io.Reader, error) {
	newBody := make(map[string]interface{})
	// Change to random value
	re := regexp.MustCompile(`\{\[(\d+)-(\d+)\]\}`)
	for key, value := range j.Body {
		switch {
		case value == "string":
			selector := NewRandomSelector(5)
			randStr, serr := selector.Select()
			if serr != nil {
				return nil, fmt.Errorf("error selecting random value: %v", serr)
			}

			newBody[key] = randStr
		case value == "int":
			randomInt := rand.Intn(10)
			newBody[key] = randomInt
		case re.MatchString(value):
			minR, maxR, err := ConvertMinMax(value)
			if err != nil {
				return nil, fmt.Errorf("error converting min-max: %v", err)
			}
			randInt := rand.Intn(maxR-minR) + minR
			newBody[key] = randInt
		default:
			return nil, fmt.Errorf("unsupported type: %s", value)
		}
	}
	// Marshal
	jsonData, err := json.Marshal(newBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling to JSON: %v", err)
	}

	return bytes.NewBuffer(jsonData), nil
}

// HasBrace
// Dynamic url을 가지고 있는지 {}를 확인
func HasBrace(urlStr string) bool {
	// Decode the URL first
	decodedUrl, err := url.QueryUnescape(urlStr)
	if err != nil {
		fmt.Println("Error decoding URL:", err)
		return false
	}

	// Check if the decoded URL contains { and }
	return strings.Contains(decodedUrl, "{") && strings.Contains(decodedUrl, "}")
}
