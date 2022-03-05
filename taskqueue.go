package taskqueue

import (
	"fmt"
	"sync"
	"time"
)

var (
	ErrStopped = fmt.Errorf("taskQueue.Stop is called")
)

const (
	defaultTimeout = 3
)

// Task is Task
type Task struct {
	f         func(...any) error
	args      []any
	retrytime int
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

type Options struct {
	retryTime int
}

type Option func(*Options)

func RetryTime(retryTime int) Option {
	return func(args *Options) {
		args.retryTime = retryTime
	}
}

// Add is add job to task queue
func (t *TaskQueue) Add(f func(...any) error, args ...any) error {
	// opts := &Options {
	//     retryTime: defaultTimeout,
	// }

	// for _, optionFunc := range optionFuncs {
	//     optionFunc(opts)
	// }

	if t.breakFlag {
		return ErrStopped
	}

	t.Lock()
	defer t.Unlock()
	t.tasks = append(t.tasks, Task{
		f:         f,
		args:      args,
		retrytime: defaultTimeout,
		// retrytime: opts.retryTime,
	})

	return nil
}

func (t *TaskQueue) addNotCheckBreakFlag(task Task) {
	t.Lock()
	defer t.Unlock()
	t.tasks = append(t.tasks, task)
}

// Run run taskqueue
func (t *TaskQueue) Run() {
L:
	for {
		if len(t.tasks) > 0 {
			tt := t.pop()
			if err := tt.f(tt.args); err != nil {
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
	tt.retrytime = tt.retrytime - 1
	if tt.retrytime >= 1 {
		t.addNotCheckBreakFlag(tt)
	}

	t.Error <- Err{
		Task: tt,
		Err:  err,
	}
}

func (t *TaskQueue) pop() Task {
	t.Lock()
	defer t.Unlock()
	task := t.tasks[0]
	t.tasks = t.tasks[1:]
	return task
}

// Stop stop taskqueue
func (t *TaskQueue) Stop() {
	t.breakFlag = true
	<-t.closeCh
	close(t.closeCh)
	close(t.Error)
}

func (t *TaskQueue) Count() int {
	t.Lock()
	defer t.Unlock()
	return len(t.tasks)
}
