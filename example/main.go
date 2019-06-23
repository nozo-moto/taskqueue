package main

import (
	"fmt"
	"time"

	taskqueue "github.com/nozo-moto/taskqueue"
)

func main() {
	taskQueue := taskqueue.New(100 * time.Millisecond)
	go taskQueue.Run()
	fmt.Println("task run ")

	var i int
	i = 0

	for {
		index := i + 1
		err := taskQueue.Add(
			func() error {
				fmt.Printf("hoge %d\n", index)
				return nil
			},
			3,
		)
		if err != nil {
			fmt.Println("error", err)
			break
		}
		err = taskQueue.Add(
			func() error {
				fmt.Printf("error %d\n", index)
				return fmt.Errorf("hoge %d", index)
			},
			2,
		)
		if err != nil {
			fmt.Println("error", err)
			break
		}
		time.Sleep(100 * time.Millisecond)
		i++

		if i > 100 {
			taskQueue.Stop()
			fmt.Println("task Stopped")
			break
		}
	}
}
