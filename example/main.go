package main

import (
	"fmt"
	"time"

	taskqueue "github.com/nozo-moto/taskqueue"
)

func main() {
	taskQueue := taskqueue.New(1000 * time.Millisecond)
	go taskQueue.Run()
	defer func() {
	R:
		for {
			if taskQueue.Count() == 0 {
				break R
			}
		}
		taskQueue.Stop()
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
				fmt.Println("error cause  ", err)
			}
		}
	}()

	for i := 0; i < 10; i++ {
		err := taskQueue.Add(
			func(args ...any) error {
				fmt.Printf("hoge %d\n", args)
				return nil
			},
			1, 2, 3, 4, "one", "two", "three",
		)
		if err != nil {
			panic(err)
		}
	}
}
