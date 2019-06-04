package taskqueue

import (
	"errors"
	"sync"
	"time"
)

type task struct {
	f          func(...interface{}) error
	retryTimes int
	args       []interface{}
}

// TaskQueue is taskqueue struct
type TaskQueue struct {
	interval time.Duration
	sync.RWMutex
	tasks     []task
	breakFlag bool
}

// New Create new struct
func New(interval time.Duration) *TaskQueue {
	return &TaskQueue{
		interval: interval,
	}
}

// Add is add job to task queue
func (t *TaskQueue) Add(f func(...interface{}) error, retryTimes int, args ...interface{}) error {
	if t.breakFlag {
		return errors.New("taskQueue.Stop is called")
	}
	t.Lock()
	t.tasks = append(t.tasks, task{
		f:          f,
		retryTimes: retryTimes,
		args:       args,
	})
	t.Unlock()
	return nil
}

func (t *TaskQueue) addNotCheckBreakFlag(f func(...interface{}) error, retryTimes int, args ...interface{}) {
	t.Lock()
	t.tasks = append(t.tasks, task{
		f:          f,
		retryTimes: retryTimes,
		args:       args,
	})
	t.Unlock()
}

// Run run taskqueue
func (t *TaskQueue) Run() {
L:
	for {
		if len(t.tasks) > 0 {
			tt := t.pop()
			if err := tt.f(tt.args); err != nil {
				t.retry(tt)
			}
		}
		if len(t.tasks) <= 0 && t.breakFlag {
			break L
		}
		time.Sleep(t.interval)
	}
}

func (t *TaskQueue) retry(tt task) {
	if tt.retryTimes > 1 {
		t.addNotCheckBreakFlag(tt.f, tt.retryTimes-1, tt.args...)
	}
}

func (t *TaskQueue) pop() task {
	t.Lock()
	task := t.tasks[0]
	t.tasks = t.tasks[1:]
	t.Unlock()
	return task
}

// Stop stop taskqueue
func (t *TaskQueue) Stop() {
	t.breakFlag = true
	for {
		if len(t.tasks) <= 0 {
			break
		}
	}
}
