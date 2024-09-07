package main

import (
	"math/rand"
	"time"
)

const (
	lineCapacity = 10
	haircutTime  = time.Second * 1
)

func main() {
	rand.Seed(time.Now().Unix())
	barbershop := BarberShop{
		isOpen:  true,
		clients: make(chan string, lineCapacity),
	}

	barbershop.addBarber()
	barbershop.addBarber()

	go func() {
		_, ok := <-time.After(12 * time.Second)
		barbershop.isOpen = !ok
	}()

	go func() {
		for {
			if barbershop.isOpen {
				time.Sleep(1200 * time.Millisecond * time.Duration(rand.Intn(5)))
				barbershop.addClient()
			} else {
				return
			}
		}
	}()

	time.Sleep(18 * time.Second)
}
