package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	openHours = time.Second * 4
	queue     = 3
	cutTime   = time.Second * 1
)

type BarberShop struct {
	c      chan int
	closeC chan bool
	wg     *sync.WaitGroup
}

func (b *BarberShop) barber() {
	b.wg.Add(1)
	defer b.wg.Done()
	for i := range b.c {
		fmt.Println("Barber is cutting hair", i)
		time.Sleep(cutTime)
	}
	// for {
	// 	select {
	// 	case i, ok := <-b.c:
	// 		if !ok {
	// 			b.wg.Done()
	// 			return
	// 		}
	// 		fmt.Println("Barber is cutting hair", i)
	// 		time.Sleep(cutTime)
	// 	default:
	// 		fmt.Printf("z")
	// 		time.Sleep(time.Millisecond * 100)
	// 	}
	// }
}

func (b *BarberShop) customer(i int) {
	select {
	case b.c <- i:
		fmt.Println("Customer arrived", i)
	default:
		fmt.Println("Queue is full", i)
	}
}

func (b *BarberShop) close() {
	<-time.After(openHours)
	fmt.Println("Shop is closing")
	b.closeC <- true
}

func (b *BarberShop) open() {
	go b.barber()
	for i := 0; ; i++ {
		select {
		case <-b.closeC:
			close(b.c)
			return
		default:
			b.customer(i)
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
		}
	}
}

func main() {
	b := BarberShop{
		c:      make(chan int, queue),
		wg:     &sync.WaitGroup{},
		closeC: make(chan bool),
	}

	go b.open()
	b.close()
	b.wg.Wait()
	fmt.Println("CLOSED")
}
