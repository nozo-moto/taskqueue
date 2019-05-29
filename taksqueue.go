package taskqueue

import (
	"sync"
	"time"
)

// TaskQueue is taskqueue struct
type TaskQueue struct {
	interval time.Duration
	sync.RWMutex
	tasks []func() error
}

// New Create new struct
func New(interval time.Duration) *TaskQueue {
	return &TaskQueue{
		interval: interval,
	}
}

// Add is add job to task queue
func (t *TaskQueue) Add(f func() error) {
	t.Lock()
	t.tasks = append(t.tasks, f)
	t.Unlock()
}

// Run run taskqueue
func (t *TaskQueue) Run() {
	go func() {
		var err error

		for {
			if len(t.tasks) > 0 {
				t.Lock()
				task := t.tasks[0]
				t.tasks = t.tasks[1:]
				t.Unlock()
				if err = task(); err != nil {
					go t.Add(task)
				}
			}
			time.Sleep(t.interval)
		}
	}()
}
