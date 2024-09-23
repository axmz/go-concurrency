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

func tee(done <-chan struct{}, in <-chan int) (_, _ <-chan int) {
	out1 := make(chan int)
	out2 := make(chan int)

	go func() {
		defer close(out1)
		defer close(out2)
		for val := range in {
			var out1, out2 = out1, out2 // avoid closing same channel twice
			select {
			case out1 <- val:
				// out1 = nil
			case <-done:
			}
			select {
			case out2 <- val:
			case <-done:
			}
		}
	}()

	return out1, out2
}

func tee2(
	done <-chan interface{},
	in <-chan int,
) (_, _ <-chan interface{}) {
	out1 := make(chan interface{})
	out2 := make(chan interface{})

	go func() {
		defer close(out1)
		defer close(out2)
		// for val := range orDone(done, in) {
		for val := range in {
			var out1, out2 = out1, out2
			for i := 0; i < 2; i++ {
				select {
				case <-done:
				case out1 <- val:
					out1 = nil
				case out2 <- val:
					out2 = nil
				}
			}
		}
	}()

	return out1, out2
}

func main() {
	g := gen(1, 2, 3, 4, 5, 6)

	// out1, out2 := tee(nil, g)
	out1, out2 := tee2(nil, g)

	go func() {
		for v := range out1 {
			fmt.Println("t1:", v)
		}
	}()

	go func() {
		for v := range out2 {
			fmt.Println("t2:", v)
		}
	}()

	time.Sleep(time.Second)
}
