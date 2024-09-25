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

// Workers는 작업을 수행하는 구조체입니다.
type Workers struct{}

// NewWorkers는 새로운 Workers 인스턴스를 생성합니다.
func NewWorkers(cnt int) (*Workers, error) {
	return &Workers{}, nil
}

// Start는 작업을 시작하고 취소 신호를 감지합니다.
func (ws *Workers) Start(ctx context.Context, cfg *Config) error {

	// 작업 파싱
	parser := NewParser(*cfg)
	jobs, err := parser.Parse()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 1) // 에러를 수집할 채널

	// 컨텍스트가 취소되지 않은 동안 루프 실행
	for ctx.Err() == nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// 에러 변수를 고루틴 내부에서 선언하여 데이터 레이스 방지
			err := request(ctx, jobs[rand.Intn(len(jobs))], cfg.Settings.SleepRange)
			if err != nil {
				// 에러 채널에 에러 전송 (버퍼 크기 주의)
				select {
				case errChan <- err:
				default:
					// 에러 채널이 가득 찬 경우 에러 로그 출력
					fmt.Println("Error channel is full:", err)
				}
			}
		}()

		// 다음 작업 전 대기 시간 (컨텍스트 취소 감지)
		select {
		case <-ctx.Done():
			fmt.Println("Context canceled during sleep. Waiting for workers to finish...")
			// 루프가 컨텍스트 취소로 인해 종료되므로 추가 조치 불필요
		case <-time.After(1 * time.Second):
			// 1초 대기 후 다음 반복 실행
		}
	}

	// 모든 고루틴이 완료될 때까지 대기
	wg.Wait()
	close(errChan) // 에러 채널 닫기

	// 에러 처리
	if len(errChan) > 0 {
		for err := range errChan {
			fmt.Println("Error occurred:", err)
		}
		return fmt.Errorf("one or more errors occurred")
	}

	return nil
}

// request는 HTTP 요청을 수행합니다.
func request(ctx context.Context, job Job, sleepRange int) error {
	// 전체 URL 생성
	url := job.Url

	// 랜덤 시간만큼 대기 (컨텍스트 취소 감지)
	sleepDuration := time.Duration(rand.Intn(sleepRange)) * time.Second
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(sleepDuration):
		// 지정된 시간만큼 대기
	}

	// HTTP 요청 생성
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
		// 컨텍스트 취소로 인한 오류 처리
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
