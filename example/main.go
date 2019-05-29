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

	i := 0

	for {
		index := i + 1
		taskQueue.Add(
			func() error {
				fmt.Printf("hoge %d\n", index)
				return nil
			},
			5,
		)
		taskQueue.Add(
			func() error {
				fmt.Printf("error %d\n", index)
				return fmt.Errorf("hoge %d", index)
			},
			2,
		)
		time.Sleep(100 * time.Millisecond)
		i++
	}
}
