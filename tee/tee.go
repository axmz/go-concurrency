package main

import (
	"fmt"
	"time"
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

func tee(in <-chan int, out ...chan<- int) {
	go func() {
		for v := range in {
			for _, ch := range out {
				ch <- v
			}
		}
		for _, ch := range out {
			close(ch)
		}
	}()
}

func main() {
	g := gen(1, 2, 3, 4, 5, 6)
	t1, t2 := make(chan int), make(chan int)

	tee(g, t1, t2)

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

	time.Sleep(time.Second)
}
