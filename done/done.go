package main

import (
	"fmt"
	"math/rand"
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

func genRand(done chan struct{}) <-chan int {
	out := make(chan int)
	go func() {
		for {
			select {
			case <-done:
				close(out)
			case out <- rand.Intn(10):
			}
		}
	}()
	return out
}

func leak() {
	doWork := func(strings <-chan string) <-chan interface{} {
		completed := make(chan interface{})
		go func() {
			fmt.Println("doWork")
			defer fmt.Println("doWork never exited.")
			defer close(completed)
			for s := range strings {
				// Do something interesting
				fmt.Println(s)
			}
			fmt.Println("if strings is not closed, this goroutine will never exit")
		}()
		return completed
	}

	// s := make(chan string)
	// go func() {
	// 	s <- "x"
	// 	close(s)
	// }()
	// doWork(s)

	doWork(nil)
	time.Sleep(time.Second)
	fmt.Println("Done.")

	// // Here we see that the main goroutine passes a nil channel into doWork. Therefore, the
	// // strings channel will never actually gets any strings written onto it, and the goroutine
	// // containing doWork will remain in memory for the lifetime of this process (we would
	// // even deadlock if we joined the goroutine within doWork and the main goroutine)
}

func leak_no_close() {
	doWork := func(strings <-chan string) <-chan interface{} {
		completed := make(chan interface{})
		go func() {
			fmt.Println("doWork")
			defer fmt.Println("doWork never exited.")
			defer close(completed)
			for s := range strings {
				// Do something interesting
				fmt.Println(s)
			}
			fmt.Println("if strings is not closed, this goroutine will never exit")
		}()
		return completed
	}

	s := make(chan string)
	go func() {
		s <- "x"
		// close(s)
	}()
	doWork(s)

	time.Sleep(time.Second)
	fmt.Println("Done.")
}

func leak_close() {
	doWork := func(strings <-chan string) <-chan interface{} {
		completed := make(chan interface{})
		go func() {
			fmt.Println("doWork")
			defer fmt.Println("doWork exited.")
			defer close(completed)
			for s := range strings {
				fmt.Println(s)
			}
			fmt.Println("if strings is not closed, this goroutine will never exit")
		}()
		return completed
	}

	s := make(chan string)
	go func() {
		s <- "x"
		close(s)
	}()
	doWork(s)

	time.Sleep(time.Second)
	fmt.Println("Done.")
}

func done() {
	doWork := func(done <-chan interface{}, strings <-chan string) <-chan interface{} {
		terminated := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(terminated)
			for {
				select {
				case s := <-strings:
					fmt.Println(s)
				case <-done:
					return
				}
			}
		}()
		return terminated
	}

	done := make(chan interface{})
	terminated := doWork(done, nil)

	go func() {
		// Cancel the operation after 1 second.
		time.Sleep(1 * time.Second)
		fmt.Println("Canceling doWork goroutine...")
		close(done)
	}()

	<-terminated
	fmt.Println("Done.")

}

func done2() {
	doWork := func(done <-chan struct{}, data <-chan int) <-chan struct{} {
		terminated := make(chan struct{})
		go func() {
			defer fmt.Println("exit doWork")
			defer close(terminated)
			for {
				select {
				case v := <-data:
					fmt.Println(v)
					time.Sleep(time.Second)
				case <-done:
					return
				}
			}

		}()
		return terminated
	}

	done := make(chan struct{})
	g := gen(1, 2, 3, 4)
	terminated := doWork(done, g)

	go func() {
		time.Sleep(time.Second)
		done <- struct{}{}
		close(done)
	}()

	<-terminated
	fmt.Println("done")
}

func done3() {
	g := gen(1, 2, 3, 4, 5)
	done := make(chan struct{})
	go func() {
		<-time.After(time.Second)
		done <- struct{}{}
		close(done)
	}()

	for {
		select {
		case v, ok := <-g:
			if !ok {
				fmt.Println("but, sometimes returns here. and it is not really intuitive")
				return
			}
			fmt.Println(v)
			time.Sleep(time.Millisecond * 200)
		case <-done:
			fmt.Println("i would expect it to always return here")
			return
		}
	}
}

func done4() {
	done := make(chan struct{})
	for i := 0; i < 3; i++ {
		fmt.Println(<-genRand(done))
	}
	// if genRand wouldn't accept done it would remain dangling as there is not way to close it.
	close(done)
	time.Sleep(time.Second)
}

func main() {
	// leak()
	// leak_no_close()
	// leak_close()
	// done()
	// done2()
	// done3()
	done4()
}
