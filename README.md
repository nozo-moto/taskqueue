# Task_Queue

[![CircleCI](https://circleci.com/gh/nozo-moto/taskqueue.svg?style=svg)](https://circleci.com/gh/nozo-moto/taskqueue)

This is task queue worker

# Usage

``` go
// set task queue interval
taskQueue := taskqueue.New(100 * time.Millisecond)
defer taskQueue.Stop()
// Run Task Queue Job in background.
go taskQueue.Run()

// Add Task, you can set retry coun
// if taskQueue.Stop called taskQueue.Add return err
err := taskQueue.Add(
	func() error {
		fmt.Printf("hoge\n")
		return nil
	},
	5,
)
```
