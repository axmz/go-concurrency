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

type Semaphore struct {
	sem chan struct{}
}

func (s *Semaphore) acquire() {
	s.sem <- struct{}{}
}

func (s *Semaphore) release() {
	<-s.sem
}

func main() {
	var wg sync.WaitGroup
	workers := 10
	res := make(chan int)

	g := gen(1, 2, 3, 4, 5)

	s := Semaphore{
		sem: make(chan struct{}, 3),
	}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.acquire()
			defer s.release()
			v, ok := <-g
			if !ok {
				return
			}
			fmt.Println("goroutine", i, "acquired the semaphore with value", v)
			fmt.Println("goroutine", i, "calculated", v*v)
			res <- v * v
			fmt.Println("goroutine", i, "released the semaphore")
		}()
	}

	wg.Add(1)
	go func() {
		wg.Done()
		for v := range res {
			fmt.Println(v)
		}
		close(res)
	}()

	wg.Wait()
}
