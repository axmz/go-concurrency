package main

import (
	"fmt"
	"time"
)

type BarberShop struct {
	isOpen       bool
	clients      chan string
	clientsCount int
	barbersCount int
}

func (shop *BarberShop) addBarber() {
	shop.barbersCount++
	barber := shop.barbersCount
	go func() {
		isSleeping := false
		fmt.Printf("Barber #%d started his shift\n", barber)
		for {
			if !shop.isOpen && len(shop.clients) == 0 {
				fmt.Printf("Barber #%d goes home\n", barber)
				return
			}

			if len(shop.clients) == 0 {
				fmt.Printf("Barber #%d takes a nap\n", barber)
				if !isSleeping {
					isSleeping = true
				}
			}

			client := <-shop.clients
			fmt.Printf("Barber #%d haircuts %s\n", barber, client)
			time.Sleep(haircutTime)
			fmt.Printf("Barber #%d finished haircutting %s\n", barber, client)
		}
	}()

}

func (shop *BarberShop) addClient() {
	shop.clientsCount++
	if len(shop.clients) < lineCapacity {
		shop.clients <- fmt.Sprintf("Client #%d", shop.clientsCount)
		fmt.Printf("Client #%d joined the line\n", shop.clientsCount)
	} else {
		fmt.Printf("Client #%d went home because the line is full\n", shop.clientsCount)
	}
}
