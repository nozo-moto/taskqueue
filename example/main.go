package main

import (
	"fmt"
	"time"

	taskqueue "github.com/nozo-moto/taskqueue"
	"github.com/pkg/profile"
)

func main() {
	defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()

	taskQueue := taskqueue.New(100 * time.Millisecond)
	go taskQueue.Run()
	fmt.Println("task run ")

	var i int
	i = 0

	for {
		index := i + 1
		err := taskQueue.Add(
			func(args ...interface{}) error {
				index, ok := args[0].([]interface{})[0].(int)
				if !ok {
					panic("error assation not good")
				}
				fmt.Printf("hoge %d\n", index)

				return nil
			},
			3,
			index,
		)
		if err != nil {
			fmt.Println("error", err)
			break
		}
		err = taskQueue.Add(
			func(args ...interface{}) error {
				index, ok := args[0].([]interface{})[0].(int)
				if !ok {
					panic("hoge")
				}
				fmt.Printf("error %d\n", index)
				return fmt.Errorf("hoge %d", index)
			},
			2,
			index,
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
