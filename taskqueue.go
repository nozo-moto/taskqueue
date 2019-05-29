package taskqueue

import (
	"sync"
	"time"
)

type task struct {
	f          func() error
	retryTimes int
}

// TaskQueue is taskqueue struct
type TaskQueue struct {
	interval time.Duration
	sync.RWMutex
	tasks []task
}

// New Create new struct
func New(interval time.Duration) *TaskQueue {
	return &TaskQueue{
		interval: interval,
	}
}

// Add is add job to task queue
func (t *TaskQueue) Add(f func() error, retryTimes int) {
	t.Lock()
	t.tasks = append(t.tasks, task{
		f:          f,
		retryTimes: retryTimes,
	})
	t.Unlock()
}

// Run run taskqueue
func (t *TaskQueue) Run() {
	for {
		if len(t.tasks) > 0 {
			tt := t.pop()
			if err := tt.f(); err != nil {
				t.retry(tt)
			}
		}
		time.Sleep(t.interval)
	}
}

func (t *TaskQueue) retry(tt task) {
	if tt.retryTimes > 1 {
		go t.Add(tt.f, tt.retryTimes-1)
	}
}

func (t *TaskQueue) pop() task {
	t.Lock()
	task := t.tasks[0]
	t.tasks = t.tasks[1:]
	t.Unlock()
	return task
}
