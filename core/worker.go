package core

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Worker  작업을 수행하는 구조체입니다.
type Worker struct{}

// NewWorker  새로운 Workers 인스턴스를 생성합니다.
func NewWorker() (*Worker, error) {
	return &Worker{}, nil
}

// Start  작업을 시작하고 취소 신호를 감지합니다.
func (ws *Worker) Start(ctx context.Context, cfg *Config) error {

	parser := NewParser(*cfg)
	jobs, err := parser.Parse()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	// 컨텍스트가 취소되지 않은 동안 실행
	for ctx.Err() == nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := request(ctx, jobs[rand.Intn(len(jobs))], cfg.Settings.SleepRange)
			// 에러 발생 시 에러 시그널 채널에 전송
			if err != nil {
				select {
				case errChan <- err:
				default:
					fmt.Println("Error channel is full:", err)
				}
			}
		}()

		select {
		case <-ctx.Done():
			fmt.Println("Context canceled during sleep. Waiting for workers to finish...")
		case <-time.After(1 * time.Second):
		}
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return fmt.Errorf("one or more errors occurred")
	}

	return nil
}

// request HTTP 요청을 수행
func request(ctx context.Context, job Job, sleepRange int) error {
	// 전체 URL 생성
	url := job.Url

	sleepDuration := time.Duration(rand.Intn(sleepRange)) * time.Second
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(sleepDuration):
	}

	var req *http.Request
	var err error
	client := &http.Client{}

	switch strings.ToUpper(job.Method) {
	case "GET":
		req, err = http.NewRequestWithContext(ctx, "GET", url, nil)
		fmt.Println("GET request to", url)
	case "POST":
		body, _ := job.ConvertTo()
		req, err = http.NewRequestWithContext(ctx, "POST", url, body)
		req.Header.Set("Content-Type", "application/json")
		fmt.Println("POST request to", url)
	default:
		return fmt.Errorf("unsupported HTTP method: %s", job.Method)
	}

	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			fmt.Println("요청이 취소되었습니다:", ctx.Err())
			return ctx.Err()
		}
		fmt.Println("요청 중 오류:", err)
		return fmt.Errorf("요청 중 오류: %v", err)
	}
	defer resp.Body.Close()

	return nil
}
