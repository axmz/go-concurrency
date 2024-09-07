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

func fanout(ch <-chan int, workers int) chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			for v := range ch {
				fmt.Println("worker:", i, "processes:", v)
				out <- v * v
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	g1 := gen(1, 2, 3, 4, 5, 6)

	f := fanout(g1, 3)
	for v := range f {
		fmt.Println(v)
	}
}
