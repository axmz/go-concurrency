package main

import "fmt"

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

func fanin(chs ...<-chan int) chan int {
	out := make(chan int)
	go func() {
		for _, ch := range chs {
			for v := range ch {
				out <- v
			}
		}
		close(out)
	}()
	return out
}

func main() {
	g1 := gen(1, 3, 5)
	g2 := gen(2, 4, 6)

	f := fanin(g1, g2)
	for v := range f {
		fmt.Println(v)
	}
}
