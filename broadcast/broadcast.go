package main

import (
	"fmt"
	"sync"
)

func gen(n ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, v := range n {
			out <- v
		}
		close(out)
	}()
	return out
}

func broadcast(in <-chan int, out ...chan<- int) {
	var wg sync.WaitGroup
	g := func(c chan<- int, v int) {
		defer wg.Done()
		c <- v
	}

	for v := range in {
		for _, ch := range out {
			wg.Add(1)
			// creates a g on each v, this seems too much
			go g(ch, v)
		}
	}

	go func() {
		wg.Wait()
		for _, ch := range out {
			close(ch)
		}
	}()
}

func main() {
	g := gen(1, 2, 3, 4, 5, 6)
	t1, t2 := make(chan int), make(chan int)

	broadcast(g, t1, t2)

	go func() {
		for v := range t1 {
			fmt.Println("t1:", v)
		}
	}()

	go func() {
		for v := range t2 {
			fmt.Println("t2:", v)
		}
	}()
}
