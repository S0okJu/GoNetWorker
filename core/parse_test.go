package core

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		want   Jobs
		isErr  bool
	}{
		{
			name: "Single work, multiple tasks",
			config: Config{
				Settings: Settings{SleepRange: 5},
				Works: []Work{
					{
						Uri:  "http://localhost",
						Port: 8000,
						Tasks: []Task{
							{
								Path:   "/users/1",
								Method: "GET",
								Body:   nil,
							},
							{
								Path:   "/users/2",
								Method: "POST",
								Body:   map[string]string{"key2": "value2"},
							},
						},
					},
				},
			},
			want: Jobs{
				{
					Url:    "http://localhost:8000/users/1",
					Method: "GET",
					Body:   nil,
				},
				{
					Url:    "http://localhost:8000/users/2",
					Method: "POST",
					Body:   map[string]string{"key2": "value2"},
				},
			},
			isErr: false,
		},
		{
			name: "No tasks available",
			config: Config{
				Settings: Settings{SleepRange: 5},
				Works:    []Work{},
			},
			want:  nil,
			isErr: true,
		},
		{
			name: "Task with no port",
			config: Config{
				Settings: Settings{SleepRange: 5},
				Works: []Work{
					{
						Uri:  "http://localhost",
						Port: 0,
						Tasks: []Task{
							{
								Path:   "/users/1",
								Method: "GET",
								Body:   nil,
							},
						},
					},
				},
			},
			want:  nil,
			isErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.config)
			got, err := p.Parse()

			if (err != nil) != tt.isErr {
				t.Errorf("Parse() error = %v, isErr %v", err, tt.isErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestFixedUrl tests the FixedUrl function
func TestFixedUrl(t *testing.T) {
	// 테스트 케이스 1: 정상적인 범위 입력
	input := "/users/{[1-5]}"
	result, err := FixedUrl(input)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// 결과 출력
	fmt.Println("Test Case 1 - Result:", result)

	// 정규 표현식으로 결과 확인 (1-5 사이의 숫자가 대체되었는지 확인)
	expectedPattern := regexp.MustCompile(`/users/\d+`)
	if !expectedPattern.MatchString(result) {
		t.Fatalf("Expected result to match /users/\\d+, but got %s", result)
	}

	// 숫자가 범위 내에 있는지 확인
	re := regexp.MustCompile(`/users/(\d+)`)
	matches := re.FindStringSubmatch(result)
	if len(matches) != 2 {
		t.Fatalf("Expected one number in the result, but got %v", matches)
	}

	number, err := strconv.Atoi(matches[1])
	if err != nil {
		t.Fatalf("Error converting matched number to integer: %v", err)
	}
	if number < 1 || number > 5 {
		t.Fatalf("Expected number between 1 and 5, but got %d", number)
	}

	// 테스트 케이스 2: 잘못된 범위 입력
	input = "/users/{[5-]}"
	_, err = FixedUrl(input)
	if err == nil {
		t.Fatalf("Expected error for invalid input, but got none")
	} else {
		fmt.Println("Test Case 2 - Error:", err)
	}

	// 테스트 케이스 3: 범위가 없는 입력
	input = "/users/"
	result, err = FixedUrl(input)
	if err == nil {
		t.Fatalf("Expected error for missing range input, but got none")
	} else {
		fmt.Println("Test Case 3 - Error:", err)
	}
}
