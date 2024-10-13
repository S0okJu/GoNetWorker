package core

import (
	"fmt"
	"regexp"
	"testing"
)

func TestUrlPatternSelect(t *testing.T) {
	// 테스트 케이스 1: Dynamic URL (with braces)
	input := "/users/{[1-5]}"
	result, err := urlPatternSelect(input)
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

	// 테스트 케이스 2: Static URL (no braces)
	input = "/static/url/path"
	result, err = urlPatternSelect(input)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Static URL은 변환되지 않고 그대로 반환
	if result != input {
		t.Fatalf("Expected %v, got %v", input, result)
	}
	fmt.Println("Test Case 2 - Result:", result)

	// 테스트 케이스 3: Invalid dynamic URL (wrong format)
	input = "/users/{[1-]}"
	_, err = urlPatternSelect(input)
	if err == nil {
		t.Fatalf("Expected an error for invalid input, but got none")
	} else {
		fmt.Println("Test Case 3 - Error:", err)
	}
}
