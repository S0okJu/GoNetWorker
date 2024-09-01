package processcontroller

import (
	"context"
	"fmt"
	"github.com/google/uuid"
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
	cancel context.CancelFunc
}

// NewProcess Create new Process Object.
func NewProcess() *Process {
	return &Process{
		id:    uuid.NewString(),
		state: NewState(Unknown),
	}
}

func (p *Process) Start(result chan<- Process) {
	fmt.Println("Process Start : %s", p.id)
	p.state.Convert(Running)

	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel
	go test(ctx)

	result <- *p
}

func (p *Process) Stop() {
	if p.state.IsInRunning() {
		p.cancel()
		p.state.Convert(Stop)
		fmt.Println("Process Stop : %s", p.id)
	} else {
		fmt.Println("Process Not Running : %s", p.id)
	}
}

// test is a function that each worker goroutine will run
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
