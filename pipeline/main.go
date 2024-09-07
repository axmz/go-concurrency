package main

import "fmt"

func gen(n int) <-chan int {
	out := make(chan int)
	go func() {
		for i := 0; i < n; i++ {
			out <- i
		}
		close(out)
	}()
	return out
}

func sum(in <-chan int, n int) <-chan int {
	out := make(chan int)
	go func() {
		for v := range in {
			out <- v + n
		}
		close(out)
	}()
	return out
}

func mult(in <-chan int, n int) <-chan int {
	out := make(chan int)
	go func() {
		for v := range in {
			out <- v * n
		}
		close(out)
	}()
	return out
}

func main() {
	g := gen(5)
	m := mult(g, 2)
	s := sum(m, 1)
	for v := range s {
		fmt.Println(v)
	}
}
