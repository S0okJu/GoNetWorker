package controller

import "sync"

type worker struct {
	channel chan bool
}

func newWorker() worker {
	return worker{channel: make(chan bool)}
}

type Workers []worker

func NewWorkers(cnt int) (*Workers, error) {
	var workers Workers
	for i := 0; i < cnt; i++ {
		workers = append(workers, newWorker())
	}
	return &workers, nil
}

func (ws *Workers) Assign(task func()) {
	var wg sync.WaitGroup

	for _, w := range *ws {
		wg.Add(1)

		go func(wk worker) {
			defer wg.Done()

			task()

			wk.channel <- true
		}(w)
	}

	wg.Wait()
}
func (ws *Workers) Done() {
	for _, worker := range *ws {
		<-worker.channel
	}
}
