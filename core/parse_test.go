package core

import (
	"encoding/json"
	"math/rand"
	"testing"
	"time"
)

// MockRandomSelector for consistent results in tests
type MockRandomSelector struct {
	value string
}

func NewMockRandomSelector(value string) *MockRandomSelector {
	return &MockRandomSelector{value: value}
}

func (m *MockRandomSelector) Select() (string, error) {
	return m.value, nil
}

// TestConvertTo tests the ConvertTo function of Job struct
func TestConvertTo(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	// Setup job with body containing "string" and "int" types
	job := &Job{
		Body: map[string]string{
			"name": "string",
			"age":  "int",
		},
	}

	// Run ConvertTo and decode JSON
	bodyReader, err := job.ConvertTo()
	if err != nil {
		t.Fatalf("ConvertTo returned an error: %v", err)
	}

	var newBody map[string]interface{}
	if err := json.NewDecoder(bodyReader).Decode(&newBody); err != nil {
		t.Fatalf("error decoding JSON: %v", err)
	}

	if age, ok := newBody["age"].(float64); !ok || age < 0 || age > 9 {
		t.Errorf("expected 'age' to be an integer between 0 and 9, got '%v'", newBody["age"])
	}
}

// TestParse tests the Parse function of Parser struct
func TestParse(t *testing.T) {
	// Define config
	config := Config{
		Settings: Settings{
			SleepRange: 5,
			CcuMax:     10,
		},
		Works: []Work{
			{
				Uri:  "http://localhost",
				Port: 8080,
				Tasks: []Task{
					{
						Path:   "/test",
						Method: "POST",
						Query:  map[string]string{"q": "query"},
						Body:   map[string]string{"name": "string"},
					},
				},
			},
		},
	}

	// Parse config into jobs
	parser := NewParser(config)
	jobs, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse returned an error: %v", err)
	}

	// Validate job count
	if len(jobs) != 1 {
		t.Errorf("expected 1 job, got %d", len(jobs))
	}

	// Check job details
	job := jobs[0]
	expectedUrl := "http://localhost:8080/test"
	if job.Url != expectedUrl {
		t.Errorf("expected URL to be '%s', got '%s'", expectedUrl, job.Url)
	}
	if job.Method != "POST" {
		t.Errorf("expected method to be 'POST', got '%s'", job.Method)
	}
	if job.Body["name"] != "string" {
		t.Errorf("expected body 'name' to be 'string', got '%s'", job.Body["name"])
	}
}

// TestHasBrace tests HasBrace function
func TestHasBrace(t *testing.T) {
	tests := []struct {
		url      string
		expected bool
	}{
		{"http://localhost:8080/users/{id}", true},
		{"http://localhost:8080/products", false},
		{"http://example.com/api/{user}/{profile}", true},
		{"https://api.example.com/test", false},
	}

	for _, test := range tests {
		result := HasBrace(test.url)
		if result != test.expected {
			t.Errorf("HasBrace(%s) = %v; expected %v", test.url, result, test.expected)
		}
	}
}
