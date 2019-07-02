package taskqueue

import (
	"errors"
	"sync"
	"time"
)

// Task is Task
type Task struct {
	F          func() error
	RetryTimes int
}

// Err is tasks fail error
type Err struct {
	Task Task
	Err  error
}

// TaskQueue is taskqueue struct
type TaskQueue struct {
	interval time.Duration
	sync.RWMutex
	tasks     []Task
	breakFlag bool
	closeCh   chan int
	Error     chan Err
}

// New Create new struct
func New(interval time.Duration) *TaskQueue {
	return &TaskQueue{
		interval: interval,
		closeCh:  make(chan int),
		Error:    make(chan Err),
	}
}

// Add is add job to task queue
func (t *TaskQueue) Add(f func() error, retryTimes int) error {
	if t.breakFlag {
		return errors.New("taskQueue.Stop is called")
	}
	t.Lock()
	t.tasks = append(t.tasks, Task{
		F:          f,
		RetryTimes: retryTimes,
	})
	t.Unlock()
	return nil
}

func (t *TaskQueue) addNotCheckBreakFlag(f func() error, retryTimes int) {
	t.Lock()
	t.tasks = append(t.tasks, Task{
		F:          f,
		RetryTimes: retryTimes,
	})
	t.Unlock()
}

// Run run taskqueue
func (t *TaskQueue) Run() {
L:
	for {
		if len(t.tasks) > 0 {
			tt := t.pop()
			if err := tt.F(); err != nil {
				t.retry(tt, err)
			}
		}
		if len(t.tasks) <= 0 && t.breakFlag {
			t.closeCh <- 1
			break L
		}
		time.Sleep(t.interval)
	}
}

func (t *TaskQueue) retry(tt Task, err error) {
	tt.RetryTimes = tt.RetryTimes - 1
	if tt.RetryTimes >= 1 {
		t.addNotCheckBreakFlag(tt.F, tt.RetryTimes)
	}

	t.Error <- Err{
		Task: tt,
		Err:  err,
	}
}

func (t *TaskQueue) pop() Task {
	t.Lock()
	task := t.tasks[0]
	t.tasks = t.tasks[1:]
	t.Unlock()
	return task
}

// Stop stop taskqueue
func (t *TaskQueue) Stop() {
	t.breakFlag = true
	<-t.closeCh
	close(t.closeCh)
	close(t.Error)
}
