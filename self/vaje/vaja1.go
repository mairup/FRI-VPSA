//go:build ignore

package main

import "fmt"

const N int = 5

var ch chan struct{} = make(chan struct{})

func worker(id int) {
	fmt.Println(id)
	ch <- struct{}{}
}

func main() {
	for range N {
		for i := range N {
			go worker(i)
		}

		for range N{
			<-ch
		}
	}
	close(ch)
}
