package main

import (
	"fmt"
	"time"

	taskqueue "github.com/nozo-moto/taskqueue"
)

func main() {
	var err error
	taskQueue := taskqueue.New(100 * time.Millisecond)
	go taskQueue.Run()
	defer func() {
		taskQueue.Stop()

		for {
		}
	}()

	fmt.Println("task run ")

	go func() {
	L:
		for {
			select {
			case err, ok := <-taskQueue.Error:
				if !ok {
					fmt.Println("Finished Taskqueue")
					break L
				}
				fmt.Println("error cause  ", err.Err, err.Task.RetryTimes, ok)
			}
		}
	}()

	err = taskQueue.Add(
		func() error {
			fmt.Printf("hoge %d\n", 1)
			return nil
		},
		1,
	)
	if err != nil {
		panic(err)
	}
	err = taskQueue.Add(
		func() error {
			return fmt.Errorf("error %d", 1)
		},
		3,
	)
	if err != nil {
		panic(err)
	}
	err = taskQueue.Add(
		func() error {
			fmt.Printf("hoge %d\n", 2)
			return nil
		},
		1,
	)
	if err != nil {
		panic(err)
	}
	err = taskQueue.Add(
		func() error {
			return fmt.Errorf("error %d", 2)
		},
		3,
	)
	if err != nil {
		panic(err)
	}
	err = taskQueue.Add(
		func() error {
			fmt.Printf("hoge %d\n", 3)
			return nil
		},
		1,
	)
	if err != nil {
		panic(err)
	}
	err = taskQueue.Add(
		func() error {
			return fmt.Errorf("error %d", 3)
		},
		3,
	)
	if err != nil {
		panic(err)
	}

}
