package main

import (
	"fmt"
	"time"
)

func server1(c chan int) {
	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		c <- i
	}
	close(c)
}

func server2(c chan int) {
	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		c <- i
	}
	close(c)
}

func main() {
	s1 := make(chan int)
	s2 := make(chan int)

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
