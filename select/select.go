package main

import (
	"fmt"
	"time"
)

func select1() {
	s1 := make(chan int)
	s2 := make(chan int)
	server1 := func(c chan int) {
		for i := 0; i < 5; i++ {
			time.Sleep(time.Second)
			c <- i
		}
		close(c)
	}

	server2 := func(c chan int) {
		for i := 0; i < 5; i++ {
			time.Sleep(time.Second)
			c <- i
		}
		close(c)
	}

	go server1(s1)
	go server2(s2)

	// loop:
	for {
		select {
		case x1, ok := <-s1:
			if ok {
				fmt.Println("1: Reading from server 1", x1)
			} else {
				return
			}
		case x2, ok := <-s1:
			if ok {
				fmt.Println("2: Reading from server 1", x2)
			}
		case x3, ok := <-s2:
			if ok {
				fmt.Println("3: Reading from server 2", x3)
			}
		case x4, ok := <-s2:
			if ok {
				fmt.Println("4: Reading from server 2", x4)
			}
			// default:
			// 	fmt.Println("Fin")
			// 	break loop
		}
	}

}

func timeout() {
	var c <-chan int
	select {
	case <-c:
	case <-time.After(1 * time.Second):
		fmt.Println("Timed out.")
	}
}

func random() {
	c1 := make(chan interface{})
	close(c1)
	c2 := make(chan interface{})
	close(c2)
	var c1Count, c2Count int
	for i := 1000; i >= 0; i-- {
		select {
		case <-c1:
			c1Count++
		case <-c2:
			c2Count++
		}
	}
	fmt.Printf("c1Count: %d\nc2Count: %d\n", c1Count, c2Count)
}

func defaultCase() {
	done := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		close(done)
	}()
	workCounter := 0
loop:
	for {
		select {
		case <-done:
			break loop
		default:
		}
		// Simulate work
		workCounter++
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("Achieved %v cycles of work before signalled to stop.\n", workCounter)
}

func blockForever() {
	select {}
}

func main() {
	// select1()
	// defaultCase()
	// timeout()
	// random()
	blockForever()
	fmt.Println("exit")
}
