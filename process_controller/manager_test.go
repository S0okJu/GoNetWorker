package processcontroller

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"
)

func TestWorkerProcessingAndStopping(t *testing.T) {
	// Create a context with cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure cancel is called to release resources

	// Create a pipe to capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the worker function in a goroutine
	go test(ctx)

	// Allow the worker to run for 3 seconds
	time.Sleep(3 * time.Second)

	// Stop the worker by canceling the context
	cancel()

	// Give some time for the worker to stop
	time.Sleep(1 * time.Second)

	// Restore original stdout
	w.Close()
	os.Stdout = oldStdout

	// Capture output from the worker
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Check if the worker processed jobs multiple times
	if !containsMultipleOccurrences(output, "Worker processing job", 2) {
		t.Errorf("expected 'Worker processing job' message multiple times, but got:\n%s", output)
	}

	// Check if the worker stopped
	if !containsOnce(output, "Worker stopping") {
		t.Errorf("expected 'Worker stopping' message once, but got:\n %s", output)
	}
}

// Helper function to check if a substring occurs at least 'count' times in a string
func containsMultipleOccurrences(s, substr string, count int) bool {
	n := 0
	for n < count {
		index := bytes.Index([]byte(s), []byte(substr))
		if index == -1 {
			return false
		}
		s = s[index+len(substr):]
		n++
	}
	return true
}

// Helper function to check if a substring occurs exactly once in a string
func containsOnce(s, substr string) bool {
	return bytes.Count([]byte(s), []byte(substr)) == 1
}
