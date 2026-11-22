//go:build ignore

package main

import (
	"fmt"
	"sync"
)

func square(x int, results chan<- int) {
	results <- x * x
	wg.Done()
}
var wg sync.WaitGroup
func main() {
	results := make(chan int)

	for i := 1; i <= 300000; i++ {
		wg.Add(1)
		go square(i, results)
	}
	go func() {
		wg.Wait()
		close(results)
	}()
	var sum float64
	for r := range results {
		sum += float64(r)
	}
	fmt.Printf("%.3v", sum)
}
