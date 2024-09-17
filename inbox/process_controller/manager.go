package processcontroller

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/s0okjug/gonetworker/core"
	"time"
)

const (
	Unknown string = "Unknown"
	Start   string = "Start"
	Pending string = "Pending"
	Running string = "Running"
	Stop    string = "Stop"
	Failed  string = "Failed"
)

type state string

func NewState(st string) state {
	return state(st)
}

func (s *state) String() string {
	return s.String()
}

func (s *state) Convert(st string) {
	*s = state(st)
}

func (s *state) IsStart() bool {
	return s.String() == Start
}

func (s *state) IsPending() bool {
	return s.String() == Pending
}

func (s *state) IsInRunning() bool {
	return s.String() == Running
}

func (s *state) IsFailed() bool {
	return s.String() == Failed
}
func (s *state) IsStop() bool {
	return s.String() == Stop

}

// Process is process object to execute worker
// id : process id(uuid)
// status : Process execution status
type Process struct {
	id     string
	state  state
	job    core.Job
	cancel context.CancelFunc
}

// NewProcess Create new Process Object.
func NewProcess(job core.Job) *Process {
	return &Process{
		id:    uuid.NewString(),
		job:   job,
		state: NewState(Unknown),
	}
}

// Convert is convert core.Jobs to Process
func Convert(jobs core.Jobs) *[]Process {
	processes := make([]Process, len(jobs))
	for i, job := range jobs {
		processes[i] = *NewProcess(job)
	}
	return &processes
}

func (p *Process) Start(result chan<- Process) {
	fmt.Printf("Process Start : %s", p.id)
	p.state.Convert(Running)

	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel
	// Start processing the job in a goroutine
	go func() {
		select {
		case <-ctx.Done(): // If context is canceled, stop the process
			fmt.Printf("Process %s canceled\n", p.id)
			p.state = NewState("Stop")
			result <- *p
			return
		default:
			// Simulate job execution by sleeping for 2 seconds
			fmt.Printf("Process %s working on job\n", p.id)
			time.Sleep(2 * time.Second)

			// Mark process as completed
			fmt.Printf("Process %s completed job\n", p.id)
			p.state = NewState("Stop")
			result <- *p
		}
	}()
}

func (p *Process) Stop() {
	if p.state.IsInRunning() {
		p.cancel()
		p.state.Convert(Stop)
		fmt.Printf("Process Stop : %s", p.id)
	} else {
		fmt.Printf("Process Not Running : %s", p.id)
	}
}

func request(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shut down now")
			return
		default:
			// request http

		}

	}
}

func test(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Worker stopping")
			return
		default:
			fmt.Println("Worker processing job")
			time.Sleep(1 * time.Second)
		}
	}
}
