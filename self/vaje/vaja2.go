//go:build ignore

package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)
	ch_end := make(chan struct{ msg string })

	go func() {
		for {
			time.Sleep(200 * time.Millisecond )
			ch1 <- "VPSA"
		}
	}()

	go func() {
		for {
			time.Sleep(300 * time.Millisecond)
			ch2 <- "vpsa"
		}
	}()

	go func() {
		fmt.Scanln()
		ch_end <- struct{ msg string }{msg: "goodbye!"}
	}()

	for {
		select {
		case msg := <-ch1:
			{
				fmt.Print("ch1: ")
				fmt.Println(msg)
			}
		case msg := <-ch2:
			{
				fmt.Print("ch2: ")
				fmt.Println(msg)
			}
		case end := <-ch_end:
			{
				fmt.Println(end.msg)
				return
			}
		}
	}
}
