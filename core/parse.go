package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
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

// FixedUrl
// Dynamic url 범위를 하나로 고정한다.
func FixedUrl(urlStr string) (string, error) {
	// 정규식으로 파싱
	decodedUrl, err := url.QueryUnescape(urlStr)
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`\{\[(\d+)-(\d+)\]\}`)
	matches := re.FindStringSubmatch(decodedUrl)
	if len(matches) != 3 {
		return "", fmt.Errorf("Invalid Regex Format")
	}

	// 숫자로 반환
	minR, err1 := strconv.Atoi(matches[1])
	maxR, err2 := strconv.Atoi(matches[2])
	if err1 != nil || err2 != nil {
		return "", fmt.Errorf("Error converting to integer")
	}

	// 랜덤 시드 생성
	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(maxR-minR+1) + minR

	// 새로운 url
	res := re.ReplaceAllString(decodedUrl, fmt.Sprintf("%d", idx))

	return res, nil

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
