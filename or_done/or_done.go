package main

import (
	"fmt"
	"math/rand"
	"time"
)

func genRand() <-chan int {
	out := make(chan int)
	go func() {
		for {
			out <- rand.Intn(10)
		}
	}()
	return out
}

// wrapper
func orDone(done chan struct{}, c <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if !ok {
					return
				}
				select {
				case out <- v:
				case <-done:
				}
			}
		}
	}()
	return out
}

func main() {
	done := make(chan struct{})
	r := genRand()
	d := orDone(done, r)

	go func() {
		time.Sleep(time.Second * 1)
		close(done)
	}()

	// Doing this allows us to get back to simple for loops, like so:
	for val := range orDone(done, d) {
		fmt.Println(val)
		time.Sleep(time.Second)
	}

	time.Sleep(time.Second)
}
