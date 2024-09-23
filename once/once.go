package main

import (
	"fmt"
	"sync"
)

func once_two_funcs() {
	var count int
	increment := func() { count++ }
	decrement := func() { count-- }
	var once sync.Once
	once.Do(increment)
	once.Do(decrement)
	fmt.Printf("Count: %d\n", count)
}

func once_deadlock() {
	var onceA, onceB sync.Once
	var initB func()
	initA := func() { fmt.Println("initA"); onceB.Do(initB) }
	initB = func() { fmt.Println("initB"); onceA.Do(initA) }
	onceA.Do(initA)
}

func main() {
	// once_deadlock()
	once_two_funcs()
}
