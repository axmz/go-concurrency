package main

import (
	"fmt"
	"sync"
	"time"
)

func signal() {
	c := sync.NewCond(&sync.Mutex{})
	queue := make([]interface{}, 0, 10)
	removeFromQueue := func(delay time.Duration) {
		time.Sleep(delay)
		c.L.Lock()
		queue = queue[1:]
		fmt.Println("Removed from queue")
		c.L.Unlock()
		c.Signal()
	}
	for i := 0; i < 10; i++ {
		c.L.Lock()
		for len(queue) == 2 {
			fmt.Println("for")
			c.Wait()
		}
		fmt.Println("Adding to queue")
		queue = append(queue, struct{}{})
		go removeFromQueue(1 * time.Second)
		c.L.Unlock()
	}
}

func broadcast() {
	type Button struct {
		Clicked *sync.Cond
	}

	button := Button{Clicked: sync.NewCond(&sync.Mutex{})}

	subscribe := func(c *sync.Cond, fn func()) {
		var goroutineRunning sync.WaitGroup
		goroutineRunning.Add(1)
		go func() {
			goroutineRunning.Done()
			c.L.Lock()
			defer c.L.Unlock()
			c.Wait()
			fn()
		}()
		goroutineRunning.Wait()
	}

	var clickRegistered sync.WaitGroup

	clickRegistered.Add(3)
	subscribe(button.Clicked, func() {
		fmt.Println("Maximizing window.")
		clickRegistered.Done()
	})

	subscribe(button.Clicked, func() {
		fmt.Println("Displaying annoying dialog box!")
		clickRegistered.Done()
	})

	subscribe(button.Clicked, func() {
		fmt.Println("Mouse clicked.")
		clickRegistered.Done()
	})

	button.Clicked.Broadcast()

	clickRegistered.Wait()
}

func broadcast2() {
	var wg sync.WaitGroup
	var mu sync.Mutex
	cond := sync.NewCond(&mu)

	// Five waiting goroutines
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			mu.Lock()
			fmt.Printf("Goroutine %d waiting\n", id)
			cond.Wait() // Wait for broadcast signal
			fmt.Printf("Goroutine %d received signal\n", id)
			mu.Unlock()
		}(i)
	}

	// Simulate some work before broadcasting
	time.Sleep(2 * time.Second)

	mu.Lock()
	fmt.Println("Broadcasting signal to all goroutines")
	cond.Broadcast() // Notify all waiting goroutines
	mu.Unlock()

	// Wait for all goroutines to finish
	wg.Wait()
	fmt.Println("All goroutines have been signaled")
}

func main() {
	// signal()
	// broadcast()
	broadcast2()
}
