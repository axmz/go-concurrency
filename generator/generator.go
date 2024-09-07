package main

import (
	"fmt"
	"time"
)

func fibonacci() <-chan int {
	ch := make(chan int)
	go func() {
		x, y := 1, 1
		for {
			ch <- x
			x, y = y, x+y
		}
	}()
	return ch
}

func main() {
	gen := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(<-gen)
	}
	time.Sleep(time.Second)
}
