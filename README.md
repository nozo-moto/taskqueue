# Task_Queue

[![CircleCI](https://circleci.com/gh/nozo-moto/taskqueue.svg?style=svg)](https://circleci.com/gh/nozo-moto/taskqueue)

This is task queue worker. I create referencing TaskQueue of GCP.
You can set task's interval time, and retry count.

# How To Use

1. Add task to goroutine worker.
2. Run Task in goroutine worker. If error happen, task will be retried as many times as retry count.

## Sample Code
``` go
// set task queue interval
taskQueue := taskqueue.New(100 * time.Millisecond)
defer taskQueue.Stop()
// Run Task Queue Job in background.
go taskQueue.Run()

// Add Task, you can set retry coun
// if taskQueue.Stop called taskQueue.Add return err
err := taskQueue.Add(
	func(args ...interface{}) error {
        // Get Agument
		index, ok := args[0].([]interface{})[0].(int) 
		if !ok {
			panic("error assation not good")
		}
		fmt.Printf("hoge %d\n", index)

		return nil
	},
	3, // Retry Count
	index, // Argument 
)

```
